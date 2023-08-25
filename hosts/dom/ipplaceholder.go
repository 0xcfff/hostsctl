package dom

import "github.com/0xcfff/hostsctl/hosts/syntax"

type IPAliasesPlaceholder struct {
	origElement *syntax.CommentLine
	changed     bool
}

func (blk *IPAliasesPlaceholder) Type() IPAliasesBlockElementType {
	return Placeholder
}

func (blk *IPAliasesPlaceholder) ClearFormatting() {
	blk.origElement = nil
	blk.changed = true
}

func (blk *IPAliasesPlaceholder) dirty() bool {
	return blk.changed
}

func newIPAliasesPlaceholderFromCommentElement(el *syntax.CommentLine) *IPAliasesPlaceholder {
	return &IPAliasesPlaceholder{
		origElement: el,
		changed:     false,
	}
}

func NewIPAliasesPlaceholder() *IPAliasesPlaceholder {
	return &IPAliasesPlaceholder{
		origElement: nil,
		changed:     true,
	}
}
