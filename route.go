package main

import (
	"io"
	"net/http"
	"regexp"
	"time"
)

type WebController struct {
	Function func(http.ResponseWriter, *http.Request)
	Method   string
	Pattern  string
}

var mux []WebController // 自己定义的路由
// ^ 匹配输入字符串的开始位置
func init() {
	mux = append(mux, WebController{post, "POST", "^/"})
	mux = append(mux, WebController{get, "GET", "^/"})
}

type httpHandler struct{} // 实际是实现了Handler interface
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t := time.Now()

	for _, webController := range mux { // 遍历路由
		// 匹配请求的   r.URL.Path  -> webController.Pattern
		if m, _ := regexp.MatchString(webController.Pattern, r.URL.Path); m { // 匹配URL

			if r.Method == webController.Method { // 匹配方法

				webController.Function(w, r) // 调用对应的处理函数

				go writeLog(r, t, "match", webController.Pattern)

				return
			}
		}
	}

	go writeLog(r, t, "unmatch", "")

	io.WriteString(w, "")
}
