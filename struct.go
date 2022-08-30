package go_chrome

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
