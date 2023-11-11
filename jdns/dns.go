package jdns

import (
	"github.com/chroblert/jlog"
	"github.com/miekg/dns"
	"net"
	"strings"
	"time"
)

func resolve(domain string) {
	c := dns.Client{
		Timeout: 5 * time.Second,
	}

	msg := dns.Msg{}
	msg.SetQuestion(domain, dns.TypeA)

	// edns
	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_SUBNET)
	e.Code = dns.EDNS0SUBNET
	e.SourceScope = 0
	e.Address = net.ParseIP("11.11.11.11")
	e.Family = 1 // IP4
	e.SourceNetmask = net.IPv4len * 8
	o.Option = append(o.Option, e)
	msg.Extra = append(msg.Extra, o)

	jlog.Info("msg", msg)

	// edns done
	jlog.Info("exchange")
	r, _, err := c.Exchange(&msg, "8.8.8.8:53")
	if err != nil {
		jlog.Fatal(err)
	}

	var dst []string
	for _, ans := range r.Answer {
		record, isType := ans.(*dns.A)
		if isType {
			jlog.Info("type A: ", record.A)
			dst = append(dst, record.A.String())
		}
		record1, isType := ans.(*dns.CNAME)
		if isType {
			jlog.Info("type CNAME: ", record1.Target)
		}
		jlog.Debug(ans.String())
	}

	for _, v := range dst {
		jlog.Info("ok: ", v)
	}
}

func ResolveA(server string, domain string, clientIP string) (domain_path []string, dest_ip []string, err error) {
	// queryType
	var qtype uint16
	qtype = dns.TypeA

	// dnsServer
	if !strings.HasSuffix(server, ":53") {
		server += ":53"
	}

	domain = dns.Fqdn(domain)

	msg := new(dns.Msg)
	msg.SetQuestion(domain, qtype)
	msg.RecursionDesired = true

	if clientIP != "" {
		opt := new(dns.OPT)
		opt.Hdr.Name = "."
		opt.Hdr.Rrtype = dns.TypeOPT
		e := new(dns.EDNS0_SUBNET)
		e.Code = dns.EDNS0SUBNET
		e.Family = 1 // ipv4
		e.SourceNetmask = 24
		e.SourceScope = 0
		e.Address = net.ParseIP(clientIP).To4()
		opt.Option = append(opt.Option, e)
		msg.Extra = []dns.RR{opt}
	}

	client := &dns.Client{
		DialTimeout:  5 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	resp, _, err := client.Exchange(msg, server)
	if err != nil {
		//jlog.Error(err)
		return
	}
	//var dst_ip []string
	//var domain_path []string
	for _, ans := range resp.Answer {
		record, isType := ans.(*dns.A)
		if isType {
			//jlog.Info("type A: ", record.A)
			dest_ip = append(dest_ip, record.A.String())
		}
		record1, isType := ans.(*dns.CNAME)
		if isType {
			//jlog.Info("type CNAME: ", record1.Target)
			domain_path = append(domain_path, record1.Target)
		}
		//jlog.Debug(ans.String())
	}
	return
}

func GetDomainMap(domain string, domain_path []string, dest_ips []string) (domain_cname_map map[string]string, domain_a_map map[string][]string) {
	//var domain_cname_map map[string]string
	//var domain_a_map map[string][]string
	last_domain := domain
	domain_cname_map = make(map[string]string, len(domain_path))
	if len(domain_path) > 0 {
		for _, v := range domain_path {
			domain_cname_map[last_domain] = v
			last_domain = v
		}
	}
	if len(dest_ips) > 0 {
		domain_a_map = make(map[string][]string, 1)
		domain_a_map[last_domain] = dest_ips
	}
	return
}
