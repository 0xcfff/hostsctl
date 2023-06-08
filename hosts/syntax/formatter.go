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
				fp := ctx.globalAliasFormatParams
				if fp == nil {
					fp = ctx.lastAliasFormatParams
				}
				line = formatAlias(el, ctx.format, ctx.autoformatSettings, fp)
				cw := calculateAliasesColumnWidths(el, ctx.autoformatSettings)
				ufp := calculateAliasesFormattingParams(cw, ctx.autoformatSettings)
				ctx.lastAliasFormatParams.updateFrom(ufp)
			}
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
			rwidths := calculateAliasesColumnWidths(el.(*IPMappingLine), settings)
			widths.updateFrom(rwidths)
			break
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

func calculateAliasesColumnWidths(el *IPMappingLine, settings *aliasAutoformattingSettings) *aliasColumnsWidths {

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
			b.WriteString(strings.Repeat(" ", b.Len()-fp.aliasPosition))
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
		if fp.commentPosition > b.Len() {
			b.WriteString(strings.Repeat(" ", b.Len()-fp.commentPosition))
		}
		if !strings.HasSuffix(b.String(), " ") {
			b.WriteString(" ")
		}
		b.WriteString(el.commentText)

		ipAlias = b.String()
	}

	return ipAlias
}
