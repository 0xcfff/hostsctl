package iptools

import "strings"

var (
	systemAliases = [][]string{{"127.0.0.1", "localhost"}}
)

func IsSystemAlias(ip string, alias string) bool {
	ipTrimmed := strings.TrimSpace(ip)
	aliasTrimmed := strings.ToLower(strings.TrimSpace(alias))

	for _, p := range systemAliases {
		if p[0] == ipTrimmed && p[1] == aliasTrimmed {
			return true
		}
	}
	return false
}
