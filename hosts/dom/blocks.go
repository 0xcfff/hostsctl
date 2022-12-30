package dom

import (
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
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
}

// Block of unrecognized lines
type UnrecognizedBlock struct {
	lines []syntax.Element
}

func (blk *UnrecognizedBlock) Type() BlockType {
	return Unknown
}

func (blk *UnrecognizedBlock) BodyElements() []syntax.Element {
	list := make([]syntax.Element, 0, len(blk.lines))
	return list
}

// Sequence of comments
type CommentsBlock struct {
	comments []*syntax.CommentLine
}

func (blk *CommentsBlock) Type() BlockType {
	return Comments
}

func (blk *CommentsBlock) LinesCount() int {
	return len(blk.comments)
}

func (blk *CommentsBlock) CommentLines() []string {
	lines := make([]string, 0, len(blk.comments))
	for _, l := range blk.comments {
		lines = append(lines, l.CommentText())
	}
	return lines
}

// Sequence of blank lines
type BlanksBlock struct {
	blanks []*syntax.EmptyLine
}

func (blk *BlanksBlock) Type() BlockType {
	return Blanks
}

func (blk *BlanksBlock) LinesCount() int {
	return len(blk.blanks)
}

// Block of IPs
type IPListBlock struct {
	header  []*syntax.CommentLine
	body    []syntax.Element
	id      int
	name    string
	comment []string
}

func (blk *IPListBlock) Type() BlockType {
	return IPList
}

func (blk *IPListBlock) Id() int {
	return blk.id
}
func (blk *IPListBlock) Name() string {
	return blk.name
}
func (blk *IPListBlock) CommentLines() []string {
	lines := make([]string, 0, len(blk.comment))
	lines = append(lines, blk.comment...)
	return lines
}

func (blk *IPListBlock) CommentText() string {
	sb := &strings.Builder{}
	first := true
	for _, s := range blk.comment {
		if first {
			first = false
		} else {
			sb.WriteString("\n")
		}
		sb.WriteString(s)
	}
	return sb.String()
}

func (blk *IPListBlock) HasHeader() bool {
	return len(blk.header) > 0
}

func (blk *IPListBlock) HeaderCommentLines() []string {
	lines := make([]string, 0, len(blk.header))
	for _, l := range blk.header {
		lines = append(lines, l.CommentText())
	}
	return lines
}

func (blk *IPListBlock) BodyElements() []syntax.Element {
	list := make([]syntax.Element, len(blk.body))
	copy(list, blk.body)
	return list
}
