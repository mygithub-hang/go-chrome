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
	osRuntime "runtime"
)

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
	GoHttp            *goHttp
}

func Create(url string, opt ...GoChromeOptions) *GoChrome {
	ctx, cancel, runOpt, runUrl := getCtX("index", url, opt...)
	var goHttp *goHttp = nil
	if runOpt.UseHttpServer {
		goHttp = getHttp(runOpt)
		go func() {
			err := goHttp.StartHttpServer()
			if err != nil {
				fmt.Println(err)
				cancel()
				os.Exit(0)
			}
		}()
	}
	//
	err := ReleaseBrowser(runOpt)
	if err != nil {
		fmt.Println(err)
		cancel()
		os.Exit(0)
		return nil
	}
	newWindow := &GoChrome{
		windowName:        "index",
		Url:               runUrl,
		ContextContext:    ctx,
		ContextCancelFunc: cancel,
		bindFunc:          map[string]interface{}{},
		Action:            chromedp.Tasks{},
		finalAction:       chromedp.Tasks{chromedp.Navigate(runUrl)},
		closeChan:         make(chan int, 1),
		GoHttp:            goHttp,
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

func (gc *GoChrome) RunUnBackup() error {
	return chromedp.Run(gc.ContextContext, gc.start())
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
		if gc.afterFunc != nil {
			go gc.afterFunc()
		}
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
		return errors.New("function may only return a value or a value error")
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
		return
	}
	payload := jsToGoPayload{}
	err = json.Unmarshal([]byte(jsonData.Payload), &payload)
	if err != nil {
		fmt.Println("Umarshal failed payload:", err)
		return
	}
	fv := reflect.ValueOf(gc.bindFunc[jsonData.Name])
	realParams := make([]reflect.Value, 0) //参数
	for _, item := range payload.Args {
		realParams = append(realParams, reflect.ValueOf(item))
	}
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	res := fv.Call(realParams)
	var jsRes string
	var jsErr error
	switch len(res) {
	case 0:
	case 1:
		// One result may be a value, or an error
		if res[0].Type().Implements(errorType) {
			if res[0].Interface() != nil {
				jsErr = res[0].Interface().(error)
			}
		}
		if _, ok := res[0].Interface().(string); ok {
			jsRes = res[0].String()
		} else if _, ok := res[0].Interface().(float64); ok {
			jsRes = fmt.Sprintf("%.f", res[0].Float())
		} else if _, ok := res[0].Interface().([]byte); ok {
			jsRes = string(res[0].Bytes())
		} else {
			jsResByte, errJs := json.Marshal(res[0].Interface())
			jsErr = errJs
			jsRes = string(jsResByte)
		}
	default:
		resArr := []interface{}{}
		for _, v := range res {
			resArr = append(resArr, v.Interface())
		}
		jsResByte, errJs := json.Marshal(resArr)
		jsErr = errJs
		jsRes = string(jsResByte)
	}
	go func() {
		jsString := func(v string) string { b, _ := json.Marshal(v); return string(b) }
		retErr := ""
		if jsErr != nil {
			retErr = jsErr.Error()
		}
		expr := fmt.Sprintf(`
							if (%[4]s) {
								window['%[1]s']['errors'].get(%[2]d)(%[4]s);
							} else {
								window['%[1]s']['callbacks'].get(%[2]d)(%[3]s);
							}
							window['%[1]s']['callbacks'].delete(%[2]d);
							window['%[1]s']['errors'].delete(%[2]d);
							`, payload.Name, payload.Seq, jsString(jsRes), jsString(retErr))
		eval := runtime.Evaluate(expr)
		eval.AwaitPromise = true
		eval.ReturnByValue = true
		_, _, err = eval.Do(gc.ContextContext)
	}()
}

func ReleaseBrowser(opt *GoChromeOptions) error {
	if opt.RestoreAssets == nil {
		return nil
	}
	packConf := getPackageConfig()
	if packConf != nil {
		if !packConf.IntegratedBrowser || packConf.ChromeExecPath != "" {
			return nil
		}
	}
	browserPath, getPathErr := GetBrowserPath(opt)
	if getPathErr != nil {
		return getPathErr
	}
	err := createDir(browserPath)
	if err != nil {
		return err
	}
	releasePath := ""
	sysType := osRuntime.GOOS
	switch sysType {
	case "darwin":
		releasePath = browserPath + "/browser/chrome-mac.zip"
	case "linux":
		releasePath = browserPath + "/browser/chrome-linux.zip"
	case "windows":
		releasePath = browserPath + "/browser/chrome-win.zip"
	}
	if !IsExist(releasePath) {
		err = createDir(releasePath)
		if err != nil {
			return err
		}
		err = opt.RestoreAssets(browserPath, "browser")
		if err != nil {
			return err
		}
	}
	if sysType == "windows" {
		runPah := getChromeExecPath(opt)
		if !IsExist(runPah) {
			return UnPackZip(releasePath, browserPath)
		}
	}
	return nil
}

func GetBrowserPath(opt *GoChromeOptions) (string, error) {
	if opt.AppName == "" {
		fmt.Println("AppName is empty")
		os.Exit(0)
		return "", errors.New("AppName is empty")
	}
	return GetHomePath() + "/AppData/Roaming/" + opt.AppName, nil
}
