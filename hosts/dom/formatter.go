package dom

import (
	"bufio"
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"github.com/0xcfff/hostsctl/iotools"
)

func constructSyntax(doc *Document) *syntax.Document {
	elements := make([]syntax.Element, 0)

	for _, block := range doc.blocks {
		switch block.Type() {
		case IPList:
			bels := constructAliases(block.(*IPAliasesBlock))
			elements = append(elements, bels...)
			break
		case Comments:
			bels := constructComments(block.(*CommentsBlock))
			elements = append(elements, bels...)
			break
		case Blanks:
			bels := constructBlanks(block.(*BlanksBlock))
			elements = append(elements, bels...)
			break
		case Unknown:
			bels := constructUnknown(block.(*UnrecognizedBlock))
			elements = append(elements, bels...)
			break
		}
	}

	return syntax.NewDocument(elements)
}

func constructAliases(block *IPAliasesBlock) []syntax.Element {
	elements := make([]syntax.Element, 0)
	for _, el := range block.origHeader {
		elements = append(elements, el)
	}
	for _, el := range block.entries {
		ipAlias := el.origElement
		if ipAlias == nil {
			ipAlias = syntax.NewIPMappingLine(el.ip, el.aliases, el.note)
			el.origElement = ipAlias
		}
		elements = append(elements, ipAlias)
	}
	return elements
}

func constructComments(block *CommentsBlock) []syntax.Element {
	elements := make([]syntax.Element, 0)
	if block.origComments != nil {
		for _, el := range block.origComments {
			elements = append(elements, el)
		}
	} else {
		s := bufio.NewScanner(strings.NewReader(block.commentsText))
		s.Split(iotools.LinesSplitterRespectEndNewLineFunc())
		lines := make([]string, 0)
		for {
			if ok := s.Scan(); !ok {
				break
			}
			lines = append(lines, s.Text())
		}
		for _, l := range lines {
			el := syntax.NewCommentsLine(l)
			elements = append(elements, el)
		}
	}
	return elements
}

func constructBlanks(block *BlanksBlock) []syntax.Element {
	elements := make([]syntax.Element, 0)
	for _, el := range block.blanks {
		elements = append(elements, el)
	}
	return elements
}

func constructUnknown(block *UnrecognizedBlock) []syntax.Element {
	elements := make([]syntax.Element, 0)
	for _, el := range block.lines {
		if !el.HasPreformattedText() {
			panic("Can not format unrecognized block, logic should never go here")
		}
		elements = append(elements, el)
	}
	return elements
}
