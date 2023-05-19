package dom

import (
	"strconv"
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"github.com/0xcfff/hostsctl/iotools"
	"golang.org/x/exp/slices"
)

const idNotSet = -1

// Block of IPs
type IPAliasesBlock struct {
	origHeader []*syntax.CommentLine
	id         int
	autoId     int
	name       string
	note       string
	entries    []*IPAliasesEntry
	changed    bool
}

func (blk *IPAliasesBlock) Type() BlockType {
	return IPList
}

func (blk *IPAliasesBlock) dirty() bool {
	return blk.changed
}

func (blk *IPAliasesBlock) Id() int {
	result := blk.id
	if blk.id == idNotSet {
		result = blk.autoId
	}
	return result
}

func (blk *IPAliasesBlock) SetId(id int) {
	blk.id = id
	blk.origHeader = nil
	blk.changed = true
}

// Returns true if real ID value is set,
// otherwise if ID is auto generated, then returns false
func (blk *IPAliasesBlock) IdSet() bool {
	return blk.id <= idNotSet
}

func (blk *IPAliasesBlock) Name() string {
	return blk.name
}

func (blk *IPAliasesBlock) SetName(name string) {
	blk.name = name
	blk.origHeader = nil
	blk.changed = true
}

func (blk *IPAliasesBlock) Note() string {
	return blk.note
}

func (blk *IPAliasesBlock) SetNote(comment string) {
	blk.note = comment
	blk.origHeader = nil
	blk.changed = true
}

func (blk *IPAliasesBlock) Entries() []*IPAliasesEntry {
	return slices.Clone(blk.entries)
}

func (blk *IPAliasesBlock) EntriesByIP(ip string) []*IPAliasesEntry {
	found := make([]*IPAliasesEntry, 0)
	for _, ent := range blk.entries {
		if ent.ip == ip {
			found = append(found, ent)
		}
	}
	return found
}

func (blk *IPAliasesBlock) EntriesByAlias(aliase string) []*IPAliasesEntry {
	found := make([]*IPAliasesEntry, 0)
	for _, ent := range blk.entries {
		if slices.Contains(ent.aliases, aliase) {
			found = append(found, ent)
			break
		}
	}
	return found
}

func newIPAliasesBlockFromElements(headerElements []*syntax.CommentLine, bodyElements []syntax.Element, autoId int) *IPAliasesBlock {
	block := &IPAliasesBlock{
		origHeader: headerElements,
	}
	parseNFillIPsBlockHeader(block, autoId)
	fillIPsBlockBody(block, bodyElements)
	return block
}

func parseNFillIPsBlockHeader(block *IPAliasesBlock, autoId int) {

	var blockId int = idNotSet
	var autoBlockId int = idNotSet
	var blockName string
	var commentsText string

	if len(block.origHeader) > 0 {
		headerLine := block.origHeader[0].CommentText()

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
		for _, line := range block.origHeader[1:] {
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
		autoBlockId = autoId
	}

	block.id = blockId
	block.autoId = autoBlockId
	block.name = blockName
	block.note = commentsText
}

func fillIPsBlockBody(block *IPAliasesBlock, bodyElements []syntax.Element) {
	entries := make([]*IPAliasesEntry, 0)
	if len(bodyElements) > 0 {
		for _, el := range bodyElements {
			item := newIPAliasesEntryFromElement(el)
			entries = append(entries, item)
		}
	}
	block.entries = entries
}
