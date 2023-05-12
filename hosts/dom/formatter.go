package dom

import (
	"github.com/0xcfff/hostsctl/hosts/syntax"
)

type FormatMode int

type formattingContext struct {
	formattedElements []syntax.Element

	// TODO: uncomment to implement minimal formatting for newly added elements
	// commentsFormat struct {
	// }
	// aliasesFormat struct {
	// }
}

func format(doc *Document) *syntax.Document {
	ctx := newFormattingContext()

	for _, block := range doc.blocks {
		switch block.Type() {
		case IPList:
			formatAliases(ctx, block.(*IPListBlock))
			break
		case Comments:
			formatComments(ctx, block.(*CommentsBlock))
			break
		case Blanks:
			formatBlanks(ctx, block.(*BlanksBlock))
			break
		case Unknown:
			formatUnknown(ctx, block.(*UnrecognizedBlock))
			break
		}
	}

	return syntax.NewDocument(ctx.formattedElements)
}

func formatAliases(ctx *formattingContext, block *IPListBlock) {
	for _, el := range block.header {
		ctx.appendElement(el)
	}
	for _, el := range block.body {
		ctx.appendElement(el)
	}
}

func formatComments(ctx *formattingContext, block *CommentsBlock) {
	for _, el := range block.origComments {
		// TODO: Add formatting to make sure newly added IPs are alligned with previously added ones
		ctx.appendElement(el)
	}
}

func formatBlanks(ctx *formattingContext, block *BlanksBlock) {
	for _, el := range block.blanks {
		ctx.appendElement(el)
	}
}

func formatUnknown(ctx *formattingContext, block *UnrecognizedBlock) {
	for _, el := range block.lines {
		if !el.HasPreformattedText() {
			panic("Can not format unrecognized block, logic should never go here")
		}
		ctx.appendElement(el)
	}
}

func newFormattingContext() *formattingContext {
	ctx := formattingContext{
		formattedElements: make([]syntax.Element, 0),
	}
	return &ctx
}

func (ctx *formattingContext) appendElement(el syntax.Element) {
	ctx.formattedElements = append(ctx.formattedElements, el)
}
