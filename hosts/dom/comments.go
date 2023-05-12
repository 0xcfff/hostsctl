package dom

import (
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"github.com/0xcfff/hostsctl/iotools"
)

// Sequence of comments
type CommentsBlock struct {
	origComments []*syntax.CommentLine
	commentsText string
	dirty        bool
}

func (blk *CommentsBlock) Type() BlockType {
	return Comments
}

func (blk *CommentsBlock) Dirty() bool {
	return blk.dirty
}

func (blk *CommentsBlock) CommentsText() string {
	return blk.commentsText
}

func (blk *CommentsBlock) SetCommentsText(commentsText string) {
	if strings.Compare(blk.commentsText, commentsText) != 0 {
		blk.commentsText = commentsText
		blk.dirty = true
	}
}

// Creates new comments block
func NewCommentsBlock(commentsText string) *CommentsBlock {
	blk := &CommentsBlock{
		commentsText: commentsText,
		dirty:        true,
	}
	return blk
}

func newCommentsBlockFromLines(comments []*syntax.CommentLine) *CommentsBlock {
	// Aggregate comments lines into one text string
	newLine := iotools.OSDependendNewLine()
	sb := &strings.Builder{}
	first := true
	for _, s := range comments {
		if first {
			first = false
		} else {
			sb.WriteString(newLine)
		}
		sb.WriteString(s.CommentText())
	}

	blk := &CommentsBlock{
		origComments: comments,
		commentsText: sb.String(),
		dirty:        false,
	}

	return blk
}
