package go_chrome

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
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
		http.HandleFunc("/"+strings.TrimLeft(k, "/"), v)
	}
	if _, ok := gh.opt.HttpRoute["/"]; !ok {
		http.HandleFunc("/", gh.indexHandler)
	}
	if gh.opt.HttpPort == 0 {
		gh.opt.HttpPort = 9527
	}
	return http.ListenAndServe(":"+fmt.Sprintf("%d", gh.opt.HttpPort), nil)
}
func (gh *goHttp) indexHandler(w http.ResponseWriter, r *http.Request) {
	gh.View(w, r, "index", gh.opt.DefHttpIndexData)
}

func (gh *goHttp) View(w http.ResponseWriter, r *http.Request, temp string, param interface{}, suffix ...string) {
	suffixStr := "html"
	if len(suffix) > 0 {
		suffixStr = suffix[0]
	}
	tempPath := fmt.Sprintf("./resources/%s.%s", temp, suffixStr)
	t, err := template.ParseFiles(tempPath)
	if err != nil {
		return
	}
	_ = t.Execute(w, param)
}

func GetUrl(urlStr string, useHttpServer bool) string {
	if strings.HasPrefix(urlStr, "http") {
		return urlStr
	}
	if (urlStr == "" || urlStr == "index") && useHttpServer {
		return "http://127.0.0.1/"
	}
	emptyUrl := "data:text/html,<html></html>"
	if !useHttpServer && urlStr == "" {
		return emptyUrl
	}
	content, err := os.ReadFile("./resources/" + urlStr + ".html")
	if err != nil {
		return emptyUrl
	}
	return "data:text/html," + url.PathEscape(string(content))
}
