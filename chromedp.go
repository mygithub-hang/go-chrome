package go_chrome

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
)

func getCtX(windowName, url string, opt ...GoChromeOptions) (context.Context, context.CancelFunc, *GoChromeOptions, string) {
	ctxInit := context.WithValue(context.Background(), "windowsName", windowName)
	runOpt := GoChromeOptions{}
	if len(opt) > 0 {
		runOpt = opt[0]
		ChromeRunCommand["headless"] = runOpt.CliModule
		url = GetUrl(url, runOpt.UseHttpServer)
		if runOpt.AppModule || strings.HasPrefix(url, "data:text/html") {
			ChromeRunCommand["app"] = url
		}
		if runOpt.WindowWidth > 0 && runOpt.WindowHeight > 0 {
			ChromeRunCommand["window-size"] = fmt.Sprintf("%d,%d", runOpt.WindowWidth, runOpt.WindowHeight)
		}
		if runOpt.WindowPositionWidth > 0 && runOpt.WindowPositionHeight > 0 {
			ChromeRunCommand["window-position"] = fmt.Sprintf("%d,%d", runOpt.WindowPositionWidth, runOpt.WindowPositionHeight)
		}
	} else {
		url = GetUrl(url, runOpt.UseHttpServer)
	}
	execAllocatorOption := []chromedp.ExecAllocatorOption{}
	if runOpt.ChromeExecPath != "" {
		execAllocatorOption = append(execAllocatorOption, chromedp.ExecPath(runOpt.ChromeExecPath))
	}
	for k, v := range ChromeRunCommand {
		execAllocatorOption = append(execAllocatorOption, chromedp.Flag(k, v))
	}
	ctx, _ := chromedp.NewExecAllocator(
		ctxInit,
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			execAllocatorOption...,
		)...,
	)
	newCtx, cancel := chromedp.NewContext(
		ctx,
		// 设置日志方法
		chromedp.WithLogf(log.Printf),
	)
	return newCtx, cancel, &runOpt, url
}
