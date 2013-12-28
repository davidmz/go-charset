// The charset package provides encoding/decoding functions for single-byte character sets.
package charset

// Charset is an object represented a single-byte character set.
type Charset struct {
	// Character ('?' by default) for replacing invalid or
	// undefined in this charset UTF-8 runes.
	ErrorChar byte

	toUtf   [256]rune
	fromUtf map[rune]byte
}

// New function created a new Charset from given symbol table.
func New(table [256]rune) *Charset {
	cs := &Charset{
		toUtf:     table,
		fromUtf:   make(map[rune]byte, 256),
		ErrorChar: '?',
	}
	for i, r := range cs.toUtf {
		cs.fromUtf[r] = byte(i)
	}
	return cs
}

// Encode encodes UFT-8 string to bytes according with charset
func (c *Charset) Encode(s string) []byte {
	a := []rune(s)
	out := make([]byte, len(a))
	for i, r := range a {
		if b, ok := c.fromUtf[r]; ok {
			out[i] = b
		} else {
			out[i] = c.ErrorChar
		}
	}
	return out
}

// Decode decodes bytes to UFT-8 string according with charset
func (c *Charset) Decode(bytes []byte) string {
	out := make([]rune, len(bytes))
	for i, b := range bytes {
		out[i] = c.toUtf[b]
	}
	return string(out)
}
