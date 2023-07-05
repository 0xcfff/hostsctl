package syntax

import (
	"io"
	"strings"
)

type formattingContext struct {
	format                  FormatMode
	autoformatSettings      *aliasAutoformattingSettings
	globalAliasFormatParams *aliasFormattingParams
	lastAliasFormatParams   *aliasFormattingParams
}

// Write the content to the document
func format(w io.Writer, doc *Document, fmt FormatMode) error {

	ctx := newFormattingContext(fmt)
	// elements := make([]syntax.Element, 0)

	if fmt == FmtReFormat {
		calculateAutoformats(ctx, doc)
	}

	isFirstLine := true

	for _, it := range doc.elements {

		var line string
		var err error

		switch it.Type() {
		case IPMapping:
			el := it.(*IPMappingLine)
			if fmt == FmtKeep && el.preformattedLineText != nil {
				line = *el.preformattedLineText
			} else {
				fp := ctx.lastAliasFormatParams
				if fp.isEmpty() {
					fp = ctx.globalAliasFormatParams
				}
				line = formatAlias(el, ctx.format, ctx.autoformatSettings, fp)
			}
			cw := calculateAliasesActualColumnWidths(line, el, ctx.autoformatSettings)
			ufp := translateColumnsToAliasesFormattingParams(cw, ctx.autoformatSettings)
			ctx.lastAliasFormatParams.updateFrom(ufp)
		default:
			if it.HasPreformattedText() {
				line = it.PreformattedLineText()
			} else {
				line = it.formatLine()
			}
		}
		if !isFirstLine {
			_, err = io.WriteString(w, "\n")
			if err != nil {
				return err
			}
		}

		_, err = io.WriteString(w, line)
		if err != nil {
			return err
		}

		isFirstLine = false
	}
	return nil
}

func newFormattingContext(fmt FormatMode) *formattingContext {
	ctx := formattingContext{
		format:                  fmt,
		autoformatSettings:      defaultAliasFormattingSettings,
		globalAliasFormatParams: newAliasFormattingParams(),
		lastAliasFormatParams:   newAliasFormattingParams(),
	}
	return &ctx
}

func calculateAutoformats(ctx *formattingContext, doc *Document) {
	settings := defaultAliasFormattingSettings
	widths := newAliasColumnWidths()
	for _, el := range doc.elements {
		switch el.Type() {
		case IPMapping:
			rwidths := calculateAliasesTheoreticalColumnWidths(el.(*IPMappingLine), settings)
			widths.updateFrom(rwidths)
		}
	}
	ctx.globalAliasFormatParams = calculateAliasesFormattingParams(widths, settings)
}

func calculateAliasesFormattingParams(cw *aliasColumnsWidths, fs *aliasAutoformattingSettings) *aliasFormattingParams {
	fmt := newAliasFormattingParams()
	fmt.ipPosition = fs.minSpacingToIP
	fmt.aliasPosition = fmt.ipPosition + cw.ip + fs.minSpacingToAlias
	fmt.commentPosition = fmt.aliasPosition + cw.alias + fs.minSpacingToComment
	return fmt
}

func translateColumnsToAliasesFormattingParams(cw *aliasColumnsWidths, fs *aliasAutoformattingSettings) *aliasFormattingParams {
	fmt := newAliasFormattingParams()
	fmt.ipPosition = fs.minSpacingToIP
	fmt.aliasPosition = cw.ip
	fmt.commentPosition = cw.ip + cw.alias
	return fmt
}

func calculateAliasesActualColumnWidths(line string, el *IPMappingLine, settings *aliasAutoformattingSettings) *aliasColumnsWidths {

	cols := newAliasColumnWidths()

	aliasLeft := len(line)
	for _, a := range el.domainNames {
		idx := strings.Index(line, a)
		if idx < aliasLeft {
			aliasLeft = idx
		}
	}
	cols.ip = aliasLeft

	if el.commentText != "" {
		cols.alias = strings.Index(line, el.commentText) - cols.ip
		cols.comment = len(el.commentText)
	} else {
		cols.alias = len(line) - cols.ip
	}

	return cols
}

func calculateAliasesTheoreticalColumnWidths(el *IPMappingLine, settings *aliasAutoformattingSettings) *aliasColumnsWidths {

	cols := newAliasColumnWidths()
	cols.ip = len(el.ip)
	cols.comment = len(el.commentText)

	first := true
	for _, a := range el.domainNames {
		if first {
			first = false
		} else {
			cols.comment += settings.minSpacingBetweenAliases
		}
		cols.comment += len(a)
	}
	return cols
}

func formatAlias(el *IPMappingLine, fm FormatMode, fs *aliasAutoformattingSettings, fp *aliasFormattingParams) string {

	var ipAlias string

	if fm == FmtKeep {
		if el.preformattedLineText != nil {
			ipAlias = *el.preformattedLineText
		}
	}

	if ipAlias == "" {
		b := strings.Builder{}
		if fp.ipPosition > 0 {
			b.WriteString(strings.Repeat(" ", fp.ipPosition))
		}
		b.WriteString(el.ip)
		if fp.aliasPosition > b.Len() {
			b.WriteString(strings.Repeat(" ", fp.aliasPosition-b.Len()))
		}
		if !strings.HasSuffix(b.String(), " ") {
			b.WriteString(" ")
		}
		first := true
		for _, a := range el.domainNames {
			if first {
				first = false
			} else {
				b.WriteString(strings.Repeat(" ", fs.minSpacingBetweenAliases))
			}
			b.WriteString(a)
		}
		if el.commentText != "" {
			if fp.commentPosition > b.Len() {
				b.WriteString(strings.Repeat(" ", fp.commentPosition-b.Len()))
			}
			orig := b.String()
			trimmed := strings.TrimRight(orig, " \t")
			spacing := len(orig) - len(trimmed)
			if fp.commentPosition != b.Len() && spacing < fs.minSpacingToAlias {
				b.WriteString(strings.Repeat(" ", fs.minSpacingToAlias-spacing))
			}

			if !strings.HasSuffix(b.String(), " ") {
				b.WriteString(" ")
			}
			b.WriteString("# ")
			b.WriteString(el.commentText)
		}
		ipAlias = b.String()
	}

	return ipAlias
}
