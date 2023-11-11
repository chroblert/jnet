package jip

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ip1,ip2,ip3-ip4,ip5/cidr
// 将ip1,ip2,ip3-ip4,ip5/cidr转换成ip字符串列表
func ParseIPStr(ipStr string) []net.IP {
	tmpList := strings.Split(ipStr, ",")
	ipStrList := make([]net.IP, 0)
	for _, v := range tmpList {
		// 是不是CIDR
		if ips, err := GetIPSFromCIDR(v); err == nil {
			ipStrList = append(ipStrList, ips...)
		} else if ips, err := getIPSFromIPRange(v); err == nil {
			ipStrList = append(ipStrList, ips...)
		} else if ip := net.ParseIP(v); ip != nil {
			ipStrList = append(ipStrList, ip)
		} else {
			//jlog.Error(err)
		}
	}
	// 清空切片
	tmpList = tmpList[:0]
	//sort.Strings(ipStrList)
	return removeDuplicateElement(ipStrList)
}

func removeDuplicateElement(addrs []net.IP) []net.IP {
	result := make([]net.IP, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item.String()]; !ok {
			temp[item.String()] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func getIPSFromIPRange(ipRange string) ([]net.IP, error) {
	tmpList := strings.Split(ipRange, "-")
	if len(tmpList) != 2 {
		return nil, fmt.Errorf("不符合格式。ip1-ip2")
	}
	if ip2Int(tmpList[0]) > ip2Int(tmpList[1]) {
		return nil, fmt.Errorf("ip1应小于ip2")
	}
	ips := make([]net.IP, 0)
	startIP := net.ParseIP(tmpList[0])
	endIP := net.ParseIP(tmpList[1])
	// 过滤最后一位为0和255的IP
	for ip := startIP; ip2Int(ip.String()) <= ip2Int(endIP.String()); {
		ip = getNextIP(ip)
		ips = append(ips, ip)
	}
	return ips, nil
}

func ip2Int(ip string) int64 {
	if len(ip) == 0 {
		return 0
	}
	bits := strings.Split(ip, ".")
	if len(bits) < 4 {
		return 0
	}
	b0 := string2Int(bits[0])
	b1 := string2Int(bits[1])
	b2 := string2Int(bits[2])
	b3 := string2Int(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func string2Int(in string) (out int) {
	out, _ = strconv.Atoi(in)
	return
}

func GetIPSFromCIDR(cidr string) (ip_list []net.IP, err error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		//fmt.Println("CIDR解析失败:", err)
		return
	}
	ip := ipNet.IP.Mask(ipNet.Mask)
	for ipNet.Contains(ip) {
		//fmt.Println("IP地址:", ip)
		ip_list = append(ip_list, ip)
		ip = getNextIP(ip)
	}
	return
}

func getNextIP(ip net.IP) net.IP {
	nextIP := make(net.IP, len(ip))
	copy(nextIP, ip)

	for i := len(ip) - 1; i >= 0; i-- {
		nextIP[i]++
		if nextIP[i] > 0 {
			break
		}
	}

	return nextIP
}
