package net

import (
	"net"
	"time"
)

// Telnet 尝试在指定的时间内连接指定主机的指定端口 用于探活网络端口 address: ip:port
func Telnet(address string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
