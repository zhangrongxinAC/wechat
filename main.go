package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"wechat/wx"
)

const (
	logLevel = "dev"
	port     = 80
	token    = "" // 生成地址：https://suijimimashengcheng.51240.com/
	/* token 需要与微信服务器配置一致 */
)

// 处理token的认证
func get(w http.ResponseWriter, r *http.Request) {

	client, err := wx.NewClient(r, w, token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403) // 校验失败
		return
	}

	if len(client.Query.Echostr) > 0 {
		w.Write([]byte(client.Query.Echostr)) // 校验成功返回的是Echostr
		return
	}

	w.WriteHeader(403)
}

// 微信平台过来消息， 处理 ，然后返回微信平台
func post(w http.ResponseWriter, r *http.Request) {

	client, err := wx.NewClient(r, w, token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}
	// 到这一步签名已经验证通过了
	client.Run()
}

// 编译方法
// go mod init wechat
// go build
// ./wechat
// 需要自己修改token，以适应自己公众号的token
func main() {
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", port), // 设置监听地址， ip:port
		Handler:        &httpHandler{},           // 用什么handler来处理
		ReadTimeout:    5 * time.Second,          // 读写超时 微信给出来5
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	log.Println(fmt.Sprintf("Listen: %d", port))
	log.Fatal(server.ListenAndServe())
	defer CloseLog()
}
