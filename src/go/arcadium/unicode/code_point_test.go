package unicode_test

import (
	"testing"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/assert"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/unicode"
)

func Test_UnicodeCodePoint_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		ucp    unicode.CodePoint
		verify func(*testing.T, string, error)
	}{
		// errors
		{
			name: "empty string",
			ucp:  unicode.CodePoint(""),
			verify: func(t *testing.T, s string, err error) {
				assert.Equal(t, s, "")
				assert.IsError(t, err, unicode.ErrInvalidCodePoint)
				assert.Error(t, err, `invalid unicode code point, ""`)
			},
		},
		{
			name: "non-canonical padding",
			ucp:  unicode.CodePoint("U+00020"),
			verify: func(t *testing.T, s string, err error) {
				assert.Equal(t, s, "")
				assert.IsError(t, err, unicode.ErrInvalidCodePoint)
				assert.Error(t, err, `invalid unicode code point, U+00020`)
			},
		},
		{
			name: "test non code point",
			ucp:  unicode.CodePoint("foobar"),
			verify: func(t *testing.T, s string, err error) {
				assert.Equal(t, s, "")
				assert.IsError(t, err, unicode.ErrInvalidCodePoint)
				assert.Error(t, err, "invalid unicode code point, foobar")
			},
		},
		{
			name: "surrogate code point error, lower boundry",
			ucp:  unicode.CodePoint("U+D800"),
			verify: func(t *testing.T, s string, err error) {
				assert.Equal(t, s, "")
				assert.IsError(t, err, unicode.ErrInvalidCodePoint)
				assert.Error(t, err, "invalid unicode code point, surrogate code point U+D800")
			},
		},
		{
			name: "surrogate code point error, upper boundry",
			ucp:  unicode.CodePoint("U+DFFF"),
			verify: func(t *testing.T, s string, err error) {
				assert.Equal(t, s, "")
				assert.IsError(t, err, unicode.ErrInvalidCodePoint)
				assert.Error(t, err, "invalid unicode code point, surrogate code point U+DFFF")
			},
		},
		{
			name: "surrogate code point error, mid boundry",
			ucp:  unicode.CodePoint("U+DCDE"),
			verify: func(t *testing.T, s string, err error) {
				assert.Equal(t, s, "")
				assert.IsError(t, err, unicode.ErrInvalidCodePoint)
				assert.Error(t, err, "invalid unicode code point, surrogate code point U+DCDE")
			},
		},
		{
			name: "code point above max error",
			ucp:  unicode.CodePoint("U+110000"),
			verify: func(t *testing.T, s string, err error) {
				assert.Equal(t, s, "")
				assert.IsError(t, err, unicode.ErrInvalidCodePoint)
				assert.Error(t, err, "invalid unicode code point, U+110000")
			},
		},
		// success
		{
			name: "U+0020",
			ucp:  unicode.CodePoint("U+0020"),
			verify: func(t *testing.T, s string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, s, " ")
			},
		},
		{
			name: "U+754C",
			ucp:  unicode.CodePoint("U+754C"),
			verify: func(t *testing.T, s string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, s, "界")
			},
		},
		{
			name: "U+1F767",
			ucp:  unicode.CodePoint("U+1F767"),
			verify: func(t *testing.T, s string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, s, "🝧")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := test.ucp.Decode()
			test.verify(t, s, err)
		})
	}
}
