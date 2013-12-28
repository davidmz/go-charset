package charset

import (
	"io"
	"unicode/utf8"
)

type encoder struct {
	w    io.Writer
	cs   *Charset
	_Buf [utf8.UTFMax]byte
	buf  []byte
}

// Encoder returns an io.Writer that takes UTF-8 data,
// converts them to charset and transmits to
// the underlying io.Writer.
//
// UTF-8 data must be aligned to runes boundaries!
func (c *Charset) Encoder(w io.Writer) io.Writer {
	return &encoder{w: w, cs: c}
}

func (e *encoder) Write(p []byte) (int, error) {
	runes := []rune(string(p))
	outBytes := make([]byte, len(runes))

	for i, ru := range runes {
		if b, ok := e.cs.fromUtf[ru]; ok && ru != utf8.RuneError {
			outBytes[i] = b
		} else {
			outBytes[i] = e.cs.ErrorChar
		}
	}
	return e.w.Write(outBytes)
}
