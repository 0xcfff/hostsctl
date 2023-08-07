package alias

import (
	"strings"

	"github.com/0xcfff/hostsctl/hosts/dom"
	"golang.org/x/exp/slices"
)

type AliasModel struct {
	IP      string          `json:"ip"                yaml:"ip"`
	Aliases []string        `json:"aliases"           yaml:"aliases"`
	Comment string          `json:"comment,omitempty" yaml:"comment,omitempty"`
	Block   AliasBlockModel `json:"block,omitempty"   yaml:"block,omitempty"`
}

type AliasBlockModel struct {
	Id      int    `json:"id"             yaml:"id"`
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Comment string `json:"-"              yaml:"-"`
}

type IPGrouping int

const (
	GrpRaw     IPGrouping = iota
	GrpUngroup IPGrouping = iota
	GrpGroup   IPGrouping = iota
)

func NewHostModels(doc *dom.Document, grouping IPGrouping) []*AliasModel {
	var result []*AliasModel = make([]*AliasModel, 0)

	for _, block := range doc.Blocks() {
		if block.Type() == dom.IPList {
			ipsBlock := block.(*dom.IPAliasesBlock)
			result = append(result, convertIPs(ipsBlock, grouping)...)
		}
	}
	return result
}

func convertIPs(ips *dom.IPAliasesBlock, grouping IPGrouping) []*AliasModel {
	switch grouping {
	case GrpUngroup:
		return ungroupAndConvert(ips)
	case GrpGroup:
		return groupAndConvert(ips)
	case GrpRaw:
		return convertOnly(ips)
	default:
		panic("unknown grouping specified")
	}
}

func ungroupAndConvert(ips *dom.IPAliasesBlock) []*AliasModel {
	result := make([]*AliasModel, 0)

	group := AliasBlockModel{
		Id:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Note(),
	}

	for _, r := range ips.Entries() {
		for _, al := range r.Aliases() {
			ip := &AliasModel{
				IP:      r.IP(),
				Aliases: []string{al},
				Comment: r.Note(),
				Block:   group,
			}
			result = append(result, ip)
		}
	}
	return result
}

func groupAndConvert(ips *dom.IPAliasesBlock) []*AliasModel {
	result := make([]*AliasModel, 0)
	ipsMap := make(map[string]*AliasModel)
	ipsComments := make(map[string][]string)

	group := AliasBlockModel{
		Id:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Note(),
	}

	for _, r := range ips.Entries() {
		ip, ok := ipsMap[r.IP()]
		if !ok {
			ip = &AliasModel{
				IP:      r.IP(),
				Block:   group,
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
		ip.Aliases = append(ip.Aliases, r.Aliases()...)

		comments := ipsComments[ip.IP]
		if ip.Comment != "" && !slices.Contains(comments, ip.Comment) {
			comments = append(comments, r.Note())
			ipsComments[ip.IP] = comments
			ip.Comment = strings.Join(comments, ", ")
		}
	}
	return result
}

func convertOnly(ips *dom.IPAliasesBlock) []*AliasModel {
	result := make([]*AliasModel, 0)

	group := AliasBlockModel{
		Id:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Note(),
	}

	for _, r := range ips.Entries() {
		ip := &AliasModel{
			IP:      r.IP(),
			Comment: r.Note(),
			Aliases: r.Aliases(),
			Block:   group,
		}
		result = append(result, ip)
	}
	return result
}
