package jip

import (
	"fmt"
	"github.com/c-robinson/iplib"
	"github.com/chroblert/jlog"
	"net"
	"strings"
)

// 比较原始CIDR的网络掩码数量，与，新掩码位数
// 比较的是原始的主机数量是否大于新掩码的主机数量
// /15 > /24
func CompareCIDR(originalCIDR string, maskLen int) int {
	// 解析原始CIDR
	_, ipNet, err := net.ParseCIDR(originalCIDR)
	if err != nil {
		fmt.Println("解析原始CIDR时发生错误:", err)
		return -2
	}

	// 获取原始CIDR的IP地址和子网掩码
	//ip := ipNet.IP
	originalMask := ipNet.Mask
	//originalLastNetNum := 0
	originalMaskLen, _ := originalMask.Size()
	if originalMaskLen < maskLen {
		return 1
	} else if originalMaskLen == maskLen {
		return 0
	} else {
		return -1
	}
}
func RebaseCIDR(originalCIDR string, newMaskLen int) ([]string, error) {
	// 解析原始CIDR
	pIp, ipNet, err := net.ParseCIDR(originalCIDR)
	if err != nil {
		fmt.Println("解析原始CIDR时发生错误:", err)
		return nil, nil
	}
	originalMask := ipNet.Mask
	//originalLastNetNum := 0
	originalMaskLen := originalMaskSize(originalMask)
	n := iplib.NewNet4(pIp, originalMaskLen)
	sub, _ := n.Subnet(newMaskLen)
	new_cidr_list := make([]string, len(sub))
	for k, v := range sub {
		new_cidr_list[k] = v.String()
	}

	return new_cidr_list, nil
}

func rebaseCIDR2(originalCIDR string, newMaskLen int) ([]string, error) {
	// 定义原始CIDR和要划分的新子网掩码长度
	//originalCIDR := originalCIDR
	newMaskLength := newMaskLen

	// 解析原始CIDR
	pIp, ipNet, err := net.ParseCIDR(originalCIDR)
	if err != nil {
		fmt.Println("解析原始CIDR时发生错误:", err)
		return nil, nil
	}

	// 获取原始CIDR的IP地址和子网掩码
	ip := ipNet.IP
	originalMask := ipNet.Mask
	//originalLastNetNum := 0
	originalMaskLen := originalMaskSize(originalMask)
	// 网络掩码到第多少个字节
	maskOrder := -1
	if (originalMaskLen % 8) != 0 {
		maskOrder = (originalMaskLen / 8) // 向下取整
	} else {
		maskOrder = (originalMaskLen / 8)
	}
	//initial_mask_order := maskOrder
	// 最后网络掩码所在字节的掩码
	// 0，1，2，3
	lastMaskOrderBit := originalMask[maskOrder]
	// 最后网络掩码所在字节的大小
	realOrderIPValue := (pIp.To4())[maskOrder]
	jlog.Info(lastMaskOrderBit, realOrderIPValue)
	jlog.Info(lastMaskOrderBit & realOrderIPValue)
	// 最后网络掩码所在字节的大小与网络掩码与后的实际大小
	baseSubnetValue := lastMaskOrderBit & realOrderIPValue
	jlog.Info("基础IP：", baseSubnetValue)
	// 算出差多少
	// 新子网掩码与源子网掩码相差多少位
	total_diff_bit := newMaskLength - originalMaskLen
	first_diff_bit := 0   // 开头相差多少个bit
	middle_diff_byte := 0 // 中间相差多少个byte
	//last_diff_bit := 0    // 最后相差多少个bit
	if originalMaskLen%8 != 0 {
		// 据补满，差多少位
		first_diff_bit = 8 - originalMaskLen%8
	} else {
		first_diff_bit = 0
	}
	if first_diff_bit >= total_diff_bit {
		first_diff_bit = total_diff_bit
		middle_diff_byte = 0
		//last_diff_bit = 0
	} else {
		// 中间有多少个完整的
		middle_diff_byte = (total_diff_bit - first_diff_bit) / 8
		// 最后还有多少个,1000,0000(2^7=(2^8-1)-(2^(8-last_diff_bit))-1)),11000000,
		//last_diff_bit = total_diff_bit - 8*middle_diff_byte
	}
	jlog.Info(middle_diff_byte)
	for i := 0; i < 1<<first_diff_bit; i++ {

	}

	// 计算新子网的数量
	subnetCount := 1 << total_diff_bit
	jlog.Debug(uint(newMaskLength)-uint(originalMaskSize(originalMask)), "共", subnetCount, "个子网")
	// 循环生成新子网
	new_cidr_list := []string{}
	for i := 0; i < subnetCount; i++ {
		// 计算每个新子网的网络地址
		newSubnetIP := net.IP(make([]byte, len(ip)))
		copy(newSubnetIP, ip)
		//jlog.Info(newSubnetIP)
		newSubnetIP[maskOrder] = byte((int(baseSubnetValue) + i) % 256)
		//jlog.Info(newSubnetIP)
		// 获取新子网的子网掩码
		newSubnetMask := net.CIDRMask(newMaskLength, 32)
		// 判断是否已经超过255
		if (int(baseSubnetValue)+i)%255 == 0 {
			//baseSubnetValue = 0
			ip[maskOrder] = 255
			maskOrder += 1
		}
		// 计算下一个新子网的IP
		//incrementIP(newSubnetIP, newSubnetMask)

		// 打印新子网的CIDR
		newSubnetCIDR := &net.IPNet{
			IP:   newSubnetIP,
			Mask: newSubnetMask,
		}
		new_cidr_list = append(new_cidr_list, newSubnetCIDR.String())
		//fmt.Printf("子网 %d: %s\n", i+1, newSubnetCIDR.String())
	}
	return new_cidr_list, nil
}

// 计算子网掩码位数
func originalMaskSize(mask net.IPMask) int {
	size, _ := mask.Size()
	return size
}

// 增加IP地址
func incrementIP(ip net.IP, mask net.IPMask) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

// 1100 0000
// num_bits: 有多少位
// value_bits: 前多少位进行组合
func getLastNum(num_bits int, value_bit int) []int {
	tmp_num := -1
	num_arr := []int{}
	// 多少种组合
	max_value := (1 << value_bit) - 1
	// 差多少个bit
	diff_bit := num_bits - value_bit
	for i := 0; i <= max_value; i++ {
		tmp_num = i << diff_bit
		num_arr = append(num_arr, tmp_num)
	}
	return num_arr

}

// 判断ip是否是IPv6
func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

// 判断CIDR是否是IPv6
func CIDRIsIPv6(cidr string) bool {
	pip, _, _ := net.ParseCIDR(cidr)
	ip := net.ParseIP(pip.String())
	return ip != nil && strings.Contains(pip.String(), ":")
}

// 判断某CIDR是否包含某IP
// 1: 包含，0：不包含；-1：出错
func CIDRContainsIP(CIDR, ip_str string) int {
	ip := net.ParseIP(ip_str)
	if ip == nil {
		return -1
	}
	_, ipnet, err := net.ParseCIDR(CIDR)
	if err != nil {
		return -1
	}
	if ipnet.Contains(ip) {
		return 1
	} else {
		return 0
	}
}

func IPInCIDRList(ip_str string, cidr_list []string) int {
	for _, cidr := range cidr_list {
		if CIDRContainsIP(cidr, ip_str) == 1 {
			return 1
		}
	}
	return 0
}
