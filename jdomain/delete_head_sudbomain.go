package jdomain

import "strings"

// 删除给定域名前面的几个字符串
// order：从1开始
// return：字符串
func DeleteHeadNumDomain(s_domain string, order int) (d_domain string) {
	word_list := strings.Split(s_domain, ".")
	if len(word_list) < order+1 {
		return ""
	}
	last_domains := word_list[order:]
	subdomain := strings.Join(last_domains, ".")
	return subdomain
}
