package ip

import (
	"net"

	"github.com/acexy/golang-toolkit/logger"
	"github.com/yl2chen/cidranger"
)

type Checker struct {
	ranger cidranger.Ranger
}

func NewIpChecker(cidrIps ...string) *Checker {
	ranger := cidranger.NewPCTrieRanger()
	for _, ip := range cidrIps {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			logger.Logrus().Warningln("parse ip error. ip:", ip)
			continue
		}
		_ = ranger.Insert(cidranger.NewBasicRangerEntry(*ipNet))
	}
	return &Checker{
		ranger: ranger,
	}
}

// AddRuleIp 添加新的ip规则
func (i *Checker) AddRuleIp(cidrIps ...string) {
	if len(cidrIps) == 0 {
		return
	}
	for _, ip := range cidrIps {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			logger.Logrus().Warningln("parse ip error. ip:", ip)
			continue
		}
		_ = i.ranger.Insert(cidranger.NewBasicRangerEntry(*ipNet))
	}
}

// RemoveRuleIp 删除ip规则
func (i *Checker) RemoveRuleIp(cidrIps ...string) {
	for _, ip := range cidrIps {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			logger.Logrus().Warningln("parse ip error. ip:", ip)
			continue
		}
		_, err = i.ranger.Remove(*ipNet)
		if err != nil {
			logger.Logrus().Warningln("remove ip error. ip:", ip)
		}
	}
}

// Match 检测ip是否匹配
func (i *Checker) Match(ip string) (bool, error) {
	return i.ranger.Contains(net.ParseIP(ip))
}
