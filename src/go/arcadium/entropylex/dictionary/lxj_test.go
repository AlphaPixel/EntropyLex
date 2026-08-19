package dictionary_test

import (
	"testing"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/dictionary"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/assert"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/require"
)

func Test_NewLXJ(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		verify   func(*testing.T, *dictionary.Dictionary, error)
	}{
		{
			name:     "empty filename",
			filename: "",
			verify: func(t *testing.T, d *dictionary.Dictionary, err error) {
				assert.Nil(t, d)
				assert.Error(t, err, "file name required")
			},
		},
		{
			name:     "file doesn't exist",
			filename: "file-does-not-exist.json",
			verify: func(t *testing.T, d *dictionary.Dictionary, err error) {
				assert.Nil(t, d)
				assert.Error(t, err, "open file-does-not-exist.json: no such file or directory")
			},
		},
		{
			name:     "valid dictionary",
			filename: "el-8-valid.lxj",
			verify: func(t *testing.T, d *dictionary.Dictionary, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			d, err := dictionary.NewLXJ(test.filename)
			test.verify(t, d, err)
		})
	}
}

func Test_NewXJValidated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		verify   func(*testing.T, *dictionary.Dictionary, error)
	}{
		{
			name:     "empty filename",
			filename: "",
			verify: func(t *testing.T, d *dictionary.Dictionary, err error) {
				assert.Nil(t, d)
				assert.Error(t, err, "file name required")
			},
		},
		{
			name:     "file doesn't exist",
			filename: "file-does-not-exist.json",
			verify: func(t *testing.T, d *dictionary.Dictionary, err error) {
				assert.Nil(t, d)
				assert.Error(t, err, "open file-does-not-exist.json: no such file or directory")
			},
		},
		{
			name:     "invalid dictionary",
			filename: "el-8-invalid.lxj",
			verify: func(t *testing.T, d *dictionary.Dictionary, err error) {
				assert.Nil(t, d)
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), "jsonschema validation failed with")
				assert.Contains(t, err.Error(), "at '/format/name': value must be 'LXJ")
			},
		},
		{
			name:     "valid dictionary",
			filename: "el-8-valid.lxj",
			verify: func(t *testing.T, d *dictionary.Dictionary, err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			d, err := dictionary.NewLXJValidated(test.filename)
			test.verify(t, d, err)
		})
	}
}
