package dom

import (
	"strconv"
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"golang.org/x/exp/slices"
)

// Container object for dom blocks found during parsing of the original document
type Document struct {
	originalDocument *syntax.Document
	blocks           []Block
}

func (doc *Document) BlocksCount() int {
	return len(doc.blocks)
}

func (doc *Document) Blocks() []Block {
	blocks := make([]Block, len(doc.blocks))
	copy(blocks, doc.blocks)
	return blocks
}

// Finds IPs block by ID
func (doc *Document) IPsBlockById(id int) *IPAliasesBlock {
	predicate := func(b *IPAliasesBlock) bool { return b.Id() == id }
	return findBlockByPredicate(doc.blocks, predicate)
}

// Finds IPs block by name
func (doc *Document) IPsBlockByName(name string) *IPAliasesBlock {
	predicateCS := func(b *IPAliasesBlock) bool { return strings.Compare(b.Name(), name) == 0 }
	predicateCI := func(b *IPAliasesBlock) bool { return strings.EqualFold(b.Name(), name) }
	result := findBlockByPredicate(doc.blocks, predicateCS)
	if result == nil {
		result = findBlockByPredicate(doc.blocks, predicateCI)
	}
	return result
}

func (doc *Document) IPsBlockByIdOrName(idOrName string) *IPAliasesBlock {
	var ipsBlock *IPAliasesBlock
	if id, err := strconv.Atoi(idOrName); err == nil {
		ipsBlock = doc.IPsBlockById(id)
	} else {
		ipsBlock = doc.IPsBlockByName(idOrName)
	}
	return ipsBlock
}

func (doc *Document) IPBlocksByIdentifiers(id int, name string) []*IPAliasesBlock {
	selectedBlocks := make([]*IPAliasesBlock, 0)
	allBlocks := doc.IPBlocks()
	for _, b := range allBlocks {
		added := false
		if id != idNotSet && b.IdSet() && b.Id() == id {
			selectedBlocks = append(selectedBlocks, b)
			added = true
		}
		if !added && name != "" && strings.EqualFold(name, b.Name()) {
			selectedBlocks = append(selectedBlocks, b)
			added = true
		}
	}

	return selectedBlocks
}

func (doc *Document) IPBlocks() []*IPAliasesBlock {
	blocks := make([]*IPAliasesBlock, 0)
	for _, blk := range doc.blocks {
		if blk.Type() == IPList {
			blocks = append(blocks, blk.(*IPAliasesBlock))
		}
	}
	return blocks
}

func (doc *Document) AddBlock(block Block) {
	doc.blocks = append(doc.blocks, block)
}

func (doc *Document) DeleteBlock(block Block) {
	idx := slices.Index(doc.blocks, block)
	if idx != -1 {
		doc.blocks = slices.Delete(doc.blocks, idx, idx+1)
	}
}

func findBlockByPredicate[B any](blocks []Block, match func(block B) bool) B {
	var result B
	for _, blk := range blocks {
		if tblk, ok := blk.(B); ok {
			if match(tblk) {
				result = tblk
			}
		}
	}
	return result
}

func NewDocument(doc *syntax.Document) *Document {
	return parse(doc)
}

func NewEmptyDocument() *Document {
	return &Document{}
}
