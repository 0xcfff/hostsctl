package dom

import (
	"strings"

	"github.com/0xcfff/hostsctl/iptools"
	"github.com/0xcfff/hostsctl/syntax"
)

type parsingState int

const (
	notStarted parsingState = iota
	unrecognized
	whitespaces
	comments
	ips
)

type parserContext struct {
	recognizedBlocks []Block

	state parsingState

	commentsList     []*syntax.CommentLine
	blanksList       []*syntax.EmptyLine
	ipsList          []syntax.Element
	unrecognizedList []syntax.Element
}

func newParseContext() parserContext {
	return parserContext{
		recognizedBlocks: make([]Block, 0),
		state:            notStarted,
		commentsList:     make([]*syntax.CommentLine, 0),
		blanksList:       make([]*syntax.EmptyLine, 0),
		ipsList:          make([]syntax.Element, 0),
		unrecognizedList: make([]syntax.Element, 0),
	}
}

func parse(doc syntax.Document) Document {

	ctx := newParseContext()

	for _, el := range doc.Elements() {
		ok := ctx.tryContinueBlock(el)
		if !ok {
			ctx.finishBlock()
			ctx.startNewBlock(el)
		}
	}
	ctx.finishBlock()

	return Document{
		originalDocument: doc,
		blocks:           ctx.recognizedBlocks,
	}
}

func (ctx *parserContext) startNewBlock(el syntax.Element) {
	switch el.Type() {
	case syntax.Unknown:
		ctx.state = unrecognized
		break
	case syntax.Comment:
		ctx.state = comments
		break
	case syntax.IPMapping:
		ctx.state = ips
		break
	case syntax.Empty:
		ctx.state = whitespaces
		break
	default:
		panic("Unknown block type")
	}

	added := ctx.tryContinueBlock(el)
	if !added {
		panic("Should always be added correctly")
	}
}

func (ctx *parserContext) tryContinueBlock(el syntax.Element) bool {
	switch ctx.state {
	case notStarted:
		return false
	case unrecognized:
		if el.Type() == syntax.Unknown {
			ctx.unrecognizedList = append(ctx.unrecognizedList, el)
			return true
		}
		return false
	case whitespaces:
		if el.Type() == syntax.Empty {
			ctx.blanksList = append(ctx.blanksList, el.(*syntax.EmptyLine))
			return true
		}
		return false
	case comments:
		if el.Type() == syntax.Comment {
			ctx.commentsList = append(ctx.commentsList, el.(*syntax.CommentLine))
			return true
		}
		if el.Type() == syntax.IPMapping {
			ctx.ipsList = make([]syntax.Element, 0)
			ctx.ipsList = append(ctx.ipsList, el)
			ctx.state = ips
			return true
		}
		return false
	case ips:
		if el.Type() == syntax.IPMapping {
			ctx.ipsList = append(ctx.ipsList, el)
			return true
		}
		if el.Type() == syntax.Comment {
			c := el.(*syntax.CommentLine)
			parts := strings.Fields(c.CommentText())
			if len(parts) >= 2 && iptools.IsIP(parts[0]) {
				ctx.ipsList = append(ctx.ipsList, el)
				return true
			}
		}
		return false
	default:
		panic("Unknown block type")
	}
}

func (ctx *parserContext) finishBlock() {
	switch ctx.state {
	case notStarted:
		break
	case whitespaces:
		block := BlanksBlock{
			blanks: ctx.blanksList,
		}
		ctx.recognizedBlocks = append(ctx.recognizedBlocks, &block)
		ctx.blanksList = make([]*syntax.EmptyLine, 0)
		break
	case comments:
		block := CommentsBlock{
			comments: ctx.commentsList,
		}
		ctx.recognizedBlocks = append(ctx.recognizedBlocks, &block)
		ctx.commentsList = make([]*syntax.CommentLine, 0)
		break
	case ips:
		block := IPListBlock{
			header: ctx.commentsList,
			body:   ctx.ipsList,
		}
		ctx.recognizedBlocks = append(ctx.recognizedBlocks, &block)
		ctx.commentsList = make([]*syntax.CommentLine, 0)
		ctx.ipsList = make([]syntax.Element, 0)
		break
	case unrecognized:
		block := UnrecognizedBlock{
			lines: ctx.unrecognizedList,
		}
		ctx.recognizedBlocks = append(ctx.recognizedBlocks, &block)
		ctx.unrecognizedList = make([]syntax.Element, 0)
		break
	default:
		panic("Unknown block type")
	}
}
