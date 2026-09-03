package main_test

import (
	"testing"

	el "github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/cmd/entropylex"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/test/assert"
)

func Test_ErrorCode(t *testing.T) {

	tests := []struct {
		name   string
		err    el.ErrorCode
		verify func(t *testing.T, err el.ErrorCode)
	}{
		{
			name: "ErrUsage",
			err:  el.ErrUsage,
			verify: func(t *testing.T, err el.ErrorCode) {
				assert.Error(t, err, "usage error")
				assert.Equal(t, err.Code(), el.UsageErrorCode)
			},
		},
		{
			name: "ErrInternal",
			err:  el.ErrInternal,
			verify: func(t *testing.T, err el.ErrorCode) {
				assert.Error(t, err, "internal error")
				assert.Equal(t, err.Code(), el.InternalErrorCode)
			},
		},
		{
			name: "ErrUnimplemented",
			err:  el.ErrUnimplemented,
			verify: func(t *testing.T, err el.ErrorCode) {
				assert.Error(t, err, "unimplemented")
				assert.Equal(t, err.Code(), el.UnimplementedErrorCode)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.verify(t, test.err)
		})
	}
}
