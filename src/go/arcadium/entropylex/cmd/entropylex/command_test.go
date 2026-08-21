package main_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/build"
	el "github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/cmd/entropylex"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/assert"
)

func Test_NewCommand(t *testing.T) {
	tests := []struct {
		name   string
		info   build.Information
		args   []string
		verify func(*testing.T, error, string)
	}{
		// --version
		{
			name: "check version",
			info: build.Info("name", "version", "branch", "commit", "date"),
			args: []string{"cmd", "--version"},
			verify: func(t *testing.T, err error, outfile string) {
				assert.Nil(t, err)
				output, err := os.ReadFile(outfile)
				assert.Nil(t, err)
				assert.Contains(t, string(output), "name version (branch: branch, commit: commit, date: date")
			},
		},
		// --bit-depth
		{
			name: "invalid bit depth",
			info: build.Info("name", "version", "branch", "commit", "date"),
			args: []string{"cmd", "-b"},
			verify: func(t *testing.T, err error, outfile string) {
				assert.Error(t, err, "flag needs an argument: -b")
				output, err := os.ReadFile(outfile)
				assert.Nil(t, err)
				assert.Contains(t, string(output), "Incorrect Usage: flag needs an argument: -b")
			},
		},
		{
			name: "invalid bit depth",
			info: build.Info("name", "version", "branch", "commit", "date"),
			args: []string{"cmd", "-b", "42"},
			verify: func(t *testing.T, err error, outfile string) {
				assert.Error(t, err, "invalid value \"42\" for flag -b: possible values are 8, 12, 14 or 16")
				output, err := os.ReadFile(outfile)
				assert.Nil(t, err)
				assert.Contains(t, string(output), `Incorrect Usage: invalid value "42" for flag -b: possible values are 8, 12, 14 or 16`)
			},
		},
		// --output
		{
			name: "existing file, w/o force",
			info: build.Info("name", "version", "branch", "commit", "date"),
			args: []string{"cmd", "-o", "test/output"},
			verify: func(t *testing.T, err error, outfile string) {
				assert.Error(t, err, "a non-empty output file exist, to overwrite use the --force option")
				output, err := os.ReadFile(outfile)
				assert.Nil(t, err)
				assert.Contains(t, string(output), `Incorrect Usage: a non-empty output file exist, to overwrite use the --force option`)
			},
		},

		// FILE
		{
			name: "multiple filename args",
			info: build.Info("name", "version", "branch", "commit", "date"),
			args: []string{"cmd", "-b", "8", "foo", "bar"},
			verify: func(t *testing.T, err error, outfile string) {
				assert.Error(t, err, "usage error")
				output, err := os.ReadFile(outfile)
				assert.Nil(t, err)
				assert.Contains(t, string(output), `Incorrect Usage: extra input file "bar"`)
			},
		},
		{
			name: "input filename is a directory",
			info: build.Info("name", "version", "branch", "commit", "date"),
			args: []string{"cmd", "dict"},
			verify: func(t *testing.T, err error, outfile string) {
				assert.Error(t, err, `input file "dict" is a directory`)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outfile := ""
			f, err := os.CreateTemp(".", "new_command_test")
			assert.Nil(t, err)
			outfile = f.Name()

			stdout, stderr := os.Stdout, os.Stderr
			os.Stdout = f
			os.Stderr = f

			defer func() {
				os.Remove(outfile)
				os.Stdout, os.Stderr = stdout, stderr
			}()

			err = el.NewCommand(test.info).Run(context.TODO(), test.args)
			test.verify(t, err, outfile)
		})
	}
}

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
