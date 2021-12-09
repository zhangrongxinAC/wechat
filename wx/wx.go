package wx

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/clbanning/mxj"
)

type weixinQuery struct {
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	EncryptType  string `json:"encrypt_type"`
	MsgSignature string `json:"msg_signature"`
	Echostr      string `json:"echostr"`
}

type WeixinClient struct {
	Token          string
	Query          weixinQuery // 请求的一些参数
	Message        map[string]interface{}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Methods        map[string]func() bool
}

/// 请求数据 Request
/// 返回数据 ResponseWriter
/// token 是自己的
func NewClient(r *http.Request, w http.ResponseWriter, token string) (*WeixinClient, error) {
	weixinClient := new(WeixinClient)
	weixinClient.Token = token // 获取本地的token
	weixinClient.Request = r
	weixinClient.ResponseWriter = w

	weixinClient.initWeixinQuery()
	log.Println("Signature:", weixinClient.Query.Signature)
	if weixinClient.Query.Signature != weixinClient.hashcode() { // 签名认证
		return nil, errors.New("Invalid Signature.")
	}

	return weixinClient, nil
}

func (this *WeixinClient) initWeixinQuery() {
	var q weixinQuery
	log.Println("URL:", this.Request.URL.Path, ", RawQuery:", this.Request.URL.RawPath)
	q.Nonce = this.Request.URL.Query().Get("nonce")
	q.Echostr = this.Request.URL.Query().Get("echostr")
	q.Signature = this.Request.URL.Query().Get("signature")
	q.Timestamp = this.Request.URL.Query().Get("timestamp")
	q.EncryptType = this.Request.URL.Query().Get("encrypt_type")
	q.MsgSignature = this.Request.URL.Query().Get("msg_signature")

	this.Query = q
}

func (this *WeixinClient) hashcode() string {
	strs := sort.StringSlice{this.Token, this.Query.Timestamp, this.Query.Nonce} // 使用本地的token生成校验
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *WeixinClient) initMessage() error {
	body, err := ioutil.ReadAll(this.Request.Body)

	if err != nil {
		return err
	}

	m, err := mxj.NewMapXml(body)

	if err != nil {
		return err
	}

	if _, ok := m["xml"]; !ok {
		return errors.New("Invalid Message.")
	}

	message, ok := m["xml"].(map[string]interface{})

	if !ok {
		return errors.New("Invalid Field `xml` Type.")
	}

	this.Message = message

	log.Println(this.Message)

	return nil
}

func (this *WeixinClient) text() {

	inMsg, ok := this.Message["Content"].(string) // 读取内容

	if !ok {
		return
	}

	var reply TextMessage

	reply.InitBaseData(this, "text")
	reply.Content = value2CDATA(fmt.Sprintf("我收到的是：%s", inMsg)) // 把消息再次封装

	replyXml, err := xml.Marshal(reply) // 序列化

	if err != nil {
		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	this.ResponseWriter.Header().Set("Content-Type", "text/xml") // 数据类型text/xml
	this.ResponseWriter.Write(replyXml)                          // 回复微信平台
}

func (this *WeixinClient) Run() {

	err := this.initMessage()

	if err != nil {

		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	MsgType, ok := this.Message["MsgType"].(string)

	if !ok {
		this.ResponseWriter.WriteHeader(403)
		return
	}

	switch MsgType {
	case "text":
		this.text() // 处理文本消息
		break
	default:
		break
	}

	return
}
