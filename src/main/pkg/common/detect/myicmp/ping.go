package myicmp

import (
	"errors"
	"fmt"
	"go-ginApp/src/main/pkg/utils/funcs/ips"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Ping self define function for ping command
func Ping(ip string, timeout int, data []byte) (int, float64, error) {
	isIPv4, isIPv6 := ips.ParseAndDistinctIP(ip)
	if !isIPv4 && !isIPv6 {
		return 1, 0.0, fmt.Errorf("input ip string is not a legal ip")
	}
	if isIPv4 {
		return PingIPv4(ip, timeout, data)
	}
	return PingIPv6(ip, timeout, data)
}

//PingIPv4 send icmp package to ipv4
func PingIPv4(ip string, timeout int, data []byte) (int, float64, error) {
	ping, err := run(ip, 8, data)
	if err != nil {
		// zap.Debug("myicmp : ping run err", "err", err)
		return 1, 0.0, err
	}
	lost, rtt, err := ping.ping(timeout)
	ping.close()
	return lost, rtt, err
}

func run(addr string, req int, data []byte) (*ping, error) {
	addr, err := lookup(addr)
	if err != nil {
		return nil, err
	}
	pid := os.Getpid() & 0xffff
	//seq := atomic.AddInt64(&common.Seq, 1)
	strs := strings.Split(addr, ".")
	if len(strs) < 4 {
		return nil, fmt.Errorf("ip format err")
	}
	tmp, _ := strconv.ParseInt(strs[3], 10, 0)
	seq := time.Now().Nanosecond() + int(tmp)

	wb, err := marshalMsg(req, data, pid, int(seq))
	if err != nil {
		return nil, err
	}
	return &ping{Data: wb, Addr: addr}, nil
}

func marshalMsg(req int, data []byte, pid int, seq int) ([]byte, error) {

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   pid,
			Seq:  seq,
			Data: data,
		},
	}
	return wm.Marshal(nil)
}

func lookup(host string) (string, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}
	if len(addrs) < 1 {
		return "", errors.New("unknown host")
	}
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return addrs[rd.Intn(len(addrs))], nil
}

type (
	ping struct {
		Addr string
		Conn net.Conn
		Data []byte
	}
	reply struct {
		Time  float64
		TTL   uint8
		Error error
	}
)

func (selfPing *ping) dail() (err error) {
	selfPing.Conn, err = net.Dial("ip4:icmp", selfPing.Addr)
	if err != nil {
		return err
	}
	return nil
}

func (selfPing *ping) setDeadline(timeout int) error {
	return selfPing.Conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (selfPing *ping) close() error {
	if selfPing.Conn == nil {
		return nil
	}
	return selfPing.Conn.Close()
}

func (selfPing *ping) ping(pingTimeout int) (int, float64, error) {
	if err := selfPing.dail(); err != nil {
		// zap.Debug("myicmp : ping dial err", "err", err)
		return 1, 0, err
	}
	selfPing.setDeadline(pingTimeout)

	r := sendPingMsg(selfPing.Conn, selfPing.Data)
	if r.Error != nil {
		// zap.Debug("myicmp : ping return err", "err", r.Error)
		return 1, 0, r.Error
	}
	//time.Sleep(1e9)
	return 0, r.Time, nil
}

func sendPingMsg(c net.Conn, wb []byte) (r reply) {
	start := time.Now()

	if _, r.Error = c.Write(wb); r.Error != nil {
		return
	}

	var rm *icmp.Message
	var duration time.Duration
	var ttl uint8
	readTime := 0
	fail := true
	for {
		rb := make([]byte, 1500)
		var n int
		n, r.Error = c.Read(rb)
		if r.Error != nil {
			return
		}
		readTime = readTime + 1
		duration = time.Now().Sub(start)

		ttl = uint8(rb[8])
		rb = func(b []byte) []byte {
			if len(b) < 20 {
				return b
			}
			hdrlen := int(b[0]&0x0f) << 2
			return b[hdrlen:]
		}(rb)

		rm, r.Error = icmp.ParseMessage(1, rb[:n])
		if r.Error != nil {
			return
		}

		if len(wb) > 8 && len(rb) > 8 && wb[4] == rb[4] && wb[5] == rb[5] && wb[6] == rb[6] && wb[7] == rb[7] {
			fail = false
			break
		}

		if len(wb) > 8 && len(rb) > 36 && wb[4] == rb[32] && wb[5] == rb[33] && wb[6] == rb[34] && wb[7] == rb[35] {
			fail = false
			break
		}
	}
	if fail {
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		t := float64(duration.Nanoseconds()) / 1000000
		r = reply{t, ttl, nil}
	case ipv4.ICMPTypeDestinationUnreachable:
		r.Error = errors.New("Destination Unreachable")
	default:
		r.Error = fmt.Errorf("Not ICMPTypeEchoReply %v", rm)
	}
	return
}
