package jdns

import (
	"github.com/miekg/dns"
	"net"
	"strconv"
	"strings"
	"time"
)

func ReverseResolveIP(server string, ip_str string, clientIP string) (r_domain string, err error) {
	// queryType
	var qtype uint16
	qtype = dns.TypePTR

	// dnsServer
	if !strings.HasSuffix(server, ":53") {
		server += ":53"
	}

	ip_reverse_str, err := ReverseAddr(ip_str)
	if err != nil {
		return
	}

	msg := new(dns.Msg)
	msg.SetQuestion(ip_reverse_str, qtype)
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
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	resp, _, err := client.Exchange(msg, server)
	if err != nil {
		//jlog.Error(err)
		return
	}
	for _, answer := range resp.Answer {
		if ptr, ok := answer.(*dns.PTR); ok {
			//fmt.Printf("%s\n", ptr.Ptr)
			return ptr.Ptr, nil
		}
	}
	return
}

type Error struct{ err string }

func (e *Error) Error() string {
	if e == nil {
		return "dns: <nil>"
	}
	return "dns: " + e.err
}

const hexDigit = "0123456789abcdef"

// ReverseAddr returns the in-addr.arpa. or ip6.arpa. hostname of the IP
// address suitable for reverse DNS (PTR) record lookups or an error if it fails
// to parse the IP address.
func ReverseAddr(addr string) (arpa string, err error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", &Error{err: "unrecognized address: " + addr}
	}
	if ip.To4() != nil {
		return strconv.Itoa(int(ip[15])) + "." + strconv.Itoa(int(ip[14])) + "." + strconv.Itoa(int(ip[13])) + "." +
			strconv.Itoa(int(ip[12])) + ".in-addr.arpa.", nil
	}
	// Must be IPv6
	buf := make([]byte, 0, len(ip)*4+len("ip6.arpa."))
	// Add it, in reverse, to the buffer
	for i := len(ip) - 1; i >= 0; i-- {
		v := ip[i]
		buf = append(buf, hexDigit[v&0xF])
		buf = append(buf, '.')
		buf = append(buf, hexDigit[v>>4])
		buf = append(buf, '.')
	}
	// Append "ip6.arpa." and return (buf already has the final .)
	buf = append(buf, "ip6.arpa."...)
	return string(buf), nil
}
