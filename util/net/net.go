package net

import (
	"fmt"
	"net"
	"time"
	"strings"
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

// GetLocalIPV4 获取本机IPV4
func GetLocalIPV4() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		// 排除掉虚拟网卡、未启用或回环接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if name := iface.Name; strings.HasPrefix(name, "docker0") || name == "lo" || strings.HasPrefix(name, "veth") || strings.HasPrefix(name, "br-") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ip4 := ipNet.IP.To4(); ip4 != nil {
					// 只返回常见内网网段
					if (ip4[0] == 10) || (ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || (ip4[0] == 192 && ip4[1] == 168) {
						return ip4.String(), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("no valid local IPv4 address found")
}
