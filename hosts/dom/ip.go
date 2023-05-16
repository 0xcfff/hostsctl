package dom

import (
	"strings"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"golang.org/x/exp/slices"
)

type IPItem struct {
	element     syntax.Element
	ip          string
	domainNames []string
	comment     string
	disabled    bool
	changed     bool
}

func (blk *IPItem) IP() string {
	return blk.ip
}

func (blk *IPItem) SetIP(ip string) {
	if strings.Compare(ip, blk.ip) != 0 {
		blk.ip = ip
		blk.element = nil
		blk.changed = true
	}
}

func (blk *IPItem) Aliases() []string {
	return slices.Clone(blk.domainNames)
}

func (blk *IPItem) AddAlias(alias string) bool {
	shouldAdd := blk.domainNames == nil || !slices.Contains(blk.domainNames, alias)
	if shouldAdd {
		blk.domainNames = append(blk.domainNames, alias)
		blk.element = nil
		blk.changed = true
	}
	return shouldAdd
}

func (blk *IPItem) RemoveAlias(alias string) bool {
	condition := func(it string) bool { return it == alias }
	newAliases, changed := removeElements(blk.domainNames, condition)
	if changed {
		blk.domainNames = newAliases
		blk.element = nil
		blk.changed = true
	}
	return changed
}

func (blk *IPItem) Comment() string {
	return blk.comment
}

func (blk *IPItem) SetComment(comment string) {
	if strings.Compare(comment, blk.comment) != 0 {
		blk.comment = comment
		blk.element = nil
		blk.changed = true
	}
}

func (blk *IPItem) Disabled() bool {
	return blk.disabled
}

func (blk *IPItem) SetDisabled(disabled bool) {
	if blk.disabled != disabled {
		blk.disabled = disabled
		blk.element = nil
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
