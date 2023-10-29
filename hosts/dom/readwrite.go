package dom

import (
	"io"

	"github.com/0xcfff/hostsctl/hosts/syntax"
)

// Documet formatting mode
type FmtMode int

const (
	FmtKeep     FmtMode = iota    // keep original formatting where possible
	FmtReFormat FmtMode = iota    // re-format document
	FmtDefault          = FmtKeep // default pormatting (Keep)
)

// Read and parse document from a reader
func Read(r io.Reader) (*Document, error) {
	sdoc, err := syntax.Read(r)
	if err != nil {
		return nil, err
	}

	return parse(sdoc), nil
}

// Write document to a writer with specified formatting
func Write(w io.Writer, doc *Document, fm FmtMode) error {
	sdoc := constructSyntax(doc)
	sfm := fm.toSyntaxFormat()
	// TODO: Add re-formatting logic at dom mode (remove unneeded spaces, etc)
	syntax.Write(w, sdoc, sfm)
	return nil
}

func (fm FmtMode) toSyntaxFormat() syntax.FormatMode {
	switch fm {
	case FmtKeep:
		return syntax.FmtKeep
	case FmtReFormat:
		return syntax.FmtReFormat
	}
	return syntax.FmtDefault
}
