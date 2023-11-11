package jdns

import (
	"fmt"
	"github.com/chroblert/jlog"
	"strings"
	"testing"
)

func TestQueryA(t *testing.T) {
	resolve("gcch-mtis.cortana.ai.")
}

var (
	client_list = map[string]string{
		"US":          "102.128.164.0",
		"Afghanistan": "103.102.220.0",
		"Hong Kong":   "101.0.24.0",
		"China":       "101.101.64.0",
		"Belarus":     "134.17.0.0",
		"BR":          "24.152.44.0",
		"IN":          "5.154.174.0",
		"TW":          "27.147.0.0",
		"LA":          "43.252.244.0",
		"blank":       "",
	}
)

func TestQuery2A(t *testing.T) {
	domain := "www.qq.com."
	for city, client_ip := range client_list {
		domain_path, dest_ip, err := ResolveA("114.114.114.114:53", domain, client_ip)
		if err != nil {
			jlog.NError(city, err)
		} else {
			jlog.NInfo(city, fmt.Sprintf("%s ->", domain), strings.Join(domain_path, " -> "), "->", strings.Join(dest_ip, ","))
		}
		domain_cname_map, domain_a_map := GetDomainMap(domain, domain_path, dest_ip)
		//lastdomain := domain
		for k, v := range domain_cname_map {
			//lastdomain = v
			jlog.Info(k, "-->", v)
		}
		for k, v := range domain_a_map {
			for _, v2 := range v {
				jlog.Info(k, "-->", v2)
			}
		}
		jlog.Info()
	}

}

func BenchmarkQuery2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		domain := "www.qq.com."
		dns_server := "8.8.4.4:53"
		for city, client_ip := range client_list {
			domain_path, dest_ip, err := ResolveA(dns_server, domain, client_ip)
			if err != nil {
				jlog.NError(city, err)
			} else {
				jlog.NInfo(city, fmt.Sprintf("%s ->", domain), strings.Join(domain_path, " -> "), "->", strings.Join(dest_ip, ","))
			}
		}
	}
	jlog.Flush()
}
