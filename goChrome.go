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
	}
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *inspector.EventDetached:
			// 界面关闭
			newWindow.ContextCancelFunc()
			os.Exit(0)
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
	gc.bindJsFunc()
	select {}
}

func (gc *GoChrome) start() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(gc.Url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			gc.ContextContext = ctx
			return nil
		}),
	}
}

func (gc *GoChrome) Close() {
	gc.ContextCancelFunc()
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

func (gc *GoChrome) JsFunc(funcName string, args ...string) map[string]interface{} {
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
	jsonMap := make(map[string]interface{}, 0)
	err := json.Unmarshal(jsonStr, &jsonMap)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
		gc.Close()
	}
	return jsonMap
	// {"type":"string","value":"ass"}
	// {"type":"object","value":["ass","ddd",122]}
	// {"type":"object","value":{"aa":"ass","dd":"ddd","cc":122}}
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
	//{"name":"gogogogo","payload":"{\"name\":\"gogogogo\",\"seq\":1,\"args\":[\"aa\",\"aas\"]}","executionContextId":1}
	jsonMap := make(map[string]string, 0)
	_ = json.Unmarshal(jsonStr, &jsonMap)
	//if err != nil {
	//	fmt.Println("Umarshal failed1:", err)
	//	gc.Close()
	//}
	paramMap := make(map[string][]interface{})
	_ = json.Unmarshal([]byte(jsonMap["payload"]), &paramMap)
	//if err != nil {
	//	fmt.Println("Umarshal failed2:", err)
	//	gc.Close()
	//}
	fv := reflect.ValueOf(gc.bindFunc[jsonMap["name"]])
	realParams := make([]reflect.Value, 0) //参数
	for _, item := range paramMap["args"] {
		realParams = append(realParams, reflect.ValueOf(item))
	}
	_ = fv.Call(realParams)
}
