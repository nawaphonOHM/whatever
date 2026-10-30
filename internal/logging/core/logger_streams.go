package core

import (
	"io"
	"os"
)

func isStandardIO(w io.Writer) bool {
	return w == nil || w == os.Stdout || w == os.Stderr
}

func isDiscardOrStdin(w io.Writer) bool {
	return w == os.Stdin || w == io.Discard
}

func isSpecialStream(w io.Writer) bool {
	if isStandardIO(w) {
		return true
	}
	return isDiscardOrStdin(w)
}

func isCloserStream(w io.Writer) io.Closer {
	if isSpecialStream(w) {
		return nil
	}
	if c, ok := w.(io.Closer); ok {
		return c
	}
	return nil
}
