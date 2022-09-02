package go_chrome

import (
	"fmt"
	"html/template"
	"net/http"
)

type goHttp struct {
	opt *GoChromeOptions
}

func getHttp(opt *GoChromeOptions) goHttp {
	return goHttp{opt: opt}
}

func (gh *goHttp) StartHttpServer() error {
	for k, v := range gh.opt.HttpRoute {
		http.HandleFunc(k, v)
	}
	if _, ok := gh.opt.HttpRoute["/"]; !ok {
		http.HandleFunc("/", gh.indexHandler)
	}
	return http.ListenAndServe("127.0.0.0:"+fmt.Sprintf("%d", gh.opt.HttpPort), nil)
}
func (gh *goHttp) indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		return
	}
	_ = t.Execute(w, gh.opt.DefHttpIndexData)
}
