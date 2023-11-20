package jip

import (
	"fmt"
	"github.com/chroblert/jlog"
	"net"
	"testing"
)

func TestRebaseCIDR(t *testing.T) {
	subnet_arr, _ := RebaseCIDR("52.123.128.0/15", 14)
	for _, v := range subnet_arr {
		jlog.NInfo(v)
	}
}

func TestCIDRContainsIP(t *testing.T) {
	jlog.Info(CIDRContainsIP("2001:0db8::/32", "2001:0db8::"))
	jlog.Info(CIDRContainsIP("2001:0db8::/32", "2001:0db9::"))

	ip, ipnet, err := net.ParseCIDR("2001:0db8::/32")
	if err != nil {
		return
	}
	jlog.Info(ip, ipnet)
}

func TestParseIPStr(t *testing.T) {
	ip_list := ParseIPStr(fmt.Sprintf("%s-%s", "1.1.1.1", "1.2.2.0"))

	for k, v := range ip_list {
		jlog.Info(k, v)
	}
}
