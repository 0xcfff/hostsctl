package syntax

import (
	"io"
)

// Parses the content and returns parsed document
func Read(r io.Reader) (*Document, error) {
	if r == nil {
		return nil, nil
	}

	return parse(r)
}

// Write the content to the document
func Write(w io.Writer, doc *Document, fm FormatMode) error {
	if doc == nil {
		return nil
	}

	return format(w, doc, fm)
}
