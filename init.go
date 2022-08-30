package go_chrome

// ChromeRunCommand 浏览器启动命令
var ChromeRunCommand map[string]interface{} = map[string]interface{}{
	"disable-background-networking":          true,
	"disable-background-timer-throttling":    true,
	"disable-backgrounding-occluded-windows": true,
	"disable-breakpad":                       true,
	"disable-client-side-phishing-detection": true,
	"disable-default-apps":                   true,
	"disable-dev-shm-usage":                  true,
	"disable-infobars":                       true,
	"disable-extensions":                     true,
	"disable-hang-monitor":                   true,
	"disable-ipc-flooding-protection":        true,
	"disable-popup-blocking":                 true,
	"disable-prompt-on-repost":               true,
	"disable-renderer-backgrounding":         true,
	"disable-sync":                           true,
	"disable-translate":                      true,
	"disable-windows10-custom-titlebar":      true,
	"metrics-recording-only":                 true,
	"no-first-run":                           true,
	"no-default-browser-check":               true,
	"safebrowsing-disable-auto-update":       true,
	"enable-automation":                      false,
	"password-store":                         "basic",
	"use-mock-keychain":                      true,
}
