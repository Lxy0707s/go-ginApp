package myicmp

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
)

//PingIPv6 send icmp package to ipv6
func PingIPv6(ip string, timeout int, data []byte) (int, float64, error) {
	ping, err := runIPv6(ip, 8, data)
	if err != nil {
		// zap.Debug("myicmp : ping run err", "err", err)
		return 1, 0.0, err
	}
	lost, rtt, err := ping.pingIPv6(timeout)
	ping.close()
	return lost, rtt, err
}

func runIPv6(addr string, req int, data []byte) (*ping, error) {
	pid := os.Getpid() & 0xffff
	//seq := atomic.AddInt64(&common.Seq, 1)
	strs := strings.Split(addr, ":")
	tmp, _ := strconv.ParseInt(strs[len(strs)-1], 16, 0)
	seq := time.Now().Nanosecond() + int(tmp)

	wb, err := marshalMsgIPv6(req, data, pid, int(seq))
	if err != nil {
		return nil, err
	}
	addr, err = lookup(addr)

	if err != nil {
		return nil, err
	}
	return &ping{Data: wb, Addr: addr}, nil
}

func marshalMsgIPv6(req int, data []byte, pid int, seq int) ([]byte, error) {

	wm := icmp.Message{
		Type: ipv6.ICMPTypeEchoRequest,
		Code: 0,
		Body: &icmp.Echo{
			ID:   pid,
			Seq:  seq,
			Data: data,
		},
	}
	return wm.Marshal(nil)
}

func (selfPing *ping) pingIPv6(pingTimeout int) (int, float64, error) {
	if err := selfPing.dailIPv6(); err != nil {
		// zap.Debug("myicmp : ping dial err", "err", err)
		return 1, 0, err
	}
	selfPing.setDeadline(pingTimeout)

	r := sendPingIPv6Msg(selfPing.Conn, selfPing.Data)
	if r.Error != nil {
		// zap.Debug("myicmp : ping return err", "err", r.Error)
		return 1, 0, r.Error
	}
	//time.Sleep(1e9)
	return 0, r.Time, nil
}

func (selfPing *ping) dailIPv6() (err error) {
	selfPing.Conn, err = net.Dial("ip6:ipv6-icmp", selfPing.Addr)
	if err != nil {
		return err
	}
	return nil
}

func sendPingIPv6Msg(c net.Conn, wb []byte) (r reply) {
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

		rm, r.Error = icmp.ParseMessage(58, rb[:n])
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
	case ipv6.ICMPTypeEchoReply:
		t := float64(duration.Nanoseconds()) / 1000000
		r = reply{t, ttl, nil}
	case ipv6.ICMPTypeDestinationUnreachable:
		r.Error = errors.New("destination Unreachable")
	default:
		r.Error = fmt.Errorf("not ICMPTypeEchoReply %v", rm)
	}
	return
}
