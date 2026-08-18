package iotools

import "bufio"

// Creates a split function which takes into account last newline
// and returns an emply line if a file ends with new line char sequence
func LinesSplitterRespectEndNewLineFunc() bufio.SplitFunc {

	hadCr := false

	ensureLastLineSplit := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		if atEOF && advance == 0 && token == nil && hadCr {
			hadCr = false
			// Return a non-nil empty token: Go 1.22 changed bufio.Scanner to
			// stop emitting the final token when ErrFinalToken carries a nil
			// token. An empty (non-nil) slice preserves the trailing-newline
			// line while keeping the same behavior on older Go versions.
			return 0, []byte{}, bufio.ErrFinalToken
		}
		hadCr = len(token) != advance
		return advance, token, err
	}

	return ensureLastLineSplit
}
