package grpc_check

// Copyright 2021 The Prometheus Authors
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
	"go-ginApp/src/main/internal/config"
	"go-ginApp/src/main/pkg/common/common_type"
	"go-ginApp/src/main/pkg/utils/protocoltool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/url"
	"strings"
)

type GRPCHealthCheck interface {
	Check(c context.Context, service string) (bool, codes.Code, *peer.Peer, string, error)
}

type gRPCHealthCheckClient struct {
	client grpc_health_v1.HealthClient
	conn   *grpc.ClientConn
}

func NewGrpcHealthCheckClient(conn *grpc.ClientConn) GRPCHealthCheck {
	client := new(gRPCHealthCheckClient)
	client.client = grpc_health_v1.NewHealthClient(conn)
	client.conn = conn
	return client
}

func (c *gRPCHealthCheckClient) Close() error {
	return c.conn.Close()
}

func (c *gRPCHealthCheckClient) Check(ctx context.Context, service string) (bool, codes.Code, *peer.Peer, string, error) {
	var res *grpc_health_v1.HealthCheckResponse
	var err error
	req := grpc_health_v1.HealthCheckRequest{
		Service: service,
	}

	serverPeer := new(peer.Peer)
	res, err = c.client.Check(ctx, &req, grpc.Peer(serverPeer))
	if err == nil {
		if res.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
			return true, codes.OK, serverPeer, res.Status.String(), nil
		}
		return false, codes.OK, serverPeer, res.Status.String(), nil
	}

	returnStatus, _ := status.FromError(err)

	return false, returnStatus.Code(), nil, "", err
}

func ProbeGRPC(ctx context.Context, target string, module Module, logger *log.Logger) (success bool) {
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		config.AppLog.Error("msg", "Could not parse target URL", "err", err)
		return false
	}

	targetHost, targetPort, err := net.SplitHostPort(targetURL.Host)
	// If split fails, assuming it's a hostname without port part.
	if err != nil {
		targetHost = targetURL.Host
	}

	tlsConfig, err := common_type.NewTLSConfig(&module.GRPC.TLSConfig)
	if err != nil {
		//level.Error(logger).Log("msg", "Error creating TLS configuration", "err", err)
		return false
	}

	ip, _, err := protocoltool.ChooseProtocol(ctx, module.GRPC.PreferredIPProtocol, module.GRPC.IPProtocolFallback, targetHost, logger)
	if err != nil {
		//level.Error(logger).Log("msg", "Error resolving address", "err", err)
		return false
	}
	//checkStart := time.Now()
	if len(tlsConfig.ServerName) == 0 {
		// If there is no `server_name` in tls_config, use
		// the hostname of the target.
		tlsConfig.ServerName = targetHost
	}

	if targetPort == "" {
		targetURL.Host = "[" + ip.String() + "]"
	} else {
		targetURL.Host = net.JoinHostPort(ip.String(), targetPort)
	}

	var opts []grpc.DialOption
	target = targetHost + ":" + targetPort
	if !module.GRPC.TLS {
		//level.Debug(logger).Log("msg", "Dialing GRPC without TLS")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if len(targetPort) == 0 {
			target = targetHost + ":80"
		}
	} else {
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
		if len(targetPort) == 0 {
			target = targetHost + ":443"
		}
	}

	conn, err := grpc.Dial(target, opts...)

	if err != nil {
		//level.Error(logger).Log("did not connect: %v", err)
	}

	client := NewGrpcHealthCheckClient(conn)
	defer conn.Close()
	ok, statusCode, serverPeer, servingStatus, err := client.Check(context.Background(), module.GRPC.Service)

	fmt.Println(statusCode)
	for servingStatusName, _ := range grpc_health_v1.HealthCheckResponse_ServingStatus_value {
		//servingStatusName
		fmt.Println("servingStatusName", servingStatusName)
	}
	if servingStatus != "" {
		fmt.Println("status:", servingStatus)
	}

	if serverPeer != nil {
		tlsInfo, tlsOk := serverPeer.AuthInfo.(credentials.TLSInfo)
		if tlsOk {
			fmt.Println(float64(getEarliestCertExpiry(&tlsInfo.State).Unix()))
			fmt.Println(&tlsInfo.State)
		} else {

		}
	}

	if !ok || err != nil {
		//level.Error(logger).Log("msg", "can't connect grpc server:", "err", err)
		success = false
	} else {
		//level.Debug(logger).Log("connect the grpc server successfully")
		success = true
	}

	return
}
