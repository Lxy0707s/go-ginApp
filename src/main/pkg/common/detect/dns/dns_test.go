package dns

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"testing"
	"time"
)

func TestNormal(t *testing.T) {
	d := NewDNS()
	var dDnsUrl = "uxks3.v.bsclink.cn"
	d.SetTarget("211.139.5.30:53")
	d.SetDNSConfig(&Config{
		IPProtocol:         "ip4",
		IPProtocolFallback: true,
		Recursion:          true,
		QueryType:          "A",
		QueryName:          dDnsUrl,
		QueryClass:         "IN",
	})
	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	d.DnsRun(testCTX)
}

func TestServfailDNSResponse(t *testing.T) {
	fmt.Println("--------", 1/10)
	d := NewDNS()
	server, addr := startDNSServer("udp", dns.HandleFailed)
	defer server.Shutdown()
	fmt.Println(addr.String())
	d.SetTarget(addr.String()) //"192.168.220.2:53"
	d.SetDNSConfig(&Config{
		IPProtocol:         "ip4",
		IPProtocolFallback: true,
		Recursion:          true,
		QueryType:          "A",
		QueryName:          "www.baidu.com.",
		QueryClass:         "IN",
	})
	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	d.DnsRun(testCTX)
}

func TestAuthoritativeDNSResponse(t *testing.T) {
	d := NewDNS()
	server, addr := startDNSServer("udp", authoritativeDNSHandler)
	defer server.Shutdown()
	d.SetTarget(addr.String()) //"192.168.220.2:53"
	d.SetDNSConfig(&Config{
		IPProtocol:         "ip4",
		IPProtocolFallback: true,
		QueryClass:         "CH",
		QueryName:          "example.com",
		QueryType:          "TXT",
		ValidateAnswer: RRValidator{
			FailIfMatchesRegexp:    []string{".*IN.*"},
			FailIfNotMatchesRegexp: []string{".*CH.*"},
		},
	})
	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	d.DnsRun(testCTX)
}

func startDNSServer(protocol string, handler func(dns.ResponseWriter, *dns.Msg)) (*dns.Server, net.Addr) {
	h := dns.NewServeMux()
	h.HandleFunc(".", handler)

	server := &dns.Server{Addr: ":0", Net: protocol, Handler: h}
	if protocol == "udp" {
		a, err := net.ResolveUDPAddr(server.Net, server.Addr)
		if err != nil {
			panic(err)
		}
		l, err := net.ListenUDP(server.Net, a)
		if err != nil {
			panic(err)
		}
		server.PacketConn = l
	} else {
		a, err := net.ResolveTCPAddr(server.Net, server.Addr)
		if err != nil {
			panic(err)
		}
		l, err := net.ListenTCP(server.Net, a)
		if err != nil {
			panic(err)
		}
		server.Listener = l
	}
	go server.ActivateAndServe()

	if protocol == "tcp" {
		return server, server.Listener.Addr()
	}
	return server, server.PacketConn.LocalAddr()
}

func authoritativeDNSHandler(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	if r.Question[0].Qtype == dns.TypeSOA {
		a, err := dns.NewRR("example.com. 3600 IN SOA ns.example.com. noc.example.com. 1000 7200 3600 1209600 3600")
		if err != nil {
			panic(err)
		}
		m.Answer = append(m.Answer, a)
	} else if r.Question[0].Qclass == dns.ClassCHAOS && r.Question[0].Qtype == dns.TypeTXT {
		txt, err := dns.NewRR("example.com. 3600 CH TXT \"goCHAOS\"")
		if err != nil {
			panic(err)
		}
		m.Answer = append(m.Answer, txt)
	} else {
		a, err := dns.NewRR("example.com. 3600 IN A 127.0.0.1")
		if err != nil {
			panic(err)
		}
		m.Answer = append(m.Answer, a)
	}

	authority := []string{
		"example.com. 7200 IN NS ns1.isp.net.",
		"example.com. 7200 IN NS ns2.isp.net.",
	}
	for _, rr := range authority {
		a, err := dns.NewRR(rr)
		if err != nil {
			panic(err)
		}
		m.Ns = append(m.Ns, a)
	}

	additional := []string{
		"ns1.isp.net. 7200 IN A 127.0.0.1",
		"ns1.isp.net. 7200 IN AAAA ::1",
		"ns2.isp.net. 7200 IN A 127.0.0.2",
	}
	for _, rr := range additional {
		a, err := dns.NewRR(rr)
		if err != nil {
			panic(err)
		}
		m.Extra = append(m.Extra, a)
	}

	if err := w.WriteMsg(m); err != nil {
		panic(err)
	}
}
