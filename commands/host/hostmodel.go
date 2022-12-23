package host

import (
	"fmt"
	"strings"

	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/hosts/syntax"
	"golang.org/x/exp/slices"
)

type HostModel struct {
	IP      string   `json:"ip"`
	Hosts   []string `json:"hosts"`
	Comment string   `json:"comment,omitempty"`
	Group   string   `json:"group,omitempty"`
}

type HostGroupModel struct {
	Name string      `json:"name"`
	IPs  []HostModel `json:"ips"`
}

type IPGrouping int

const (
	GrpOriginal IPGrouping = iota
	GrpUngroup  IPGrouping = iota
	GrpGroup    IPGrouping = iota
)

func NewHostModels(doc *dom.Document, grouping IPGrouping) []*HostModel {
	var result []*HostModel = make([]*HostModel, 0)

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

func convertIPs(ips *dom.IPListBlock, source string, grouping IPGrouping) []*HostModel {
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

func ungroupAndConvert(ips *dom.IPListBlock, source string) []*HostModel {
	result := make([]*HostModel, 0)

	for _, r := range ips.BodyElements() {
		if r.Type() == syntax.IPMapping {
			rr := r.(*syntax.IPMappingLine)
			for _, al := range rr.DomainNames() {
				ip := &HostModel{
					IP:      rr.IPAddress(),
					Hosts:   []string{al},
					Comment: rr.CommentText(),
					Group:   source,
				}
				result = append(result, ip)
			}
		}
	}
	return result
}

func groupAndConvert(ips *dom.IPListBlock, source string) []*HostModel {
	result := make([]*HostModel, 0)
	ipsMap := make(map[string]*HostModel)
	ipsComments := make(map[string][]string)

	for _, r := range ips.BodyElements() {
		if r.Type() == syntax.IPMapping {
			rr := r.(*syntax.IPMappingLine)
			ip, ok := ipsMap[rr.IPAddress()]
			if !ok {
				ip := &HostModel{
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
			ip.Hosts = append(ip.Hosts, rr.DomainNames()...)

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

func convertOnly(ips *dom.IPListBlock, source string) []*HostModel {
	result := make([]*HostModel, 0)

	for _, r := range ips.BodyElements() {
		if r.Type() == syntax.IPMapping {
			rr := r.(*syntax.IPMappingLine)
			ip := &HostModel{
				IP:      rr.IPAddress(),
				Comment: rr.CommentText(),
				Hosts:   rr.DomainNames(),
				Group:   source,
			}
			result = append(result, ip)
		}
	}
	return result
}
