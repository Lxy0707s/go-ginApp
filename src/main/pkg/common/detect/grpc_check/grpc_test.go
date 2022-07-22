package grpc_check

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go-ginApp/src/main/pkg/common/common_type"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"
)

func TestGRPCConnection(t *testing.T) {

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error listening on socket: %s", err)
	}
	defer ln.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("Error retrieving port for socket: %s", err)
	}
	s := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("service", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	go func() {
		if err := s.Serve(ln); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()
	defer s.GracefulStop()

	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := ProbeGRPC(testCTX, "localhost:"+port,
		Module{Timeout: time.Second, GRPC: GRPCProbe{
			IPProtocolFallback: false,
		},
		}, nil)

	if !result {
		t.Fatalf("GRPC probe failed")
	}

}

func TestMultipleGRPCservices(t *testing.T) {

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error listening on socket: %s", err)
	}
	defer ln.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("Error retrieving port for socket: %s", err)
	}
	s := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("service1", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("service2", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	go func() {
		if err := s.Serve(ln); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()
	defer s.GracefulStop()

	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resultService1 := ProbeGRPC(testCTX, "localhost:"+port,
		Module{Timeout: time.Second, GRPC: GRPCProbe{
			IPProtocolFallback: false,
			Service:            "service1",
		},
		}, nil)

	if !resultService1 {
		t.Fatalf("GRPC probe failed for service1")
	}

	resultService2 := ProbeGRPC(testCTX, "localhost:"+port,
		Module{Timeout: time.Second, GRPC: GRPCProbe{
			IPProtocolFallback: false,
			Service:            "service2",
		},
		}, nil)

	if resultService2 {
		t.Fatalf("GRPC probe succeed for service2")
	}
}

func TestGRPCTLSConnection(t *testing.T) {

	certExpiry := time.Now().AddDate(0, 0, 1)
	testCertTmpl := generateCertificateTemplate(certExpiry, false)
	testCertTmpl.IsCA = true
	_, testcertPem, testKey := generateSelfSignedCertificate(testCertTmpl)

	// CAFile must be passed via filesystem, use a tempfile.
	tmpCaFile, err := ioutil.TempFile("", "cafile.pem")
	if err != nil {
		t.Fatalf("Error creating CA tempfile: %s", err)
	}
	if _, err = tmpCaFile.Write(testcertPem); err != nil {
		t.Fatalf("Error writing CA tempfile: %s", err)
	}
	if err = tmpCaFile.Close(); err != nil {
		t.Fatalf("Error closing CA tempfile: %s", err)
	}
	defer os.Remove(tmpCaFile.Name())

	testKeyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(testKey)})
	testcert, err := tls.X509KeyPair(testcertPem, testKeyPem)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode TLS testing keypair: %s\n", err))
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{testcert},
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS12,
	}

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error listening on socket: %s", err)
	}
	defer ln.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("Error retrieving port for socket: %s", err)
	}

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	healthServer := health.NewServer()
	healthServer.SetServingStatus("service", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	go func() {
		if err := s.Serve(ln); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()
	defer s.GracefulStop()

	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := ProbeGRPC(testCTX, "localhost:"+port,
		Module{Timeout: time.Second, GRPC: GRPCProbe{
			TLS:                true,
			TLSConfig:          common_type.TLSConfig{InsecureSkipVerify: true},
			IPProtocolFallback: false,
		},
		}, nil)

	if !result {
		t.Fatalf("GRPC probe failed")
	}
}

func TestNoTLSConnection(t *testing.T) {

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error listening on socket: %s", err)
	}
	defer ln.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("Error retrieving port for socket: %s", err)
	}
	s := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("service", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	go func() {
		if err := s.Serve(ln); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()
	defer s.GracefulStop()

	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := ProbeGRPC(testCTX, "localhost:"+port,
		Module{Timeout: time.Second, GRPC: GRPCProbe{
			TLS:                true,
			TLSConfig:          common_type.TLSConfig{InsecureSkipVerify: true},
			IPProtocolFallback: false,
		},
		}, nil)

	if result {
		t.Fatalf("GRPC probe succeed")
	}
}

func TestGRPCServiceNotFound(t *testing.T) {

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error listening on socket: %s", err)
	}
	defer ln.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("Error retrieving port for socket: %s", err)
	}
	s := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("service", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	go func() {
		if err := s.Serve(ln); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()
	defer s.GracefulStop()

	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := ProbeGRPC(testCTX, "localhost:"+port,
		Module{Timeout: time.Second, GRPC: GRPCProbe{
			IPProtocolFallback: false,
			Service:            "NonExistingService",
		},
		}, nil)

	if result {
		t.Fatalf("GRPC probe succeed")
	}
}

func TestGRPCHealthCheckUnimplemented(t *testing.T) {

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error listening on socket: %s", err)
	}
	defer ln.Close()

	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("Error retrieving port for socket: %s", err)
	}
	s := grpc.NewServer()

	go func() {
		if err := s.Serve(ln); err != nil {
			t.Errorf("failed to serve: %v", err)
			return
		}
	}()
	defer s.GracefulStop()

	testCTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := ProbeGRPC(testCTX, "localhost:"+port,
		Module{Timeout: time.Second, GRPC: GRPCProbe{
			IPProtocolFallback: false,
			Service:            "NonExistingService",
		},
		}, nil)

	if result {
		t.Fatalf("GRPC probe succeed")
	}
}
