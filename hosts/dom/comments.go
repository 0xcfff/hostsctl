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
}

func (blk *CommentsBlock) Type() BlockType {
	return Comments
}

func (blk *CommentsBlock) dirty() bool {
	return blk.origComments == nil
}

func (blk *CommentsBlock) CommentsText() string {
	return blk.commentsText
}

func (blk *CommentsBlock) SetCommentsText(commentsText string) {
	if strings.Compare(blk.commentsText, commentsText) != 0 {
		blk.commentsText = commentsText
		blk.origComments = nil
	}
}

// Creates new comments block
func NewCommentsBlock(commentsText string) *CommentsBlock {
	blk := &CommentsBlock{
		commentsText: commentsText,
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
	}

	return blk
}
