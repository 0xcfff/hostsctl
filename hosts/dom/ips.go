package dom

import (
	"github.com/0xcfff/hostsctl/hosts/syntax"
)

const idNotSet = -1

// Block of IPs
type IPListBlock struct {
	header       []*syntax.CommentLine
	body         []syntax.Element
	id           int
	autoId       int
	name         string
	commentsText string
	changed      bool
}

func (blk *IPListBlock) Type() BlockType {
	return IPList
}

func (blk *IPListBlock) dirty() bool {
	return blk.changed
}

func (blk *IPListBlock) Id() int {
	result := blk.id
	if blk.id == idNotSet {
		result = blk.autoId
	}
	return result
}

func (blk *IPListBlock) SetId(id int) {
	blk.id = id
	blk.header = nil
	blk.changed = true
}

// Returns true if real ID value is set,
// otherwise if ID is auto generated, then returns false
func (blk *IPListBlock) IdSet() bool {
	return blk.id <= idNotSet
}

func (blk *IPListBlock) Name() string {
	return blk.name
}

func (blk *IPListBlock) SetName(name string) {
	blk.name = name
	blk.header = nil
	blk.changed = true
}

func (blk *IPListBlock) Comment() string {
	return blk.commentsText
}

func (blk *IPListBlock) SetComment(comment string) {
	blk.commentsText = comment
	blk.header = nil
	blk.changed = true
}

func (blk *IPListBlock) BodyElements() []syntax.Element {
	list := make([]syntax.Element, len(blk.body))
	copy(list, blk.body)
	return list
}
