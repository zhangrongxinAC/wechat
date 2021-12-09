package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var LogFile *os.File

func init() {
	// fmt.Println("log init")
	// LogFile, err := os.OpenFile("wechat.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //打开日志文件，不存在则创建
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// log.SetOutput(LogFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func CloseLog() {
	if LogFile != nil {
		LogFile.Close()
	}

}

func writeLog(r *http.Request, t time.Time, match string, pattern string) {

	if logLevel != "prod" {

		d := time.Now().Sub(t)

		l := fmt.Sprintf("[ACCESS] | % -10s | % -40s | % -16s | % -10s | % -40s |", r.Method, r.URL.Path, d.String(), match, pattern)

		log.Println(l)
	}
}

func func_log2fileAndStdout() {
	//创建日志文件
	f, err := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//完成后，延迟关闭
	defer f.Close()
	// 设置日志输出到文件
	// 定义多个写入器
	writers := []io.Writer{
		f,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	// 创建新的log对象
	logger := log.New(fileAndStdoutWriter, "", log.Ldate|log.Ltime|log.Lshortfile)
	// 使用新的log对象，写入日志内容
	logger.Println("--> logger :  check to make sure it works")
}
