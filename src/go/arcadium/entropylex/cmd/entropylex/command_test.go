package main_test

import (
	"context"
	"os"
	"testing"

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
