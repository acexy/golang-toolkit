package ip

import (
	"fmt"
	"testing"
)

func TestIpChecker(t *testing.T) {
	ipChecker := NewIpChecker("91.231.222.0/24", "192.168.1.2/32")
	fmt.Println(ipChecker.Match("91.231.222.173"))
	fmt.Println(ipChecker.Match("192.168.1.3"))

	ipChecker.AddRuleIp("192.168.1.3/32")
	ipChecker.AddRuleIp("192.168.1.3/32")

	fmt.Println(ipChecker.Match("192.168.1.3"))
	ipChecker.RemoveRuleIp("192.168.1.3/32")
	fmt.Println(ipChecker.Match("192.168.1.3"))
}
