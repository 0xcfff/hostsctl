package dom

import (
	"strconv"
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"github.com/0xcfff/hostsctl/iotools"
	"golang.org/x/exp/slices"
)

const idNotSet = -1
const (
	Alias IPAliasesBlockElementType = iota
	Placeholder
)

type IPAliasesBlockElementType int

type IPAliasesBlockElement interface {
	Type() IPAliasesBlockElementType
	ClearFormatting()
	dirty() bool
}

// Block of IPs
type IPAliasesBlock struct {
	origHeader []*syntax.CommentLine
	id         int
	autoId     int
	name       string
	note       string
	entries    []IPAliasesBlockElement
	changed    bool
}

func (blk *IPAliasesBlock) Type() BlockType {
	return IPList
}

func (blk *IPAliasesBlock) dirty() bool {
	if blk.changed {
		return true
	}

	for _, ent := range blk.entries {
		if ent.dirty() {
			return true
		}
	}

	return false
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
	return blk.id >= idNotSet
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

func (blk *IPAliasesBlock) ClearFormatting() {
	blk.origHeader = nil
	blk.changed = true

	for _, el := range blk.entries {
		el.ClearFormatting()
	}
}

func (blk *IPAliasesBlock) Entries() []IPAliasesBlockElement {
	return slices.Clone(blk.entries)
}

func (blk *IPAliasesBlock) AliasEntries() []*IPAliasesEntry {
	aliasEntries := filterSliceByTypeAndPredicate(blk.entries, func(ent *IPAliasesEntry) bool { return true })
	return aliasEntries
}

func (blk *IPAliasesBlock) AliasEntriesByIP(ip string) []*IPAliasesEntry {
	found := filterSliceByTypeAndPredicate(blk.entries, func(ent *IPAliasesEntry) bool { return ent.ip == ip })
	return found
}

func (blk *IPAliasesBlock) AliasEntriesByAlias(alias string) []*IPAliasesEntry {
	found := filterSliceByTypeAndPredicate(blk.entries, func(ent *IPAliasesEntry) bool { return slices.Contains(ent.aliases, alias) })
	return found
}

func (blk *IPAliasesBlock) AliasEntriesByIPOrAlias(ipOrAlias string) []*IPAliasesEntry {
	entries := blk.AliasEntriesByIP(ipOrAlias)
	if len(entries) == 0 {
		entries = blk.AliasEntriesByAlias(ipOrAlias)
	}
	return entries
}

func (blk *IPAliasesBlock) AddEntry(entry IPAliasesBlockElement) {
	blk.entries = append(blk.entries, entry)
}

func (blk *IPAliasesBlock) RemoveEntry(entry IPAliasesBlockElement) bool {
	condition := func(it IPAliasesBlockElement) bool { return it == entry }
	newEntries, changed := removeElements(blk.entries, condition)
	if changed {
		blk.entries = newEntries
		blk.changed = true
	}
	return changed
}

func (blk *IPAliasesBlock) normalize() bool {
	normalized := false
	var placeholder IPAliasesBlockElement
	hasIPs := false
	for _, ent := range blk.entries {
		switch ent.Type() {
		case Placeholder:
			placeholder = ent
		case Alias:
			hasIPs = true
		default:
		}
	}

	if hasIPs && placeholder != nil {
		blk.RemoveEntry(placeholder)
		normalized = true
	} else if !hasIPs && placeholder == nil {
		blk.AddEntry(NewIPAliasesPlaceholder())
		normalized = true
	}

	return normalized
}

func NewIPAliasesBlock() *IPAliasesBlock {
	return &IPAliasesBlock{
		id:      idNotSet,
		autoId:  idNotSet,
		changed: true,
	}
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
		newLine := iotools.NewLine
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
	entries := make([]IPAliasesBlockElement, 0)
	if len(bodyElements) > 0 {
		for _, el := range bodyElements {
			item := newIPAliasesEntryFromElement(el)
			entries = append(entries, item)
		}
	}
	block.entries = entries
}

func filterSliceByTypeAndPredicate[B any, S any](items []S, match func(block B) bool) []B {
	result := make([]B, 0)
	for _, blk := range items {
		if tblk, ok := any(blk).(B); ok {
			if match(tblk) {
				result = append(result, tblk)
			}
		}
	}
	return result
}
