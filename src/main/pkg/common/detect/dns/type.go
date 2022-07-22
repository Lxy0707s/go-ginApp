package dns

import (
	"go-ginApp/src/main/pkg/common/common_type"
)

type DNS struct {
	Target string
	Dns    *Config `yaml:"dns,omitempty"`
}

type (
	Config struct {
		IPProtocol         string                `yaml:"preferred_ip_protocol,omitempty"`
		IPProtocolFallback bool                  `yaml:"ip_protocol_fallback,omitempty"`
		DNSOverTLS         bool                  `yaml:"dns_over_tls,omitempty"`
		TLSConfig          common_type.TLSConfig `yaml:"tls_config,omitempty"`
		SourceIPAddress    string                `yaml:"source_ip_address,omitempty"`
		TransportProtocol  string                `yaml:"transport_protocol,omitempty"`
		QueryClass         string                `yaml:"query_class,omitempty"` // Defaults to IN.
		QueryName          string                `yaml:"query_name,omitempty"`
		QueryType          string                `yaml:"query_type,omitempty"`        // Defaults to ANY.
		Recursion          bool                  `yaml:"recursion_desired,omitempty"` // Defaults to true.
		ValidRcodes        []string              `yaml:"valid_rcodes,omitempty"`      // Defaults to NOERROR.
		ValidateAnswer     RRValidator           `yaml:"validate_answer_rrs,omitempty"`
		ValidateAuthority  RRValidator           `yaml:"validate_authority_rrs,omitempty"`
		ValidateAdditional RRValidator           `yaml:"validate_additional_rrs,omitempty"`
	}
	RRValidator struct {
		FailIfMatchesRegexp     []string `yaml:"fail_if_matches_regexp,omitempty"`
		FailIfAllMatchRegexp    []string `yaml:"fail_if_all_match_regexp,omitempty"`
		FailIfNotMatchesRegexp  []string `yaml:"fail_if_not_matches_regexp,omitempty"`
		FailIfNoneMatchesRegexp []string `yaml:"fail_if_none_matches_regexp,omitempty"`
	}
)
