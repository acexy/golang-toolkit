package net

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Telnet probes a TCP address within timeout and closes successful connections; it does not test application health.
func Telnet(address string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// GetLocalIPV4 returns the first private IPv4 on an eligible active interface, or an error if none is found.
func GetLocalIPV4() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		// Skip inactive, loopback, and selected virtual interfaces.
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
					// Select only common private IPv4 ranges.
					if (ip4[0] == 10) || (ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || (ip4[0] == 192 && ip4[1] == 168) {
						return ip4.String(), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("no valid local IPv4 address found")
}
