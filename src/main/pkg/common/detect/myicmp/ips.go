package myicmp

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
)

type (
	//IPRange ip range
	IPRange struct {
		FromIP, ToIP string
	}
)

// LocalIP to get first IP
func LocalIP() string {
	address, err := GetIP(true)
	if err != nil {
		return ""
	}
	if len(address) == 0 {
		return ""
	}
	for _, localIP := range address {
		return localIP[0].IP.String()
	}
	return ""
}

// LocalIPs to get local IPs
func LocalIPs() []string {
	address, err := GetIP(true)
	if err != nil {
		return nil
	}
	if len(address) == 0 {
		return nil
	}
	result := make([]string, 0)
	for _, localIP := range address {
		result = append(result, localIP[0].IP.String())
	}
	sort.Strings(result)
	return result
}

// GetIP return all ip of this computer.
// when skipIntranet is true, return only external ip
// when skipIntranet is false, return all ip include internal ip
func GetIP(skipIntranet bool) (map[string][]*net.IPNet, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	result := make(map[string][]*net.IPNet)
	for _, iface := range ifaces {
		ipList := make([]*net.IPNet, 0)
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		if strings.HasPrefix(iface.Name, "docker") {
			continue
		}
		if strings.HasPrefix(iface.Name, "w-") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				// todo: ipv6
				continue
			}
			if !skipIntranet {
				ipList = append(ipList, addr.(*net.IPNet))
				continue
			}
			if !IsIntranet(ip.String()) {
				ipList = append(ipList, addr.(*net.IPNet))
			}
		}
		if len(ipList) > 0 {
			name := iface.Name
			result[name] = ipList
		}
	}
	return result, nil
}

//GetIPV6 get ip v6
func GetIPV6(skipIntranet bool) (map[string]*net.IPNet, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	result := make(map[string]*net.IPNet)
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		if strings.HasPrefix(iface.Name, "docker") {
			continue
		}
		if strings.HasPrefix(iface.Name, "w-") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip.To4() != nil {
				// todo: ipv6
				continue
			}
			name := iface.Name
			if !skipIntranet {
				result[name] = addr.(*net.IPNet)
				continue
			}
			if !IsIntranet(ip.String()) {
				result[name] = addr.(*net.IPNet)
			}
		}
	}
	return result, nil
}

// IsIntranet return whether the ip is internal or not
func IsIntranet(ipStr string) bool {
	if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") {
		return true
	}
	if strings.HasPrefix(ipStr, "172.") {
		// 172.16.0.0-172.31.255.255
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}
		second, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return false
		}
		if second >= 16 && second <= 31 {
			return true
		}
	}
	return false
}

//ToDotted change ip net to dotted string
func ToDotted(ip *net.IPNet) string {
	if ip == nil {
		return ""
	}
	ipBytes := []byte(ip.IP.To4())
	return fmt.Sprintf("%x.%x.%x.%x", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
}

// SegmentIPs return all ip in subnet
func SegmentIPs(ipNet *net.IPNet) ([]string, error) {
	var ips []string
	ip := ipNet.IP.Mask(ipNet.Mask)
	for ; ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	if len(ips) <= 1 {
		return nil, errors.New("no segement ips")
	}
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

//GetIPRangesOfCSegment get ip ranges of c segment
func GetIPRangesOfCSegment(fromIP, toIP string) []*IPRange {
	result := make([]*IPRange, 0)
	ipBegin := net.ParseIP(fromIP)
	ipEnd := net.ParseIP(toIP)
	equal := false
	ipSegs := make(map[string]*IPRange)
	for {
		if ipBegin.Equal(ipEnd) {
			equal = true
		}
		ip := ipBegin.String()
		ipseg := GetNetSegment(ip)
		if ipSegs[ipseg] == nil {
			ipSegs[ipseg] = &IPRange{
				FromIP: ip,
				ToIP:   ip,
			}
		} else {
			ipSegs[ipseg].ToIP = ip
		}
		if equal {
			break
		}
		Inc(ipBegin)
	}
	for _, ipRange := range ipSegs {
		result = append(result, ipRange)
	}
	return result
}

// GetNetSegment return c segment information of the ip
func GetNetSegment(ip string) string {
	idx := strings.LastIndex(ip, ".")
	if idx < 0 {
		return ""
	}
	return ip[:idx]
}

// Inc increase input ip addr
func Inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

//GetLocalIP get local ip
func GetLocalIP() (localIP string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
			}
		}
	}
	if len(localIP) == 0 {
		err = fmt.Errorf("local ip is empty")
		return
	}
	return
}

//ParseAndDistinctIP parse the input string is a ip address, and is ipv4 or ipv6
func ParseAndDistinctIP(ipstr string) (isIPv4, isIPv6 bool) {
	// ip := net.ParseIP(ipstr)
	// if ip == nil {
	// 	return false, false
	// }
	if strings.Contains(ipstr, ".") {
		return true, false
	}
	return false, true
}
