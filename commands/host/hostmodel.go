package host

import (
	"strings"

	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/hosts/syntax"
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
			ipsBlock := block.(*dom.IPListBlock)
			result = append(result, convertIPs(ipsBlock, grouping)...)
		}
	}
	return result
}

func convertIPs(ips *dom.IPListBlock, grouping IPGrouping) []*HostModel {
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

func ungroupAndConvert(ips *dom.IPListBlock) []*HostModel {
	result := make([]*HostModel, 0)

	group := HostGroupModel{
		Id:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Comment(),
	}

	for _, r := range ips.BodyElements() {
		if r.Type() == syntax.IPMapping {
			rr := r.(*syntax.IPMappingLine)
			for _, al := range rr.DomainNames() {
				ip := &HostModel{
					IP:      rr.IPAddress(),
					Hosts:   []string{al},
					Comment: rr.CommentText(),
					Group:   group,
				}
				result = append(result, ip)
			}
		}
	}
	return result
}

func groupAndConvert(ips *dom.IPListBlock) []*HostModel {
	result := make([]*HostModel, 0)
	ipsMap := make(map[string]*HostModel)
	ipsComments := make(map[string][]string)

	group := HostGroupModel{
		Id:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Comment(),
	}

	for _, r := range ips.BodyElements() {
		if r.Type() == syntax.IPMapping {
			rr := r.(*syntax.IPMappingLine)
			ip, ok := ipsMap[rr.IPAddress()]
			if !ok {
				ip = &HostModel{
					IP:      rr.IPAddress(),
					Group:   group,
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

func convertOnly(ips *dom.IPListBlock) []*HostModel {
	result := make([]*HostModel, 0)

	group := HostGroupModel{
		Id:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Comment(),
	}

	for _, r := range ips.BodyElements() {
		if r.Type() == syntax.IPMapping {
			rr := r.(*syntax.IPMappingLine)
			ip := &HostModel{
				IP:      rr.IPAddress(),
				Comment: rr.CommentText(),
				Hosts:   rr.DomainNames(),
				Group:   group,
			}
			result = append(result, ip)
		}
	}
	return result
}
