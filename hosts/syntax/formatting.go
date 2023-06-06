package syntax

type FormatMode int

const (
	FmtKeep     FormatMode = iota    // Keep original formatting where possible
	FmtReFormat FormatMode = iota    // Re-format document
	FmtDefault             = FmtKeep // Same as FmtKeep
)

var (
	// Default alias autoformatting settings
	defaultAliasFormattingSettings *aliasAutoformattingSettings = &aliasAutoformattingSettings{
		minSpacingToIP:           0,
		minSpacingToAlias:        8,
		minSpacingBetweenAliases: 2,
		minSpacingToComment:      1,
	}
)

// Settings for automatic document formatting
type aliasAutoformattingSettings struct {
	minSpacingToIP           int // minimal space before IP
	minSpacingToAlias        int // minimal space between IP and its aliases
	minSpacingBetweenAliases int // minimal space between aliases to the same IP
	minSpacingToComment      int // minimal space between IP aliases and comment line
}

// IP aliases block formatting parameters
type aliasFormattingParams struct {
	ipPosition      int // position where IP part is starting
	aliasPosition   int // position where alias part is starting
	commentPosition int // position where comments line is starting
}

// Widths of alias columns
type aliasColumnsWidths struct {
	ip      int
	alias   int
	comment int
}

func newAliasFormattingParams() *aliasFormattingParams {
	return &aliasFormattingParams{
		ipPosition:      -1,
		aliasPosition:   -1,
		commentPosition: -1,
	}
}

// Updates formatting parameters from the specified parameters
func (fmt *aliasFormattingParams) updateFrom(other *aliasFormattingParams) {
	ipDiff := fmt.ipPosition - other.ipPosition
	if ipDiff < 0 {
		fmt.ipPosition -= ipDiff
	}
	aliasDiff := fmt.aliasPosition - other.aliasPosition
	if aliasDiff < 0 {
		fmt.aliasPosition -= aliasDiff
	}
	commentsDiff := fmt.commentPosition - other.commentPosition
	if commentsDiff < 0 {
		fmt.commentPosition -= commentsDiff
	}
}

// Makes a copy of the params object
func (fmt *aliasFormattingParams) clone() *aliasFormattingParams {
	return &aliasFormattingParams{
		ipPosition:      fmt.ipPosition,
		aliasPosition:   fmt.aliasPosition,
		commentPosition: fmt.commentPosition,
	}
}

func newAliasColumnWidths() *aliasColumnsWidths {
	return &aliasColumnsWidths{
		ip:      0,
		alias:   0,
		comment: 0,
	}
}

func (widths *aliasColumnsWidths) updateFrom(other *aliasColumnsWidths) {
	if widths.ip < other.ip {
		widths.ip = other.ip
	}
	if widths.alias < other.alias {
		widths.alias = other.alias
	}
	if widths.comment < other.comment {
		widths.comment = other.comment
	}
}
