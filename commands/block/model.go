package block

import (
	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/0xcfff/hostsctl/iptools"
)

type BlockModel struct {
	ID                 int    `json:"id"                yaml:"id"`
	Name               string `json:"name"              yaml:"name"`
	Comment            string `json:"comment,omitempty" yaml:"comment,omitempty"`
	AliasesCount       int    `json:"count,omitempty"   yaml:"count,omitempty"`
	SystemAliasesCount int    `json:"-"   				yaml:"-"`
}

func NewBlocksModels(doc *dom.Document) []*BlockModel {
	var result []*BlockModel = make([]*BlockModel, 0)

	for _, block := range doc.Blocks() {
		if block.Type() == dom.IPList {
			ipsBlock := block.(*dom.IPAliasesBlock)
			result = append(result, convertIPs(ipsBlock))
		}
	}
	return result
}

func convertIPs(ips *dom.IPAliasesBlock) *BlockModel {
	block := &BlockModel{
		ID:      ips.Id(),
		Name:    ips.Name(),
		Comment: ips.Note(),
	}

	cntAll := 0
	cntSys := 0

	for _, ip := range ips.Entries() {
		for _, a := range ip.Aliases() {
			cntAll += 1
			if iptools.IsSystemAlias(ip.IP(), a) {
				cntSys += 1
			}
		}
	}

	block.AliasesCount = cntAll
	block.SystemAliasesCount = cntSys

	return block
}
