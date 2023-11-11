package jdomain

import "strings"

func Get3ldDomain(domain string) (sld_domain string) {
	return GetSpecificLayerDomain(domain, 3)
}

func GetSldDomain(domain string) (sld_domain string) {
	return GetSpecificLayerDomain(domain, 2)
}

func GetTldDomain(domain string) (tld_domain string) {
	return GetSpecificLayerDomain(domain, 1)
}

// ww.test.com为三级
// t4.www.test.com为四级
// 若不足layer，则返回空字符串
func GetSpecificLayerDomain(domain string, layer int) (domain2 string) {
	domain = strings.TrimRight(domain, ".")
	domain_list := strings.Split(domain, ".")
	if len(domain_list) < layer {
		return ""
	}
	s_domain := strings.Join(domain_list[len(domain_list)-layer:], ".")
	return s_domain
}
