package charset

import (
	"io"
	"unicode/utf8"
)

// Decoder returns an io.Reader that reads charset-encoded data from
// the underlying io.Reader and converts them  to UTF-8.
func (c *Charset) Decoder(r io.Reader) io.Reader {
	return &decoder{r: r, cs: c}
}

type decoder struct {
	r    io.Reader
	cs   *Charset
	_Buf [utf8.UTFMax]byte
	buf  []byte // слайс в _Buf
}

// читаем байты, возвращаем байты UTF-8
func (d *decoder) Read(p []byte) (int, error) {
	bytesReaded := 0
	for {
		// Если в buf лежит прочитанная руна
		if len(d.buf) > 0 {
			copied := copy(p, d.buf)
			d.buf, p = d.buf[copied:], p[copied:]
			bytesReaded += copied
		}
		if len(p) == 0 {
			// больше писать некуда
			return bytesReaded, nil
		} else {
			if err := d.readNextRune(); err != nil {
				return bytesReaded, err
			}
		}
	}
}

// readNextRune читает следующий байт из потока, преобразует его
// в UTF-8-последовательность и сохраняет в _Buf и buf
func (d *decoder) readNextRune() error {
	var b [1]byte
	if _, err := d.r.Read(b[:]); err != nil {
		return err
	}
	ru := d.cs.toUtf[b[0]]
	n := utf8.EncodeRune(d._Buf[:], ru)
	d.buf = d._Buf[0:n]
	return nil
}
