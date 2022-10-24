package spf

import (
	"context"
	"fmt"
	"net"
	"strings"
)

func networkCIDR(ip, prefix string) (*net.IPNet, error) {
	if prefix == "" {
		ip := net.ParseIP(ip)

		if ip.To4() != nil {
			prefix = "32"
		} else {
			prefix = "128"
		}
	}

	cidrStr := fmt.Sprintf("%s/%s", ip, prefix)

	_, network, err := net.ParseCIDR(cidrStr)
	return network, err
}

func ipInNetworks(ip net.IP, networks []*net.IPNet) bool {
	for _, network := range networks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

func buildNetworks(ips []string, prefix string) []*net.IPNet {
	var networks []*net.IPNet

	for _, ip := range ips {
		network, err := networkCIDR(ip, prefix)
		if err == nil {
			networks = append(networks, network)
		}
	}

	return networks
}

func aNetworks(ctx context.Context, m *Mechanism) []*net.IPNet {
	resolver := &net.Resolver{}
	ips, _ := resolver.LookupHost(ctx, m.Domain)

	return buildNetworks(ips, m.Prefix)
}

func mxNetworks(ctx context.Context, m *Mechanism) []*net.IPNet {
	var networks []*net.IPNet
	resolver := &net.Resolver{}

	mxs, _ := resolver.LookupMX(ctx, m.Domain)

	for _, mx := range mxs {
		ips, _ := net.LookupHost(mx.Host)
		networks = append(networks, buildNetworks(ips, m.Prefix)...)
	}

	return networks
}

func testPTR(ctx context.Context, m *Mechanism, ip string) bool {
	resolver := &net.Resolver{}
	names, err := resolver.LookupAddr(ctx, ip)

	if err != nil {
		return false
	}

	for _, name := range names {
		if strings.HasSuffix(name, m.Domain) {
			return true
		}
	}

	return false
}
