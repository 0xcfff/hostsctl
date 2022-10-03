package printutil

import (
	"bytes"
	"io"

	"github.com/liggitt/tabwriter"
)

func CalulateWidthsTabbed(padding int, f func(w io.Writer) error) ([]int, error) {
	tw := tabwriter.NewWriter(io.Discard, 0, 4, padding, ' ', 0)
	err := f(tw)

	if err != nil {
		return nil, err
	}
	tw.Flush()
	res := tw.RememberedWidths()
	return res, nil

}

func FormatTabbed(widths []int, padding int, f func(w io.Writer) error) (string, error) {
	buff := &bytes.Buffer{}
	tw := tabwriter.NewWriter(buff, 0, 4, padding, ' ', 0)
	if widths != nil {
		tw.SetRememberedWidths(widths)
	}
	err := f(tw)

	if err != nil {
		return "", err
	}
	tw.Flush()
	res := buff.String()
	return res, nil
}

func PrintTabbed(w io.Writer, widths []int, padding int, f func(w io.Writer) error) error {
	tw := tabwriter.NewWriter(w, 0, 4, padding, ' ', 0)
	if widths != nil {
		tw.SetRememberedWidths(widths)
	}
	err := f(tw)

	if err != nil {
		return err
	}
	tw.Flush()
	return nil
}
