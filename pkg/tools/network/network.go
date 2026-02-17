// Package network provides network utility functions.
package network

import (
	"net"
	"time"
)

// DefaultDNSEndpoints contains common DNS endpoints used for network connectivity checks.
var DefaultDNSEndpoints = []string{
	"8.8.8.8:53",        // Google DNS
	"1.1.1.1:53",        // Cloudflare DNS
	"208.67.222.222:53", // OpenDNS
}

// DefaultTimeout is the default timeout for network connectivity checks.
const DefaultTimeout = 3 * time.Second

// HasNetwork checks if the system has network connectivity by attempting to
// connect to multiple DNS endpoints. Returns true if any connection succeeds.
// This provides better reliability in restricted network environments.
func HasNetwork() bool {
	return HasNetworkWithEndpoints(DefaultDNSEndpoints, DefaultTimeout)
}

// HasNetworkWithEndpoints checks network connectivity using the specified DNS endpoints
// and timeout. Returns true if any connection succeeds.
func HasNetworkWithEndpoints(endpoints []string, timeout time.Duration) bool {
	for _, endpoint := range endpoints {
		conn, err := net.DialTimeout("tcp", endpoint, timeout)
		if err != nil {
			continue
		}
		conn.Close()
		return true
	}

	return false
}
