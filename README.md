# golang-toolkit

Golang 工具包，封装了缓存、HTTP 客户端、日志、加密、系统和常用工具函数，方便快速组装后端服务基础能力。

## 目录

- [安装](#安装)
- [快速开始](#快速开始)
    - [缓存](#缓存)
    - [HTTP 客户端](#http-客户端)
    - [日志](#日志)
    - [邮件](#邮件)
- [模块概览](#模块概览)
- [第三方依赖](#第三方依赖)

## 安装

```
go get -u github.com/acexy/golang-toolkit
```

## 快速开始

### 缓存

```go
import (
"time"

"github.com/acexy/golang-toolkit/caching"
)

bucket := caching.NewSimpleBigCache(5 * time.Minute)
manager := caching.NewCacheBucketManager("default", bucket)
if err := manager.Put("user:%d", userData, userID); err != nil {
// 处理错误
}
```

### HTTP 客户端

```go
import (
"github.com/acexy/golang-toolkit/httpclient"
)

client := httpclient.NewRestyClient()
var respData MyResponse
resp, err := client.R().SetReturnStruct(&respData).Get("https://api.example.com/status")
```

### 日志

```go
import "github.com/acexy/golang-toolkit/logger"

logger.EnableConsole(logger.InfoLevel)
logger.Logrus().Infoln("服务启动完成")
```

### 邮件

```go
import "github.com/acexy/golang-toolkit/email"

mail := email.NewGoMail("smtp.example.com", 465, "user", "pass", "noreply@example.com", true)
content := email.NewContent([]*email.ToAddress{{Address: "dev@example.com"}}, "测试邮件")
content.SetContent("text/plain", "Hello from golang-toolkit")
if err := mail.SendMail(content); err != nil {
// 处理发送异常
}
```

## 模块概览

### caching

对 `github.com/allegro/bigcache/v3` 的轻量封装，包括 `CacheManager`/`CacheBucket` 接口定义，配合可序列化的 `MemCacheKey`
便于统一管理多个缓存桶。

### httpclient

基于 Resty 的 `RestyClient` 包装，支持代理池、超时、自动下载、Request/Method 链式设置、JSON/表单、header、重定向控制等常用逻辑。

### logger

封装 `logrus`，提供可定制的控制台/文件输出、JSON/text formatter、trace id 追踪、等级判断、日志级别配置等工厂方法。

### email

基于 `gomail.v2`  的 `GoMail` 单例：支持发件人昵称、HTML/text 邮件正文、附件、批量收件人组装。

### crypto

- `symmetric`：AES 实现（CBC/GCM 模式、可插拔 IV/result/padding 组件）。
- `asymmetric`：RSA、ECDSA 的密钥生成、签名、验签以及 PEM 编码/解码。
- `hashing`：MD5/SHA256 简易封装。

### math

- `conversion`：多种基本类型与字节/字符串之间的转换（number/byte/binary/string 等）。
- `random`：整型、浮点、字符串随机值生成（支持范围指定）。

### sys

提供 CPU 数量探测、常用信号优雅关机、GOMAXPROCS 设置、goroutine 本地存储、TraceId 生成等系统级辅助函数。

### util

包含常用工具：

- `coll`：切片/映射帮助 (过滤、去重、随机、一致性判断)。
- `date`：时间格式化/解析、时间戳、区间计算、日历辅助等。
- `gob`：GOB 编解码快捷函数。
- `json`：JSON 序列化/反序列化、结构体复制、格式化。
- `net`：端口探活、常用网络信息。
- `reflect`：结构体字段枚举/赋值、深拷贝、非零值提取等。
- `str`：字符串操作、Builder、命名风格转换。

## 第三方依赖

| 库                             | 作用         |
|-------------------------------|------------|
| github.com/go-resty/resty/v2  | HTTP 客户端   |
| github.com/sirupsen/logrus    | 日志         |
| github.com/allegro/bigcache   | 内存缓存       |
| gopkg.in/gomail.v2            | 邮件发送       |
| github.com/jinzhu/copier      | 结构体复制      |
| github.com/timandy/routine    | 协程本地存储     |
| github.com/tidwall/gjson      | JSON 解析    |
| github.com/google/uuid        | UUID       |
| github.com/iancoleman/strcase | 字符串大小写转换   |
| github.com/yl2chen/cidranger  | IP/CIDR 匹配 |

