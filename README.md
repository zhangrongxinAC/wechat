### 基本配置
1. 注册/登录公众号
2. 在公众号中左栏 [设置与开发]/[基本配置] 中设置自己的服务器配置（需要服务器启动验证）
3. token 需要与服务器的 token 设置为一样

### 注册公众号/订阅号
```
https://mp.weixin.qq.com/cgi-bin/registermidpage?action=index&lang=zh_CN&token=
```

### 原理
微信后台向服务器发送 {token, timestamp, nonce} 的三元组进行验证
timestamp: 发送的时间戳, nonce: 随机数
验证成功后, 服务器解析 xml 数据包

### 功能点
1. 验证 token（与微信后台服务器）
2. 解析微信后台发送的 xml 消息

### 运行代码
```
go mod init wechat
go build
sudo ./wechat
```
运行成功后, 在第一步基本配置中点击提交即可