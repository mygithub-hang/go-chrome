package go_chrome

import "github.com/chromedp/chromedp"

// GoChromeOptions 启动配置
type GoChromeOptions struct {
	CliModule            bool   // 命令行模式
	AppModule            bool   // 应用模式启动
	WindowWidth          int    // 宽度
	WindowHeight         int    // 高度
	WindowPositionWidth  int    // 横向位置
	WindowPositionHeight int    // 竖向位置
	ChromeExecPath       string // 指定运行的浏览器
}

type ActionTask chromedp.Tasks

type jsToGoData struct {
	Name               string `json:"name"`
	Payload            string `json:"payload"`
	ExecutionContextId int    `json:"executionContextId"`
}

type jsToGoPayload struct {
	Name string        `json:"name"`
	Seq  int           `json:"seq"`
	Args []interface{} `json:"args"`
}

type goToJs struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
