package dom

import (
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
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
