package syntax

import (
	"io"
)

// Write the content to the document
func format(w io.Writer, doc *Document) error {

	isFirstLine := true

	for _, it := range doc.elements {

		var line string
		var err error

		if it.HasPreformattedText() {
			line = it.PreformattedLineText()
		} else {
			line = it.formatLine()
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
