package hosts

import (
	"fmt"
	"strings"

	"github.com/0xcfff/dnspipe/backend/hosts"
	"golang.org/x/exp/slices"
)

type IPModel struct {
	IP      string   `json:"ip"`
	Aliases []string `json:"aliases"`
	Comment string   `json:"comment,omitempty"`
	Source  string   `json:"source,omitempty"`
}

type IPGroupModel struct {
	Name string    `json:"name"`
	IPs  []IPModel `json:"ips"`
}

type IPGrouping int

const (
	GrpOriginal IPGrouping = iota
	GrpUngroup  IPGrouping = iota
	GrpGroup    IPGrouping = iota
)

func NewIPModels(hostsContent *hosts.HostsFileContent, grouping IPGrouping) []*IPModel {
	var result []*IPModel
	if len(hostsContent.IPRecords) > 0 {
		result = append(result, convertIPs(hostsContent.IPRecords, "(default)", grouping)...)
	}
	for idx, sb := range hostsContent.SyncBlocks {
		if sb.Data != nil && len(sb.Data.IPRecords) > 0 {
			result = append(result, convertIPs(sb.Data.IPRecords, fmt.Sprintf("TBD group %v", idx), grouping)...)
		}
	}
	return result
}

func convertIPs(ips []*hosts.IPRecord, source string, grouping IPGrouping) []*IPModel {
	switch grouping {
	case GrpUngroup:
		return ungroupAndConvert(ips, source)
	case GrpGroup:
		return groupAndConvert(ips, source)
	case GrpOriginal:
		return convertOnly(ips, source)
	default:
		panic("unknown grouping specified")
	}
}

func ungroupAndConvert(ips []*hosts.IPRecord, source string) []*IPModel {
	result := make([]*IPModel, 0)

	for _, r := range ips {
		for _, al := range r.Aliases {
			ip := &IPModel{
				IP:      r.IP,
				Aliases: []string{al},
				Comment: r.Notes,
				Source:  source,
			}
			result = append(result, ip)
		}
	}
	return result
}

func groupAndConvert(ips []*hosts.IPRecord, source string) []*IPModel {
	result := make([]*IPModel, 0)
	ipsMap := make(map[string]*IPModel)
	ipsComments := make(map[string][]string)

	for _, r := range ips {
		ip, ok := ipsMap[r.IP]
		if !ok {
			ip := &IPModel{
				IP:      r.IP,
				Source:  source,
				Comment: r.Notes,
			}
			result = append(result, ip)
			ipsMap[r.IP] = ip

			comments := make([]string, 0)
			if r.Notes != "" {
				comments = append(comments, r.Notes)
			}
			ipsComments[r.IP] = comments
		}
		ip.Aliases = append(ip.Aliases, r.Aliases...)

		comments := ipsComments[ip.IP]
		if ip.Comment != "" && !slices.Contains(comments, ip.Comment) {
			comments = append(comments, r.Notes)
			ipsComments[ip.IP] = comments
			ip.Comment = strings.Join(comments, ", ")
		}
	}
	return result
}

func convertOnly(ips []*hosts.IPRecord, source string) []*IPModel {
	result := make([]*IPModel, 0)

	for _, r := range ips {
		ip := &IPModel{
			IP:      r.IP,
			Comment: r.Notes,
			Aliases: r.Aliases,
			Source:  source,
		}
		result = append(result, ip)
	}
	return result
}
