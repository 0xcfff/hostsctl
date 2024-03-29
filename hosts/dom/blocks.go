package dom

import (
	"github.com/0xcfff/hostsctl/hosts/syntax"
	"golang.org/x/exp/slices"
)

type BlockType int

const (
	Unknown BlockType = iota
	Comments
	Blanks
	IPList
)

type Block interface {
	Type() BlockType
	dirty() bool
}

// Block of unrecognized lines
type UnrecognizedBlock struct {
	lines []syntax.Element
}

func (blk *UnrecognizedBlock) Type() BlockType {
	return Unknown
}

func (blk *UnrecognizedBlock) dirty() bool {
	return false
}

func (blk *UnrecognizedBlock) BodyElements() []syntax.Element {
	return slices.Clone(blk.lines)
}

// Sequence of blank lines
type BlanksBlock struct {
	blanks []*syntax.EmptyLine
}

func (blk *BlanksBlock) Type() BlockType {
	return Blanks
}

func (blk *BlanksBlock) dirty() bool {
	return false
}

func (blk *BlanksBlock) LinesCount() int {
	return len(blk.blanks)
}
