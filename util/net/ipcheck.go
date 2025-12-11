package net

import (
	"net"

	"github.com/acexy/golang-toolkit/logger"
	"github.com/yl2chen/cidranger"
)

type IpChecker struct {
	ranger cidranger.Ranger
}

func NewIpChecker(ruleIps ...string) *IpChecker {
	ranger := cidranger.NewPCTrieRanger()
	for _, ip := range ruleIps {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			logger.Logrus().Warningln("parse ip error. ip:", ip)
			continue
		}
		_ = ranger.Insert(cidranger.NewBasicRangerEntry(*ipNet))
	}
	return &IpChecker{
		ranger: ranger,
	}
}

// NewRuleIp 添加新的ip规则
func (i *IpChecker) NewRuleIp(ruleIps ...string) {
	if len(ruleIps) == 0 {
		return
	}
	for _, ip := range ruleIps {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			logger.Logrus().Warningln("parse ip error. ip:", ip)
			continue
		}
		_ = i.ranger.Insert(cidranger.NewBasicRangerEntry(*ipNet))
	}
}

// RemoveRuleIp 删除ip规则
func (i *IpChecker) RemoveRuleIp(ruleIps ...string) {
	for _, ip := range ruleIps {
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
func (i *IpChecker) Match(ip string) (bool, error) {
	return i.ranger.Contains(net.ParseIP(ip))
}
