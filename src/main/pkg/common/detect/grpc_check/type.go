package grpc_check

import (
	"go-ginApp/src/main/pkg/common/common_type"
	"time"
)

type (
	Module struct {
		Timeout time.Duration `yaml:"timeout,omitempty"`
		GRPC    GRPCProbe     `yaml:"grpc,omitempty"`
	}
	GRPCProbe struct {
		Service             string                `yaml:"service,omitempty"`
		TLS                 bool                  `yaml:"tls,omitempty"`
		TLSConfig           common_type.TLSConfig `yaml:"tls_config,omitempty"`
		IPProtocolFallback  bool                  `yaml:"ip_protocol_fallback,omitempty"`
		PreferredIPProtocol string                `yaml:"preferred_ip_protocol,omitempty"`
	}
)
