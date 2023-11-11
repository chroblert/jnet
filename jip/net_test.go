package jip

import (
	"github.com/chroblert/jlog"
	"testing"
)

func TestRebaseCIDR(t *testing.T) {
	subnet_arr, _ := RebaseCIDR("52.123.128.0/15", 14)
	for _, v := range subnet_arr {
		jlog.NInfo(v)
	}
}
