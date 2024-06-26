package net

import (
	"fmt"
	"net"
	"time"
)

// Telnet 尝试在指定的时间内连接指定主机的指定端口 用于探活网络端口
func Telnet(host string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout) // 2秒超时
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
