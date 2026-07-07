# golang-toolkit

`golang-toolkit` 是一个 Go 后端工具包，封装缓存、HTTP 客户端、日志、邮件、加密、系统能力和常用工具函数。

go version >=1.25.0

## 安装

```bash
go get github.com/acexy/golang-toolkit
```

## 模块概览

- `caching`：基于 BigCache 的缓存封装，支持 `CacheManager`、多 bucket、类型化 `CacheKey`/`BucketName` 和自定义 `Codec`。
- `httpclient`：基于 Resty 的 HTTP 客户端封装，支持超时、代理、请求构造、返回值绑定和下载等能力。
- `logger`：基于 logrus 的日志封装，支持控制台、文件、格式化器、日志级别和 trace id。
- `email`：基于 go-mail 的邮件发送封装，支持 HTML/text 正文、附件和批量收件人。
- `crypto`：提供 AES、RSA、ECDSA、MD5、SHA256 等加密和摘要能力。
- `math`：提供基础类型转换、字节转换、随机数和随机字符串工具。
- `sys`：提供 CPU 探测、GOMAXPROCS 设置、优雅关机和 goroutine 本地存储等系统辅助能力。
- `util`：包含集合、时间、JSON、GOB、网络、反射和字符串等通用工具。
- `error`：集中定义项目公共错误变量，便于跨模块统一判断。

本项目使用 Go module 管理依赖。修改依赖后再运行 `go mod tidy`。
