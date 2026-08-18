package entropylex8

import (
	"io"
)

type (
	Encoding interface {
		Decode(dst, src []byte) (n int, err error)
		Encode(dst, src []byte)
	}
)

// NewEncoder returns a new entropy lex stream encoder. Data written to the
// returned writer will be encoded using enc and then written to w. When
// finished writing, the caller must Close the returned encoder.
func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser {
	return nil
}

// NewDecoder constructs a new entropy lex stream decoder.
func NewDecoder(enc *Encoding, r io.Reader) io.Reader {
	return nil
}
