package ip

import (
	"net"

	"github.com/acexy/golang-toolkit/logger"
	"github.com/yl2chen/cidranger"
)

// Checker matches IP addresses against mutable CIDR rules.
type Checker struct {
	ranger cidranger.Ranger
}

// NewIpChecker creates a CIDR checker, logging and skipping invalid rules.
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

// AddRuleIp adds CIDR rules, logging and skipping invalid input.
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

// RemoveRuleIp removes CIDR rules, logging parsing or removal errors.
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

// Match checks IP membership and returns any ranger error.
func (i *Checker) Match(ip string) (bool, error) {
	return i.ranger.Contains(net.ParseIP(ip))
}
