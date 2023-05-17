package host

import (
	"strings"

	"github.com/0xcfff/hostsctl/hosts/dom"
	"golang.org/x/exp/slices"
)

type HostModel struct {
	IP      string         `json:"ip"                yaml:"ip"`
	Hosts   []string       `json:"aliases"           yaml:"aliases"`
	Comment string         `json:"comment,omitempty" yaml:"comment,omitempty"`
	Group   HostGroupModel `json:"group,omitempty"   yaml:"group,omitempty"`
}

type HostGroupModel struct {
	Id      int    `json:"id"             yaml:"id"`
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Comment string `json:"-"              yaml:"-"`
}

type IPGrouping int

const (
	GrpOriginal IPGrouping = iota
	GrpUngroup  IPGrouping = iota
	GrpGroup    IPGrouping = iota
)

func NewHostModels(doc *dom.Document, grouping IPGrouping) []*HostModel {
	var result []*HostModel = make([]*HostModel, 0)

	for _, block := range doc.Blocks() {
		if block.Type() == dom.IPList {
			ipsBlock := block.(*dom.IPAliasesBlock)
			result = append(result, convertIPs(ipsBlock, grouping)...)
		}
	}
	return result
}

func convertIPs(ips *dom.IPAliasesBlock, grouping IPGrouping) []*HostModel {
	switch grouping {
	case GrpUngroup:
		return ungroupAndConvert(ips)
	case GrpGroup:
		return groupAndConvert(ips)
	case GrpOriginal:
		return convertOnly(ips)
	default:
		panic("unknown grouping specified")
	}
}

func ungroupAndConvert(ips *dom.IPAliasesBlock) []*HostModel {
	result := make([]*HostModel, 0)

	group := HostGroupModel{
		Id:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Note(),
	}

	for _, r := range ips.Entries() {
		for _, al := range r.Aliases() {
			ip := &HostModel{
				IP:      r.IP(),
				Hosts:   []string{al},
				Comment: r.Note(),
				Group:   group,
			}
			result = append(result, ip)
		}
	}
	return result
}

func groupAndConvert(ips *dom.IPAliasesBlock) []*HostModel {
	result := make([]*HostModel, 0)
	ipsMap := make(map[string]*HostModel)
	ipsComments := make(map[string][]string)

	group := HostGroupModel{
		Id:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Note(),
	}

	for _, r := range ips.Entries() {
		ip, ok := ipsMap[r.IP()]
		if !ok {
			ip = &HostModel{
				IP:      r.IP(),
				Group:   group,
				Comment: r.Note(),
			}
			result = append(result, ip)
			ipsMap[r.IP()] = ip

			comments := make([]string, 0)
			if r.Note() != "" {
				comments = append(comments, r.Note())
			}
			ipsComments[r.IP()] = comments
		}
		ip.Hosts = append(ip.Hosts, r.Aliases()...)

		comments := ipsComments[ip.IP]
		if ip.Comment != "" && !slices.Contains(comments, ip.Comment) {
			comments = append(comments, r.Note())
			ipsComments[ip.IP] = comments
			ip.Comment = strings.Join(comments, ", ")
		}
	}
	return result
}

func convertOnly(ips *dom.IPAliasesBlock) []*HostModel {
	result := make([]*HostModel, 0)

	group := HostGroupModel{
		Id:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Note(),
	}

	for _, r := range ips.Entries() {
		ip := &HostModel{
			IP:      r.IP(),
			Comment: r.Note(),
			Hosts:   r.Aliases(),
			Group:   group,
		}
		result = append(result, ip)
	}
	return result
}
