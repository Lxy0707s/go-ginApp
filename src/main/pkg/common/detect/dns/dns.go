package dns

import (
	"context"
	"fmt"
	"github.com/miekg/dns"
	"go-ginApp/src/main/pkg/utils/myfile"
	"net"
	"strconv"
	"time"
)

var Instance *DNS

func NewDNS() *DNS {
	if Instance == nil {
		Instance = &DNS{
			Target: "",
			Dns:    &Config{},
		}
	}
	return Instance
}

func (ds *DNS) SetTarget(target string) {
	if ds == nil {
		fmt.Println("dns instance is nil")
	}
	ds.Target = target
}

func (ds *DNS) SetDNSConfig(dNSConfig *Config) {
	if ds == nil {
		fmt.Println("dNSConfig is nil")
		return
	}
	ds.Dns = dNSConfig
}

func (ds *DNS) DnsRun(ctx context.Context) {
	client := new(dns.Client)
	if ds.Dns == nil {
		fmt.Println("dns config is error.")
		return
	}
	d := ds.Dns
	qt := dns.TypeANY
	if d.QueryType != "" {
		var ok bool
		qt, ok = dns.StringToType[d.QueryType]
		if !ok {
			fmt.Println("msg", "Invalid query type", "Type seen", d.QueryType, "Existing types", dns.TypeToString)
			return
		}
	}
	qc := uint16(dns.ClassINET)
	if d.QueryClass != "" {
		var ok bool
		qc, ok = dns.StringToClass[d.QueryClass]
		if !ok {
			fmt.Println("msg", "Invalid query class", "Class seen", d.QueryClass, "Existing classes", dns.ClassToString)
			return
		}
	}

	if d.TransportProtocol == "" {
		d.TransportProtocol = "udp"
	}
	if !(d.TransportProtocol == "udp" || d.TransportProtocol == "tcp") {
		fmt.Println("protocol is wrong")
		return
	}

	var ip *net.IPAddr
	targetAddr, port, err := net.SplitHostPort(ds.Target)
	if err != nil {
		// Target only contains host so fallback to default port and set targetAddr as target.
		if d.DNSOverTLS {
			port = "853"
		} else {
			port = "53"
		}
	}
	ip, _, err = ds.chooseProtocol(ctx, d.IPProtocol, d.IPProtocolFallback, targetAddr)
	if err != nil {
		fmt.Println("msg", "Error resolving address", "err", err)
		return
	}
	//判断v46
	targetIP := net.JoinHostPort(ip.String(), port)

	var dialProtocol string
	if d.DNSOverTLS {
		if d.TransportProtocol == "tcp" {
			dialProtocol += "-tls"
		} else {
			return
		}
	}
	// Use configured SourceIPAddress.
	if len(d.SourceIPAddress) > 0 {
		srcIP := net.ParseIP(d.SourceIPAddress)
		if srcIP == nil {
			fmt.Println("msg", "Error parsing source ip address", "srcIP", d.SourceIPAddress)
			return
		}
		fmt.Println("msg", "Using local address", "srcIP", srcIP)
		client.Dialer = &net.Dialer{}
		if d.TransportProtocol == "tcp" {
			client.Dialer.LocalAddr = &net.TCPAddr{IP: srcIP}
		} else {
			client.Dialer.LocalAddr = &net.UDPAddr{IP: srcIP}
		}
	}
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = d.Recursion
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{Name: dns.Fqdn(d.QueryName), Qtype: qt, Qclass: qc}
	client.Timeout = 5 * time.Second
	response, rtt, err := client.Exchange(msg, targetIP)
	if err != nil {
		fmt.Println("dns error")
		return
	}
	var result = make(map[string][]string)
	var dst []string
	for _, ans := range response.Answer {
		switch qt {
		case dns.TypeA:
			record, isType := ans.(*dns.A)
			if isType {
				fmt.Println("Header:", record.Header(), "A:", record.A, "RTT:", rtt)
				dst = append(dst, record.Hdr.Name+"@"+record.A.String()+"@"+strconv.Itoa(int(record.Hdr.Ttl)))
			}
			record1, isType := ans.(*dns.CNAME)
			if isType {
				fmt.Println("type cname:", record1.Target)
			}
		case dns.TypeSOA:
			for _, a := range response.Answer {
				if soa, ok := a.(*dns.SOA); ok {
					fmt.Println(soa.Ns)
				}
			}
		case dns.TypeAAAA:
			for _, a := range response.Answer {
				if soa, ok := a.(*dns.AAAA); ok {
					fmt.Println(soa.AAAA.String())
					fmt.Println("----------")
				}
			}
		}
		//...
	}
	result[ip.String()] = dst
	if len(result) != 0 {
		myfile.WriteJSON(result, "dns_result.json")
	}
}

func (ds *DNS) chooseProtocol(ctx context.Context, IPProtocol string, fallbackIPProtocol bool, target string) (ip *net.IPAddr, lookupTime float64, err error) {
	if IPProtocol == "ip6" || IPProtocol == "" {
		IPProtocol = "ip6"
	} else {
		IPProtocol = "ip4"
	}

	fmt.Println("msg", "Resolving target address", "target", target, "ip_protocol", IPProtocol)
	resolveStart := time.Now()

	defer func() {
		lookupTime = time.Since(resolveStart).Seconds()
	}()

	resolver := &net.Resolver{}
	if !fallbackIPProtocol {
		ips, err := resolver.LookupIP(ctx, IPProtocol, target)
		if err == nil {
			for _, ip := range ips {
				return &net.IPAddr{IP: ip}, lookupTime, nil
			}
		}
		fmt.Println("msg", "Resolution with IP protocol failed", "target", target, "ip_protocol", IPProtocol, "err", err)
		return nil, 0.0, err
	}
	// target 是host主机
	ips, err := resolver.LookupIPAddr(ctx, target)
	if err != nil {
		fmt.Println("msg", "Resolution with IP protocol failed", "target", target, "err", err)
		return nil, 0.0, err
	}
	// Return the IP in the requested protocol.
	var fallback *net.IPAddr
	for _, ip := range ips {
		switch IPProtocol {
		case "ip4":
			if ip.IP.To4() != nil {
				fmt.Println("msg", "Resolved target address", "target", target, "ip", ip.String())
				return &ip, lookupTime, nil
			}
			// ip4 as fallback
			fallback = &ip

		case "ip6":
			if ip.IP.To4() == nil {
				fmt.Println("msg", "Resolved target address", "target", target, "ip", ip.String())
				return &ip, lookupTime, nil
			}
			// ip6 as fallback
			fallback = &ip
		}
	}
	// Unable to find ip and no fallback set.
	if fallback == nil || !fallbackIPProtocol {
		fmt.Println("unable to find ip; no fallback")
		return nil, 0.0, fmt.Errorf("unable to find ip; no fallback")
	}
	return fallback, lookupTime, nil
}
