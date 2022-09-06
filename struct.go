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
	UseHttpServer        bool                                                // 是否使用http服务
	HttpPort             int                                                 // http服务的端口
	chromeExecPath       string                                              // 指定运行浏览器
	HttpRoute            map[string]func(http.ResponseWriter, *http.Request) // 额外的http路由
	AssetFile            http.FileSystem
	RestoreAssets        func(dir, name string) error // 解压资源方法
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

type PackageConf struct {
	Name           string   `json:"name"`
	ChromeExecPath string   `json:"chrome_exec_path"`
	ChromePackPath Platform `json:"chrome_pack_path"`
	ChromeVersion  int      `json:"chrome_version"`
	BuildCachePath string   `json:"build_cache_path"`
	Icons          Platform `json:"icons"`
	RunBuildPath   string   `json:"run_build_path"`
}

type Platform struct {
	Linux   string `json:"linux"`
	Windows string `json:"windows"`
	Darwin  string `json:"darwin"`
}
