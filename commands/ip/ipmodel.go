package ip

import (
	"fmt"
	"strings"

	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/hosts/syntax"
	"golang.org/x/exp/slices"
)

type IPModel struct {
	IP      string   `json:"ip"`
	Aliases []string `json:"aliases"`
	Comment string   `json:"comment,omitempty"`
	Group   string   `json:"group,omitempty"`
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

func NewIPModels(doc *dom.Document, grouping IPGrouping) []*IPModel {
	var result []*IPModel = make([]*IPModel, 0)

	ipsBlockIdx := 0
	for _, block := range doc.Blocks() {
		if block.Type() == dom.IPList {
			ipsBlock := block.(*dom.IPListBlock)
			blockName := fmt.Sprintf("TBD group %v", ipsBlockIdx)
			result = append(result, convertIPs(ipsBlock, blockName, grouping)...)
			ipsBlockIdx += 1
		}
	}
	return result
}

func convertIPs(ips *dom.IPListBlock, source string, grouping IPGrouping) []*IPModel {
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

func ungroupAndConvert(ips *dom.IPListBlock, source string) []*IPModel {
	result := make([]*IPModel, 0)

	for _, r := range ips.BodyElements() {
		if r.Type() == syntax.IPMapping {
			rr := r.(*syntax.IPMappingLine)
			for _, al := range rr.DomainNames() {
				ip := &IPModel{
					IP:      rr.IPAddress(),
					Aliases: []string{al},
					Comment: rr.CommentText(),
					Group:   source,
				}
				result = append(result, ip)
			}
		}
	}
	return result
}

func groupAndConvert(ips *dom.IPListBlock, source string) []*IPModel {
	result := make([]*IPModel, 0)
	ipsMap := make(map[string]*IPModel)
	ipsComments := make(map[string][]string)

	for _, r := range ips.BodyElements() {
		if r.Type() == syntax.IPMapping {
			rr := r.(*syntax.IPMappingLine)
			ip, ok := ipsMap[rr.IPAddress()]
			if !ok {
				ip := &IPModel{
					IP:      rr.IPAddress(),
					Group:   source,
					Comment: rr.CommentText(),
				}
				result = append(result, ip)
				ipsMap[rr.IPAddress()] = ip

				comments := make([]string, 0)
				if rr.CommentText() != "" {
					comments = append(comments, rr.CommentText())
				}
				ipsComments[rr.IPAddress()] = comments
			}
			ip.Aliases = append(ip.Aliases, rr.DomainNames()...)

			comments := ipsComments[ip.IP]
			if ip.Comment != "" && !slices.Contains(comments, ip.Comment) {
				comments = append(comments, rr.CommentText())
				ipsComments[ip.IP] = comments
				ip.Comment = strings.Join(comments, ", ")
			}
		}
	}
	return result
}

func convertOnly(ips *dom.IPListBlock, source string) []*IPModel {
	result := make([]*IPModel, 0)

	for _, r := range ips.BodyElements() {
		if r.Type() == syntax.IPMapping {
			rr := r.(*syntax.IPMappingLine)
			ip := &IPModel{
				IP:      rr.IPAddress(),
				Comment: rr.CommentText(),
				Aliases: rr.DomainNames(),
				Group:   source,
			}
			result = append(result, ip)
		}
	}
	return result
}
