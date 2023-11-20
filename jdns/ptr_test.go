package jdns

import (
	"fmt"
	"github.com/chroblert/jlog"
	"github.com/miekg/dns"
	"net"
	"os"
	"testing"
)

func TestReverseResolveIP(t *testing.T) {
	jlog.Info(ReverseResolveIP("127.0.0.1", "8.8.8.8", ""))
	//213.133.196.229
	r_domain, err := ReverseResolveIP("8.8.8.81", "213.133.196.229", "")
	jlog.Info("len:", len(r_domain), "err:", err)

}

func TestReverseResolveIP2(t *testing.T) {

	ip := "8.8.8.8"
	jlog.Info(ReverseAddr(ip))
	ip_a, _ := ReverseAddr(ip)
	jlog.Warn(dns.Fqdn(ip))
	// 指定要查询的IP地址
	//ip := "8.8.8.8"

	// 创建一个DNS客户端
	client := dns.Client{}

	// 创建一个DNS消息
	msg := dns.Msg{}
	msg = *msg.SetQuestion(ip_a, dns.TypePTR)
	//msg.RecursionDesired = true

	// 发送DNS查询并接收响应
	response, _, err := client.Exchange(&msg, "8.8.8.8:53") // 使用8.8.8.8作为DNS服务器，你可以更改为其他DNS服务器

	// 处理查询结果
	if err != nil {
		fmt.Printf("DNS查询错误: %v\n", err)
		os.Exit(1)
	}

	// 打印查询结果
	fmt.Printf("IP地址 %s 的PTR查询结果:\n", ip)
	jlog.Info(response.Answer)
	for _, answer := range response.Answer {
		if ptr, ok := answer.(*dns.PTR); ok {
			fmt.Printf("%s\n", ptr.Ptr)
		}
	}
}

// 有效
func TestLookupRPT(t *testing.T) {
	ptr, e := net.LookupAddr("8.8.8.8")
	if e != nil {
		return
	}
	for _, ptrval := range ptr {
		jlog.Info(ptrval, nil)
	}
	return
}
