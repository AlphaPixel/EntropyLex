package dictionary_test

import (
	"io"
	"os"
	"testing"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/dictionary"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/assert"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/require"
)

func Test_NewLXJ(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reader func(*testing.T) io.ReadSeeker
		verify func(*testing.T, *dictionary.LXJ, error)
	}{
		{
			name:   "nil reader",
			reader: func(*testing.T) io.ReadSeeker { return nil },
			verify: func(t *testing.T, dict *dictionary.LXJ, err error) {
				assert.Nil(t, dict)
				assert.Error(t, err, "invalid dictionary file")
			},
		},
		{
			name: "valid lxj dictionary",
			reader: func(t *testing.T) io.ReadSeeker {
				f, err := os.Open("el-8-valid.lxj")
				assert.Nil(t, err)
				return f
			},
			verify: func(t *testing.T, dict *dictionary.LXJ, err error) {
				assert.Nil(t, err)
				assert.Equal(t, dict.Format.Name, "LXJ")
				require.Equal(t, len(dict.Tokens), 256)
				assert.Equal(t, dict.Tokens[0], "able")
				assert.Equal(t, dict.Tokens[255], "power")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dict, err := dictionary.NewLXJ(test.reader(t))
			test.verify(t, dict, err)
		})
	}
}

func Test_NewLXJValidated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reader func(*testing.T) io.ReadSeeker
		verify func(*testing.T, *dictionary.LXJ, error)
	}{
		{
			name:   "nil reader",
			reader: func(*testing.T) io.ReadSeeker { return nil },
			verify: func(t *testing.T, dict *dictionary.LXJ, err error) {
				assert.Nil(t, dict)
				assert.Error(t, err, "invalid dictionary file")
			},
		},
		{
			name: "invalid dictionary",
			reader: func(t *testing.T) io.ReadSeeker {
				f, err := os.Open("el-8-invalid.lxj")
				assert.Nil(t, err)
				return f
			},
			verify: func(t *testing.T, dict *dictionary.LXJ, err error) {
				assert.Nil(t, dict)
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), "jsonschema validation failed with")
				assert.Contains(t, err.Error(), "at '/format/name': value must be 'LXJ")
			},
		},
		{
			name: "valid lxj dictionary",
			reader: func(t *testing.T) io.ReadSeeker {
				f, err := os.Open("el-8-valid.lxj")
				assert.Nil(t, err)
				return f
			},
			verify: func(t *testing.T, dict *dictionary.LXJ, err error) {
				assert.Nil(t, err)
				assert.Equal(t, dict.Format.Name, "LXJ")
				require.Equal(t, len(dict.Tokens), 256)
				assert.Equal(t, dict.Tokens[0], "able")
				assert.Equal(t, dict.Tokens[255], "power")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			d, err := dictionary.NewLXJValidated(test.reader(t))
			test.verify(t, d, err)
		})
	}
}
