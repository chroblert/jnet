package jip

import (
	"github.com/chroblert/jlog"
	"testing"
)

func TestIpRangeToCIDR(t *testing.T) {
	jlog.Info(IpRangeToCIDR("fec0::", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffe"))
	jlog.Info(IpRangeToCIDR("192.168.0.0", "192.168.0.255"))
	jlog.Info(IpRangeToCIDR("52.93.69.220", "52.93.69.220"))
	jlog.Info(IPInCIDRList("52.93.69.220", []string{"52.93.69.220/32"}))
}
