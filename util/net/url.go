package net

import "net/url"

// UrlEncode url 编码
func UrlEncode(raw string) string {
	return url.QueryEscape(raw)
}

// UrlDecode url 解码
func UrlDecode(raw string) (string, error) {
	return url.QueryUnescape(raw)
}
