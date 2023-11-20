package jip

import (
	"fmt"
	"math/rand"
	"net"
)

// 从cidr中随机取出指定数量的IP
func GetRandomIPFromCidr(cidr string, random_count int) ([]net.IP, error) {
	if random_count < 1 {
		return nil, fmt.Errorf("random count can not < 1")
	}
	ip_list, err := GetIPSFromCIDR(cidr)
	if err != nil {
		return nil, err
	}
	// 随机取出IP列表中的5个
	if len(ip_list) <= random_count {
		return ip_list, nil
	}
	random_ip_list, err := GetRandomIPValuesInArr(ip_list, random_count)
	return random_ip_list, nil
}

func GetRandomIPValuesInArr(arr []net.IP, count int) (randomValues []net.IP, err error) {
	// 随机选择5个值的索引
	randomIndices := getRandomIndices(len(arr), count)

	// 打印随机选择的值
	randomValues = getIPValuesByIndices(arr, randomIndices)
	return
}

// getRandomIndices 生成n个不重复的随机索引
func getRandomIndices(n, count int) []int {
	// 生成所有索引的切片
	allIndices := make([]int, n)
	for i := range allIndices {
		allIndices[i] = i
	}

	// 随机打乱索引顺序
	rand.Shuffle(n, func(i, j int) {
		allIndices[i], allIndices[j] = allIndices[j], allIndices[i]
	})

	// 截取前count个索引
	return allIndices[:count]
}

// getValuesByIndices 根据索引从列表中获取值
func getValuesByIndices(list []string, indices []int) []string {
	values := make([]string, len(indices))
	for i, index := range indices {
		values[i] = list[index]
	}
	return values
}

// getValuesByIndices 根据索引从列表中获取值
func getIPValuesByIndices(list []net.IP, indices []int) []net.IP {
	values := make([]net.IP, len(indices))
	for i, index := range indices {
		values[i] = list[index]
	}
	return values
}
