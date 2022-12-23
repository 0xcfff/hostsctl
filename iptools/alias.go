package iptools

import (
	"net"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	systemAliases = []struct{
			ip net.IP
			alias string
	}{
		{net.IPv4(127,0,0,1), "localhost"},
		{net.ParseIP("::1"), "ip6-localhost"},
		{net.ParseIP("::1"), "ip6-loopback"},
	}
)

func IsSystemAlias(ip string, alias string) bool {
	ipTrimmed := strings.TrimSpace(ip)
	aliasTrimmed := strings.ToLower(strings.TrimSpace(alias))

	if IsIP(ipTrimmed) {
		ip := net.ParseIP(ipTrimmed)

		for _, p := range systemAliases {
			if slices.Compare(p.ip, ip) == 0 && p.alias == aliasTrimmed {
				return true
			}
		}
	}

	return false
}
