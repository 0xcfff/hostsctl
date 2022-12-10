package dom

import (
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

func NewDocument(doc *syntax.Document) *Document {
	return parse(doc)
}

func NewEmptyDocument() *Document {
	return &Document{}
}
