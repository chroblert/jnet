package jip

// 获取给定ip列表中不在指定的CIDR列表中的ip
func GetIPListNotInCIDRList(cidr_list []string, ip_list []string) []string {
	not_ip_list := []string{}
	for _, ip_str := range ip_list {
		if IPInCIDRList(ip_str, cidr_list) == 0 {
			not_ip_list = append(not_ip_list, ip_str)
		}
	}
	return not_ip_list
}
