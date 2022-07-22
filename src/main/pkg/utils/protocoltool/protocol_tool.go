package protocoltool

// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"net"
	"time"
)

var protocolToGauge = map[string]float64{
	"ip4": 4,
	"ip6": 6,
}

// Returns the IP for the IPProtocol and lookup time.
func ChooseProtocol(ctx context.Context, IPProtocol string, fallbackIPProtocol bool, target string, logger *log.Logger) (ip *net.IPAddr, lookupTime float64, err error) {

	if IPProtocol == "ip6" || IPProtocol == "" {
		IPProtocol = "ip6"
	} else {
		IPProtocol = "ip4"
	}

	//level.Info(logger).Log("msg", "Resolving target address", "target", target, "ip_protocol", IPProtocol)
	resolveStart := time.Now()

	defer func() {
		lookupTime = time.Since(resolveStart).Seconds()
	}()

	resolver := &net.Resolver{}
	if !fallbackIPProtocol {
		ips, err := resolver.LookupIP(ctx, IPProtocol, target)
		if err == nil {
			for _, ip := range ips {
				//level.Info(logger).Log("msg", "Resolved target address", "target", target, "ip", ip.String())
				return &net.IPAddr{IP: ip}, lookupTime, nil
			}
		}
		//level.Error(logger).Log("msg", "Resolution with IP protocol failed", "target", target, "ip_protocol", IPProtocol, "err", err)
		return nil, 0.0, err
	}

	ips, err := resolver.LookupIPAddr(ctx, target)
	if err != nil {
		//level.Error(logger).Log("msg", "Resolution with IP protocol failed", "target", target, "err", err)
		return nil, 0.0, err
	}

	// Return the IP in the requested protocol.
	var fallback *net.IPAddr
	for _, ip := range ips {
		switch IPProtocol {
		case "ip4":
			if ip.IP.To4() != nil {
				//level.Info(logger).Log("msg", "Resolved target address", "target", target, "ip", ip.String())
				return &ip, lookupTime, nil
			}

			// ip4 as fallback
			fallback = &ip

		case "ip6":
			if ip.IP.To4() == nil {
				//level.Info(logger).Log("msg", "Resolved target address", "target", target, "ip", ip.String())
				return &ip, lookupTime, nil
			}

			// ip6 as fallback
			fallback = &ip
		}
	}

	// Unable to find ip and no fallback set.
	if fallback == nil || !fallbackIPProtocol {
		return nil, 0.0, fmt.Errorf("unable to find ip; no fallback")
	}
	//level.Info(logger).Log("msg", "Resolved target address", "target", target, "ip", fallback.String())
	return fallback, lookupTime, nil
}

func ipHash(ip net.IP) float64 {
	h := fnv.New32a()
	if ip.To4() != nil {
		h.Write(ip.To4())
	} else {
		h.Write(ip.To16())
	}
	return float64(h.Sum32())
}
