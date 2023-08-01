package dom

import (
	"bufio"
	"fmt"
	"strconv"
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
			if len(elements) > 0 && elements[len(elements)-1].Type() != syntax.Empty {
				elements = append(elements, syntax.NewEmptyLine())
			}
			elements = append(elements, bels...)
		case Comments:
			bels := constructComments(block.(*CommentsBlock))
			elements = append(elements, bels...)
		case Blanks:
			bels := constructBlanks(block.(*BlanksBlock))
			elements = append(elements, bels...)
		case Unknown:
			bels := constructUnknown(block.(*UnrecognizedBlock))
			elements = append(elements, bels...)
		}
	}

	return syntax.NewDocument(elements)
}

func constructAliases(block *IPAliasesBlock) []syntax.Element {
	elements := make([]syntax.Element, 0)
	if block.origHeader != nil {
		for _, el := range block.origHeader {
			elements = append(elements, el)
		}
	} else if block.id != idNotSet || block.name != "" || block.note != "" {
		sb := strings.Builder{}

		// format block id and name prefix
		blockId := "*"
		if block.id != idNotSet {
			blockId = strconv.Itoa(block.id)
		}
		sb.WriteString(fmt.Sprintf("[%s]", blockId))
		if block.name != "" {
			sb.WriteRune(' ')
			sb.WriteString(block.name)
		}

		// format notes
		firstLine := true
		s := bufio.NewScanner(strings.NewReader(block.note))
		s.Split(iotools.LinesSplitterRespectEndNewLineFunc())
		lines := make([]string, 0)
		for {
			if ok := s.Scan(); !ok {
				break
			}
			lines = append(lines, s.Text())
		}
		for _, l := range lines {
			lt := l
			if firstLine {
				lt = fmt.Sprintf("%s - %s", sb.String(), l)
				firstLine = false
			}
			el := syntax.NewCommentsLine(lt)
			elements = append(elements, el)
		}

		// ensure header element was added even if no comments for the block
		if len(elements) == 0 {
			el := syntax.NewCommentsLine(sb.String())
			elements = append(elements, el)
		}
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
