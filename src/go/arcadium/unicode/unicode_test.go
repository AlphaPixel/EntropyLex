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
		verify func(*testing.T, string)
	}{
		{
			name: "test non code point",
			ucp:  unicode.CodePoint("foobar"),
			verify: func(t *testing.T, s string) {
				assert.Equal(t, s, "foobar")
			},
		},
		{
			name: "U+0020",
			ucp:  unicode.CodePoint("U+0020"),
			verify: func(t *testing.T, s string) {
				assert.Equal(t, s, " ")
			},
		},
		{
			name: "U+754C",
			ucp:  unicode.CodePoint("U+754C"),
			verify: func(t *testing.T, s string) {
				assert.Equal(t, s, "界")
			},
		},
		{
			name: "U+1F767",
			ucp:  unicode.CodePoint("U+1F767"),
			verify: func(t *testing.T, s string) {
				assert.Equal(t, s, "🝧")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.verify(t, test.ucp.String())
		})
	}
}
