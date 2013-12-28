package charset

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

var testData = []struct {
	utf8,
	encoded string
	charset *Charset
}{
	//{"", "", CP1251},
	{"Привет", "\xcf\xf0\xe8\xe2\xe5\xf2", CP1251},
	{"hello!", "hello!", CP1251},
}

func TestDecoding(t *testing.T) {
	for i, data := range testData {
		charset := data.charset

		if u := charset.Decode([]byte(data.encoded)); u != data.utf8 {
			t.Errorf("#%d: 'Decode' error: %v (expected %v)", i, u, data.utf8)
		}

		dec := charset.Decoder(bytes.NewBufferString(data.encoded))
		b := make([]byte, 1)
		if _, err := dec.Read(b[:]); err != nil {
			t.Errorf("#%d:'Decoder' error while reading byte: %v", i, err)
		} else if b[0] != []byte(data.utf8)[0] {
			t.Errorf("#%d:'Decoder' unexpected result: %#v (expected %#v)", i, b[0], []byte(data.utf8)[0])
		}
		if _, err := dec.Read(b[:]); err != nil {
			t.Errorf("#%d:'Decoder' error while reading byte: %v", i, err)
		} else if b[0] != []byte(data.utf8)[1] {
			t.Errorf("#%d:'Decoder' unexpected result: %#v (expected %#v)", i, b[0], []byte(data.utf8)[0])
		}
		if _, err := dec.Read(b[:]); err != nil {
			t.Errorf("#%d:'Decoder' error while reading byte: %v", i, err)
		} else if b[0] != []byte(data.utf8)[2] {
			t.Errorf("#%d:'Decoder' unexpected result: %#v (expected %#v)", i, b[0], []byte(data.utf8)[0])
		}

		dec = charset.Decoder(bytes.NewBufferString(data.encoded))
		out, err := ioutil.ReadAll(dec)
		if err != nil {
			t.Errorf("#%d:'Decoder' error while reading: %v", i, err)
		}
		if string(out) != data.utf8 {
			t.Errorf("#%d:'Decoder' unexpected result: %#v (expected %#v)", i, out, []byte(data.utf8))
		}
	}
}

func TestEncoding(t *testing.T) {
	for i, data := range testData {
		charset := data.charset

		if e := string(charset.Encode(data.utf8)); e != data.encoded {
			t.Errorf("#%d: 'Encode' error: %#v (expected %#v)", i, e, data.encoded)
		}

		buf := new(bytes.Buffer)
		encoder := charset.Encoder(buf)

		io.Copy(encoder, bytes.NewBufferString(data.utf8))
		if s := buf.String(); s != data.encoded {
			t.Errorf("#%d:'Decoder' unexpected result: '% x' (expected '% x')", i, s, data.encoded)
		}
	}
}
