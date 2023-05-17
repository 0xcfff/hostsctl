package dom

import (
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"github.com/0xcfff/hostsctl/iptools"
	"golang.org/x/exp/slices"
)

type IPAliasesLine struct {
	origElement syntax.Element
	ip          string
	aliases     []string
	note        string
	disabled    bool
	changed     bool
}

func (blk *IPAliasesLine) IP() string {
	return blk.ip
}

func (blk *IPAliasesLine) SetIP(ip string) {
	if strings.Compare(ip, blk.ip) != 0 {
		blk.ip = ip
		blk.origElement = nil
		blk.changed = true
	}
}

func (blk *IPAliasesLine) Aliases() []string {
	return slices.Clone(blk.aliases)
}

func (blk *IPAliasesLine) AddAlias(alias string) bool {
	shouldAdd := blk.aliases == nil || !slices.Contains(blk.aliases, alias)
	if shouldAdd {
		blk.aliases = append(blk.aliases, alias)
		blk.origElement = nil
		blk.changed = true
	}
	return shouldAdd
}

func (blk *IPAliasesLine) RemoveAlias(alias string) bool {
	condition := func(it string) bool { return it == alias }
	newAliases, changed := removeElements(blk.aliases, condition)
	if changed {
		blk.aliases = newAliases
		blk.origElement = nil
		blk.changed = true
	}
	return changed
}

func (blk *IPAliasesLine) Note() string {
	return blk.note
}

func (blk *IPAliasesLine) SetNote(comment string) {
	if strings.Compare(comment, blk.note) != 0 {
		blk.note = comment
		blk.origElement = nil
		blk.changed = true
	}
}

func (blk *IPAliasesLine) Disabled() bool {
	return blk.disabled
}

func (blk *IPAliasesLine) SetDisabled(disabled bool) {
	if blk.disabled != disabled {
		blk.disabled = disabled
		blk.origElement = nil
		blk.changed = true
	}
}

func removeElements[T any](l []T, remove func(T) bool) ([]T, bool) {
	out := make([]T, 0)
	changed := false
	for _, element := range l {
		if remove(element) {
			changed = true
		} else {
			out = append(out, element)
		}
	}
	return out, changed
}

func newIPMappingFromElement(element syntax.Element) *IPAliasesLine {
	switch element.Type() {
	case syntax.IPMapping:
		ip := element.(*syntax.IPMappingLine)
		return newIPMappingFromIPElement(ip)
	case syntax.Comment:
		comment := element.(*syntax.CommentLine)
		if !isCommentedIPMapping(comment) {
			panic("Specified comment line is not IP")
		}
		return newIPMappingFromCommentElement(comment)
	default:
		panic("Not supported element type")
	}
}

func newIPMappingFromIPElement(ip *syntax.IPMappingLine) *IPAliasesLine {
	item := &IPAliasesLine{
		origElement: ip,
		ip:          ip.IPAddress(),
		aliases:     slices.Clone(ip.DomainNames()),
		note:        ip.CommentText(),
	}
	return item
}

func newIPMappingFromCommentElement(comment *syntax.CommentLine) *IPAliasesLine {
	if !isCommentedIPMapping(comment) {
		panic("Specified comment is not an IP")
	}

	origCommentText := comment.CommentText()
	parts := strings.Fields(origCommentText)
	aliases := make([]string, 0)
	commentText := ""

	for _, it := range parts[1:] {
		if strings.HasPrefix(it, "#") {
			idx := strings.Index(origCommentText, "#")
			commentText = strings.TrimSpace(origCommentText[idx+1:])
			break
		}
		aliases = append(aliases, it)
	}

	item := &IPAliasesLine{
		origElement: comment,
		ip:          parts[0],
		aliases:     aliases,
		note:        commentText,
		disabled:    true,
	}
	return item
}

func isCommentedIPMapping(comment *syntax.CommentLine) bool {
	parts := strings.Fields(comment.CommentText())
	return len(parts) >= 2 && iptools.IsIP(parts[0])
}
