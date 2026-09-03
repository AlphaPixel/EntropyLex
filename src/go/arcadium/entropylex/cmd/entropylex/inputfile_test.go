package main_test

import (
	"os"
	"testing"

	el "github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/cmd/entropylex"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/assert"
)

func Test_InputFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		verify   func(*testing.T, *os.File, error)
	}{
		{
			name:     "empty filename",
			filename: "",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, err)
				assert.Nil(t, f)
			},
		},
		{
			name:     "- filename",
			filename: "-",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, err)
				assert.Nil(t, f)
			},
		},
		{
			name:     "unknown filename",
			filename: "xyq.pdq",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, f)
				assert.Error(t, err, "stat xyq.pdq: no such file or directory")
			},
		},
		{
			name:     "directory",
			filename: "test",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, f)
				assert.Error(t, err, `input file "test" is a directory`)
			},
		},
		{
			name:     "link to directory",
			filename: "./test/dir_link",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, f)
				assert.Error(t, err, `input file "./test/dir_link" is a directory`)
			},
		},
		{
			name:     "bad link",
			filename: "./test/bad_link",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, f)
				assert.Error(t, err, "stat ./test/bad_link: no such file or directory")
			},
		},
		{
			name:     "good file",
			filename: "main.go",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.NotNil(t, f)
				assert.Nil(t, err)
				assert.Nil(t, f.Close())
			},
		},
		{
			name:     "link to good file",
			filename: "test/good_link",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.NotNil(t, f)
				assert.Nil(t, err)
				assert.Nil(t, f.Close())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			f, err := el.InputFile(test.filename)
			test.verify(t, f, err)
		})
	}
}
