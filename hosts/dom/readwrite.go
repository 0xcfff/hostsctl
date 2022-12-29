package dom

import (
	"io"

	"github.com/0xcfff/hostsctl/hosts/syntax"
)

func Read(r io.Reader) (*Document, error) {
	sdoc, err := syntax.Read(r)
	if err != nil {
		return nil, err
	}

	return parse(sdoc), nil
}

func Write(w io.Writer, doc *Document) error {
	sdoc := format(doc)
	syntax.Write(w, sdoc)
	return nil
}
