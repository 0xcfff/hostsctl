package iptools

import "regexp"

var (
	rxIPv4Address = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	rxIPv6Address = regexp.MustCompile(`^[\da-fA-F]{0,}:[\da-fA-F]{0,}:[\da-fA-F]{0,}(:[\da-fA-F]{0,}){0,5}$`)
)


// Returns true if specified value is IPv4
func IsIPv4(value string) bool {
	return rxIPv4Address.MatchString(value)
}

// Returns true if specified value is IPv6
func IsIPv6(value string) bool {
	return rxIPv6Address.MatchString(value)
}

// Returns true if specified value is IPv6
func IsIP(value string) bool {
	return rxIPv4Address.MatchString(value) || rxIPv6Address.MatchString(value)
}
