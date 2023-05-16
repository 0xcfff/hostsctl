package dom

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"github.com/0xcfff/hostsctl/iotools"
	"github.com/0xcfff/hostsctl/iptools"
	"golang.org/x/exp/slices"
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
	case comments:
		block := newCommentsBlockFromLines(ctx.commentsList)
		ctx.recognizedBlocks = append(ctx.recognizedBlocks, block)
		ctx.commentsList = make([]*syntax.CommentLine, 0)
	case ips:
		block := IPListBlock{
			header: ctx.commentsList,
			body:   ctx.ipsList,
		}
		parseNFillIPsBlockValues(ctx, &block)
		ctx.recognizedBlocks = append(ctx.recognizedBlocks, &block)
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

func parseNFillIPsBlockValues(ctx *parserContext, block *IPListBlock) {

	var blockId int = idNotSet
	var autoBlockId int = idNotSet
	var blockName string
	var commentsText string

	if len(block.header) > 0 {
		headerLine := block.header[0].CommentText()

		// try extract block ID
		matches := rxBlockId.FindAllStringSubmatch(headerLine, -1)
		if matches != nil {
			match := matches[0]
			fullstr := match[0]
			idstr := match[1]
			if idstr == "*" {
				blockId = idNotSet
			} else {
				var err error
				blockId, err = strconv.Atoi(idstr)
				if err != nil {
					blockId = idNotSet
				}
			}
			headerLine = strings.TrimSpace(headerLine[len(fullstr):])
		}

		// try extract block name
		parts := strings.Fields(headerLine)
		if len(parts) == 1 {
			blockName = parts[0]
			headerLine = ""
		} else if len(parts) > 1 {
			dividers := []string{"-", ":", "|", "*", "#"}
			div := parts[1]
			if slices.Contains(dividers, div) {
				blockName = parts[0]
				noName := headerLine[len(blockName):]
				divIdx := strings.Index(noName, div)
				headerLine = strings.TrimSpace(noName[divIdx+len(div):])
			}
		}

		var blockComment []string
		blockComment = append(blockComment, headerLine)
		for _, line := range block.header[1:] {
			blockComment = append(blockComment, line.CommentText())
		}

		// Construct comment text
		newLine := iotools.OSDependendNewLine()
		sb := &strings.Builder{}
		first := true
		for _, s := range blockComment {
			if first {
				first = false
			} else {
				sb.WriteString(newLine)
			}
			sb.WriteString(s)
		}
		commentsText = sb.String()
	}

	if blockId == idNotSet {
		autoBlockId = 1
		for {
			found := true
			for _, b := range ctx.recognizedBlocks {
				if b.Type() == IPList {
					bb := b.(*IPListBlock)
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
	}

	block.id = blockId
	block.autoId = autoBlockId
	block.name = blockName
	block.commentsText = commentsText
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
