package go_chrome

import (
	"github.com/chromedp/chromedp"
	"net/http"
)

// GoChromeOptions 启动配置
type GoChromeOptions struct {
	CliModule            bool                                                // 命令行模式
	AppModule            bool                                                // 应用模式启动
	WindowWidth          int                                                 // 宽度
	WindowHeight         int                                                 // 高度
	WindowPositionWidth  int                                                 // 横向位置
	WindowPositionHeight int                                                 // 竖向位置
	ChromeExecPath       string                                              // 指定运行的浏览器内核
	BrowserRunPath       Platform                                            // 浏览器内核解压目录
	UseHttpServer        bool                                                // 是否使用http服务
	HttpPort             int                                                 // http服务的端口
	HttpRoute            map[string]func(http.ResponseWriter, *http.Request) // 额外的http路由
	DefHttpIndexData     interface{}                                         // 默认渲染主页模板参数
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

type Platform struct {
	Linux   string
	Windows string
	Darwin  string
}
