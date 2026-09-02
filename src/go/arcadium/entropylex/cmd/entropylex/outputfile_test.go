package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"

	el "github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/cmd/entropylex"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/assert"
)

func Test_OutputFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		force    bool
		before   func(*testing.T, string)
		verify   func(*testing.T, *os.File, error)
		after    func(*testing.T, string)
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
			name:     "file does not exist",
			filename: fmt.Sprintf("test/%s", uuid.NewString()),
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, f)
				assert.Nil(t, f.Close())
			},
			after: func(t *testing.T, output string) {
				assert.Nil(t, os.Remove(output))
			},
		},
		{
			name:     "file does not exist, it non-existant path",
			filename: "test/foo/bar/this_should_fail",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, f)
				assert.Error(t, err, "open test/foo/bar/this_should_fail: no such file or directory")
			},
		},
		{
			name:     "file exists, directory",
			filename: "test",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, f)
				assert.Error(t, err, `output file "test" is a directory`)
			},
		},
		{
			name:     "file exists, bad link",
			filename: "test/bad_link",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, f)
				assert.Error(t, err, "open test/bad_link: no such file or directory")
			},
		},
		{
			name:     "file exists, empty",
			filename: "test/output_empty",
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, f)
				fs, e := f.Stat()
				assert.Nil(t, e)
				assert.Equal(t, fs.Size(), 0)
			},
		},
		{
			name:     "file exists, not empty w/o force",
			filename: fmt.Sprintf("test/%s", uuid.NewString()),
			force:    false,
			before: func(t *testing.T, output string) {
				f, err := os.Create(output)
				assert.NotNil(t, f)
				assert.Nil(t, err)
				_, err = f.Write([]byte("testing 1 2 3 4"))
				assert.Nil(t, err)
				assert.Nil(t, f.Close())
			},
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, f)
				assert.Error(t, err, "a non-empty output file exist, to overwrite use the --force option")
			},
			after: func(t *testing.T, output string) {
				assert.Nil(t, os.Remove(output))
			},
		},
		{
			name:     "file exists, not empty w/force",
			filename: fmt.Sprintf("test/%s", uuid.NewString()),
			force:    true,
			before: func(t *testing.T, output string) {
				f, err := os.Create(output)
				assert.NotNil(t, f)
				assert.Nil(t, err)
				_, err = f.Write([]byte("testing 1 2 3 4"))
				assert.Nil(t, err)
				assert.Nil(t, f.Close())
			},
			verify: func(t *testing.T, f *os.File, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, f)
				fs, e := f.Stat()
				assert.Nil(t, e)
				assert.Equal(t, fs.Size(), 0)
				assert.Nil(t, f.Close())
			},
			after: func(t *testing.T, output string) {
				assert.Nil(t, os.Remove(output))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if test.before != nil {
				test.before(t, test.filename)
			}

			f, err := el.OutputFile(test.filename, test.force)
			test.verify(t, f, err)

			if test.after != nil {
				test.after(t, test.filename)
			}
		})
	}
}
