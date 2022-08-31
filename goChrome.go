package go_chrome

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/inspector"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"reflect"
)

var goChromeCtxMap = make(map[string]context.Context, 0)

type GoChrome struct {
	windowName        string
	Url               string
	ContextContext    context.Context
	ContextCancelFunc context.CancelFunc
	bindFunc          map[string]interface{}
	Action            chromedp.Tasks
	finalAction       chromedp.Tasks
	afterFunc         func()
	closeChan         chan int
}

func Create(url string, opt ...GoChromeOptions) *GoChrome {
	if url == "" {
		url = "data:text/html,<html></html>"
	}
	ctx, cancel := getCtX("index", url, opt...)
	newWindow := &GoChrome{
		windowName:        "index",
		Url:               url,
		ContextContext:    ctx,
		ContextCancelFunc: cancel,
		bindFunc:          map[string]interface{}{},
		Action:            chromedp.Tasks{},
		finalAction:       chromedp.Tasks{chromedp.Navigate(url)},
		closeChan:         make(chan int, 1),
	}
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *inspector.EventDetached:
			// 界面关闭
			newWindow.Close()
		case *runtime.EventBindingCalled:
			// 处理函数
			jsonByte, err := ev.MarshalJSON()
			if err == nil {
				newWindow.runBindFunc(jsonByte)
			}
		}
	})
	return newWindow
}

func (gc *GoChrome) Run() {
	if err := chromedp.Run(gc.ContextContext, gc.start()); err != nil {
		log.Fatal(err)
		return
	}
	defer gc.Close()
	for {
		select {
		case <-gc.closeChan:
			gc.ContextCancelFunc()
			os.Exit(0)
		}
	}
}

func (gc *GoChrome) RunTask(myTask chromedp.Tasks) {
	if err := chromedp.Run(gc.ContextContext, myTask); err != nil {
		log.Fatal(err)
		return
	}
	defer gc.Close()
	for {
		select {
		case <-gc.closeChan:
			gc.ContextCancelFunc()
			os.Exit(0)
		}
	}
}

func (gc *GoChrome) start() chromedp.Tasks {
	gc.finalAction = append(gc.finalAction, chromedp.ActionFunc(func(ctx context.Context) error {
		gc.ContextContext = ctx
		gc.bindJsFunc()
		return nil
	}))
	gc.finalAction = append(gc.finalAction, gc.Action...)
	gc.finalAction = append(gc.finalAction, chromedp.ActionFunc(func(ctx context.Context) error {
		gc.ContextContext = ctx
		go gc.afterFunc()
		return nil
	}))
	return gc.finalAction
}

func (gc *GoChrome) SetAction(actionArr ActionTask) {
	gc.Action = chromedp.Tasks(actionArr)
}

func (gc *GoChrome) Close() {
	gc.closeChan <- 1
}

func (gc *GoChrome) ListenTarget(fn func(ev interface{})) {
	chromedp.ListenTarget(gc.ContextContext, fn)
}

func (gc *GoChrome) ListenBrowser(fn func(ev interface{})) {
	chromedp.ListenTarget(gc.ContextContext, fn)
}

func (gc *GoChrome) Bind(jsFuncName string, f interface{}) error {
	v := reflect.ValueOf(f)
	// f must be a function
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}
	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}
	gc.bindFunc[jsFuncName] = f
	return nil
}

func (gc *GoChrome) OpenAfter(f func()) {
	gc.afterFunc = f
}

func (gc *GoChrome) JsFunc(funcName string, args ...string) reflect.Value {
	param := ""
	for _, v := range args {
		if param == "" {
			param = `'` + v + `'`
		} else {
			param += `, '` + v + `'`
		}
	}
	fn := fmt.Sprintf("%s(%s);", funcName, param)
	eval := runtime.Evaluate(fn)
	eval.AwaitPromise = true
	eval.ReturnByValue = true
	do, _, _ := eval.Do(gc.ContextContext)
	jsonStr, _ := do.MarshalJSON()
	jsonData := goToJs{}
	err := json.Unmarshal(jsonStr, &jsonData)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
		return reflect.Value{}
	}
	return reflect.ValueOf(jsonData.Value)
}

func (gc *GoChrome) bindJsFunc() {
	for fn, _ := range gc.bindFunc {
		_ = runtime.AddBinding(fn).Do(gc.ContextContext)
		script := fmt.Sprintf(`(() => {
	const bindingName = '%s';
	const binding = window[bindingName];
	window[bindingName] = async (...args) => {
		const me = window[bindingName];
		let errors = me['errors'];
		let callbacks = me['callbacks'];
		if (!callbacks) {
			callbacks = new Map();
			me['callbacks'] = callbacks;
		}
		if (!errors) {
			errors = new Map();
			me['errors'] = errors;
		}
		const seq = (me['lastSeq'] || 0) + 1;
		me['lastSeq'] = seq;
		const promise = new Promise((resolve, reject) => {
			callbacks.set(seq, resolve);
			errors.set(seq, reject);
		});
		binding(JSON.stringify({name: bindingName, seq, args}));
		return promise;
	}})();
	`, fn)
		_, _ = page.AddScriptToEvaluateOnNewDocument(script).Do(gc.ContextContext)
		eval := runtime.Evaluate(script)
		eval.AwaitPromise = true
		eval.ReturnByValue = true
		_, _, _ = eval.Do(gc.ContextContext)
	}
}

func (gc *GoChrome) runBindFunc(jsonStr []byte) {
	jsonData := jsToGoData{}
	err := json.Unmarshal(jsonStr, &jsonData)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
		//gc.Close()
		return
	}
	payload := jsToGoPayload{}
	err = json.Unmarshal([]byte(jsonData.Payload), &payload)
	if err != nil {
		fmt.Println("Umarshal failed payload:", err)
		//gc.Close()
		return
	}
	fv := reflect.ValueOf(gc.bindFunc[jsonData.Name])
	realParams := make([]reflect.Value, 0) //参数
	for _, item := range payload.Args {
		realParams = append(realParams, reflect.ValueOf(item))
	}
	res := fv.Call(realParams)
	fmt.Println(res)
}
