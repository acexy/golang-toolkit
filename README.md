# golang-toolkit

`golang-toolkit` is a foundational utility library for Go backend projects. It provides reusable abstractions for caching, HTTP clients, logging, email, cryptography, system utilities, and common helper functions.

## Requirements

- Go `>= 1.25.0`
- Module: `github.com/acexy/golang-toolkit`

## Installation

```bash
go get github.com/acexy/golang-toolkit
```

## Key Features

- **Centralized error definitions**: Common errors are defined in the `error` package. Feature packages return independently comparable error variables whenever possible, making `errors.Is` checks and cross-package reuse straightforward.
- **HTTP client abstraction**: Built on Resty, `httpclient` supports request construction, JSON bodies, query and path parameters, proxies, randomized proxy selection, TLS configuration, file downloads, and response binding.
- **Email delivery**: Built on `github.com/wneessen/go-mail`, `email` supports sender names, recipient display names, HTML and text bodies, attachments, and integration tests against a real SMTP server.
- **Cache management**: Built on BigCache, `caching` supports multiple buckets, typed keys, cache encoding and decoding, and a unified `CacheManager`.
- **Logging abstraction**: Built on logrus, `logger` supports console and file output, log levels, custom formatters, trace IDs, and log rotation.
- **Cryptography and hashing**: `crypto` provides AES symmetric encryption, RSA and ECDSA asymmetric operations, and digest functions such as MD5 and SHA-256.
- **Collection utilities**: `util/coll` provides common slice and map operations, including iteration, lookup, filtering, mapping, deduplication, grouping, merging, random selection, and two-dimensional slice flattening.
- **JSON utilities**: `util/json` provides JSON serialization and deserialization, struct copying, fast reads with gjson, timestamp wrapper types, and one-time global timestamp configuration.
- **Math and random utilities**: `math` provides numeric and byte conversions, hexadecimal parsing, random numbers, random strings, and probability-based selection.
- **System utilities**: `sys` and `sys/routine` provide CPU and GOMAXPROCS helpers, graceful shutdown handling, and goroutine-local context.

## Package Index

| Package | Description |
| --- | --- |
| `caching` | BigCache abstraction, multiple buckets, cache keys, and codecs |
| `crypto/asymmetric` | RSA, ECDSA, and other asymmetric encryption and signature operations |
| `crypto/hashing` | MD5, SHA-256, and other digest utilities |
| `crypto/symmetric` | AES and other symmetric encryption operations |
| `email` | SMTP delivery, message bodies, attachments, and address wrappers |
| `error` | Shared project error variables |
| `httpclient` | Resty client abstraction |
| `logger` | logrus logging abstraction |
| `math` | Byte, numeric, random, and probability utilities |
| `math/conversion` | String, numeric, and byte conversions |
| `math/random` | Random numbers, random strings, and probability-based selection |
| `sys` | System signal, CPU, and process utilities |
| `sys/routine` | Goroutine-local storage |
| `util/coll` | Slice and map collection utilities |
| `util/date` | Date and time utilities |
| `util/gob` | GOB serialization utilities |
| `util/json` | JSON, gjson, and timestamp wrapper utilities |
| `util/net` | Network and IP utilities |
| `util/reflect` | Reflection utilities |
| `util/str` | String utilities |

## Usage Conventions

- Prefer existing broad error definitions from the `error` package instead of creating overly specific errors for clearly related cases.
- Initialize the global JSON timestamp configuration only once. Use explicit `WithType` conversion functions when a temporary seconds or milliseconds format is required.
- Add new utility functions to the smallest relevant package and follow Go naming conventions, including uppercase common initialisms such as `JSON` and `URL`.
- Dependencies are managed with Go modules. Run `go mod tidy` only when dependencies change.

## Testing

```bash
go test ./...
```

To iterate on a single package, run:

```bash
go test ./util/json
go test ./httpclient
go test ./email
```
