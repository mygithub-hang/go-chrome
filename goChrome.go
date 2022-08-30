package go_chrome

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/inspector"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"os"
)

var goChromeCtxMap = make(map[string]context.Context, 0)

type GoChrome struct {
	windowName        string
	Url               string
	ContextContext    context.Context
	ContextCancelFunc context.CancelFunc
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
	}
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *inspector.EventDetached:
			// 界面关闭
			newWindow.ContextCancelFunc()
			os.Exit(0)
		case *runtime.EventBindingCalled:
			// 处理函数
			json, _ := ev.MarshalJSON()
			//{"name":"gogogogo","payload":"{\"name\":\"gogogogo\",\"seq\":1,\"args\":[\"aa\",\"aas\"]}","executionContextId":1}
			fmt.Println(string(json))
		}
	})
	return newWindow
}

func (gc *GoChrome) Run() {
	if err := chromedp.Run(gc.ContextContext, gc.start()); err != nil {
		log.Fatal(err)
		return
	}
	select {}
}

func (gc *GoChrome) start() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(gc.Url),
	}
}

func (gc *GoChrome) ListenTarget(fn func(ev interface{})) {
	chromedp.ListenTarget(gc.ContextContext, fn)
}

func (gc *GoChrome) ListenBrowser(fn func(ev interface{})) {
	chromedp.ListenTarget(gc.ContextContext, fn)
}
