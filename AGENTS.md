# Repository Guidelines

## Project Structure & Module Organization

This repository is a Go toolkit module: `github.com/acexy/golang-toolkit`.
Packages are organized by capability at the repository root. Examples include
`caching/`, `httpclient/`, `logger/`, `email/`, and `sys/`. Broader utility
areas use nested packages, such as `crypto/asymmetric`, `crypto/symmetric`,
`math/conversion`, `util/json`, `util/net`, and `util/str`.

Tests live beside the package they validate and use Go's standard
`*_test.go` naming convention, for example `caching/caching_test.go` and
`logger/logrus_test.go`. Keep new code in the smallest relevant package and
avoid cross-package changes unless the public API requires them.

## Build, Test, and Development Commands

- `go test ./...`: run the full test suite for all packages.
- `go test ./caching`: run tests for one package while iterating.
- `go test -cover ./...`: run tests with coverage reporting.
- `go test ./... -run TestName`: run a focused test by name.
- `go mod tidy`: update module requirements after dependency or import changes.

Use the shared Maven, Go, and npm caches configured for the workspace when
tooling needs cache paths. Do not introduce new dependency caches.

## Coding Style & Naming Conventions

Follow idiomatic Go style and keep files formatted with `gofmt`. Use tabs for
Go indentation as produced by the formatter. Package names should be short,
lowercase, and descriptive, such as `caching` or `httpclient`.

Exported identifiers must be documented when they form part of the package API.
Prefer clear constructor names like `NewRestyClient` and focused helper names
that describe behavior. Keep comments concise and useful for maintenance.

## Testing Guidelines

Use Go's standard `testing` package. Name tests as `TestXxx` and place them in
the same package directory as the implementation. Add or update tests for new
public behavior, bug fixes, and edge cases around serialization, networking,
logging, crypto, and time-sensitive utilities.

## Commit & Pull Request Guidelines

Recent history uses concise conventional-style prefixes, for example
`feat: add slice flat helper` and `bugfix: logger输出的日志文件丢失调用细节`.
Use a short imperative subject and include the affected package when helpful.

Pull requests should explain the change, list impacted packages, mention any
API compatibility concerns, and include the test command used. Link related
issues when available. Screenshots are generally unnecessary unless a change
adds visual documentation or generated assets.

## Security & Configuration Tips

Never commit secrets, credentials, private keys, `.env` files, or generated log
files. Treat `logger/logs/` as runtime output, not source. Avoid changing
dependency versions or module files unless the task explicitly requires it.
