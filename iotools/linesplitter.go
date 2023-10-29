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
			return 0, nil, bufio.ErrFinalToken
		}
		hadCr = len(token) != advance
		return advance, token, err
	}

	return ensureLastLineSplit
}
