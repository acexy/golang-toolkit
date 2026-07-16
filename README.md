# golang-toolkit

`golang-toolkit` 是一个面向 Go 后端项目的基础工具库，集中封装缓存、HTTP 客户端、日志、邮件、加密、系统辅助能力以及常用工具函数。

## 环境要求

- Go `>= 1.25.0`
- Module：`github.com/acexy/golang-toolkit`

## 安装

```bash
go get github.com/acexy/golang-toolkit
```

## 重点能力

- **统一错误定义**：公共错误集中放在 `error` 包，业务包尽量返回可直接比较的独立错误变量，便于 `errors.Is` 判断和跨模块复用。
- **HTTP 客户端封装**：`httpclient` 基于 Resty，提供请求构造、JSON body、query/path 参数、代理、多代理随机选择、TLS 配置、下载文件和响应绑定等能力。
- **邮件发送**：`email` 基于 `github.com/wneessen/go-mail`，支持发件人名称、收件人显示名、HTML/text 正文、附件和真实 SMTP 发送测试。
- **缓存管理**：`caching` 基于 BigCache，支持多 bucket 管理、类型化 key、缓存编解码和统一 `CacheManager`。
- **日志封装**：`logger` 基于 logrus，支持控制台、文件、日志级别、自定义 formatter、trace id 和滚动文件输出。
- **加密与摘要**：`crypto` 提供 AES 对称加密、RSA/ECDSA 非对称能力，以及 MD5、SHA256 等摘要函数。
- **集合工具**：`util/coll` 提供 slice/map 的常用操作，包括遍历、查找、过滤、映射、去重、分组、合并、随机选择和二维切片展开。
- **JSON 工具**：`util/json` 提供 JSON 序列化/反序列化、结构体复制、gjson 快速读取、时间戳包装类型和全局一次性时间戳配置。
- **数学与随机工具**：`math` 提供数字/字节转换、十六进制解析、随机数、随机字符串和概率选择。
- **系统辅助能力**：`sys` 和 `sys/routine` 提供 CPU/GOMAXPROCS、优雅退出和 goroutine 本地上下文等工具。

## 包索引

| 包 | 说明 |
| --- | --- |
| `caching` | BigCache 封装、多 bucket、缓存 key、Codec |
| `crypto/asymmetric` | RSA、ECDSA 等非对称加密/签名能力 |
| `crypto/hashing` | MD5、SHA256 等摘要工具 |
| `crypto/symmetric` | AES 等对称加密能力 |
| `email` | SMTP 邮件发送、正文、附件、地址封装 |
| `error` | 项目公共错误变量 |
| `httpclient` | Resty 客户端封装 |
| `logger` | logrus 日志封装 |
| `math` | 字节、数值、随机和概率相关工具 |
| `math/conversion` | 字符串、数字、字节转换 |
| `math/random` | 随机数、随机字符串、概率选择 |
| `sys` | 系统信号、CPU、进程辅助工具 |
| `sys/routine` | goroutine 本地存储 |
| `util/coll` | slice/map 集合工具 |
| `util/date` | 日期时间工具 |
| `util/gob` | GOB 序列化工具 |
| `util/json` | JSON、gjson、时间戳包装工具 |
| `util/net` | 网络、IP 辅助工具 |
| `util/reflect` | 反射辅助工具 |
| `util/str` | 字符串辅助工具 |

## 使用约定

- 错误优先引用 `error` 包中已有的大类错误，避免为明显同类场景创建过细错误。
- JSON 时间戳全局配置只应初始化一次；需要临时指定秒/毫秒格式时，使用显式 `WithType` 转换函数。
- 新增工具函数应放在最小相关包中，命名保持 Go 风格，例如 `JSON`、`URL` 等常见缩写使用全大写。
- 依赖由 Go module 管理；只有调整依赖时才需要运行 `go mod tidy`。

## 测试

```bash
go test ./...
```

针对单个包迭代时可以执行：

```bash
go test ./util/json
go test ./httpclient
go test ./email
```
