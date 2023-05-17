package dom

import (
	"regexp"

	"github.com/0xcfff/hostsctl/hosts/syntax"
)

type parsingState int

const (
	notStarted parsingState = iota
	unrecognized
	whitespaces
	comments
	ips
)

var (
	rxBlockId = regexp.MustCompile(`^\s*\[\s*(\d+|\*)\s*\]`)
)

type parserContext struct {
	recognizedBlocks []Block

	state parsingState

	commentsList     []*syntax.CommentLine
	blanksList       []*syntax.EmptyLine
	ipsList          []syntax.Element
	unrecognizedList []syntax.Element
}

func parse(doc *syntax.Document) *Document {

	ctx := newParseContext()

	for _, el := range doc.Elements() {
		ok := ctx.tryContinueBlock(el)
		if !ok {
			ctx.finishBlock()
			ctx.startNewBlock(el)
		}
	}
	ctx.finishBlock()

	return &Document{
		originalDocument: doc,
		blocks:           ctx.recognizedBlocks,
	}
}

func (ctx *parserContext) startNewBlock(el syntax.Element) {
	switch el.Type() {
	case syntax.Unknown:
		ctx.state = unrecognized
	case syntax.Comment:
		ctx.state = comments
	case syntax.IPMapping:
		ctx.state = ips
	case syntax.Empty:
		ctx.state = whitespaces
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
			if isCommentedIPMapping(c) {
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
	case comments:
		block := newCommentsBlockFromLines(ctx.commentsList)
		ctx.recognizedBlocks = append(ctx.recognizedBlocks, block)
		ctx.commentsList = make([]*syntax.CommentLine, 0)
	case ips:
		autoId := calcNextIPsBlockAutoId(ctx)
		block := newIPAliasesBlockFromElements(ctx.commentsList, ctx.ipsList, autoId)
		ctx.recognizedBlocks = append(ctx.recognizedBlocks, block)
		ctx.commentsList = make([]*syntax.CommentLine, 0)
		ctx.ipsList = make([]syntax.Element, 0)
	case unrecognized:
		block := UnrecognizedBlock{
			lines: ctx.unrecognizedList,
		}
		ctx.recognizedBlocks = append(ctx.recognizedBlocks, &block)
		ctx.unrecognizedList = make([]syntax.Element, 0)
	default:
		panic("Unknown block type")
	}
}

func calcNextIPsBlockAutoId(ctx *parserContext) int {
	autoBlockId := 1
	for {
		found := true
		for _, b := range ctx.recognizedBlocks {
			if b.Type() == IPList {
				bb := b.(*IPAliasesBlock)
				if bb.Id() == autoBlockId {
					found = false
					break
				}
			}
		}
		if found {
			break
		}
		autoBlockId += 1
	}
	return autoBlockId
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
