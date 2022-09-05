package go_chrome

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type goHttp struct {
	opt *GoChromeOptions
}

func getHttp(opt *GoChromeOptions) *goHttp {
	return &goHttp{opt: opt}
}

func (gh *goHttp) StartHttpServer() error {
	for k, v := range gh.opt.HttpRoute {
		if k != "/" {
			http.HandleFunc("/"+strings.TrimLeft(k, "/"), v)
		}
	}
	httpFileSystem := http.FileServer(http.Dir("./resources"))
	if gh.opt.AssetFile != nil {
		httpFileSystem = http.FileServer(gh.opt.AssetFile)
	}
	http.Handle("/", httpFileSystem)
	if gh.opt.HttpPort == 0 {
		gh.opt.HttpPort = 9527
	}
	return http.ListenAndServe(":"+fmt.Sprintf("%d", gh.opt.HttpPort), nil)
}

func (gh *goHttp) GetUrl(urlStr string, useHttpServer bool) string {
	if strings.HasPrefix(urlStr, "http") {
		return urlStr
	}
	if useHttpServer {
		port := gh.opt.HttpPort
		if gh.opt.HttpPort == 0 {
			port = 9527
		}
		return fmt.Sprintf("http://127.0.0.1:%d/", port)
	}
	emptyUrl := "data:text/html,<html></html>"
	if !useHttpServer && urlStr == "" {
		return emptyUrl
	}
	return "data:text/html," + url.PathEscape(urlStr)
}
