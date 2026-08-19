package main

import (
	"context"
	"errors"
	"io"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/dictionary"
)

const (
	el8LXJFile = "dict/aentropylex-en-8-test-v1.lxj" // FIXME: Change to production version at some point.
)

type (
	EntropyLex8 struct {
		input  io.ReadCloser
		output io.WriteCloser
		decode bool
		dict   *dictionary.LXJ
	}
)

func NewEntropyLex8(input io.ReadCloser, output io.WriteCloser, decode bool) (*EntropyLex8, error) {
	lxjfile, err := dictFS.Open(el8LXJFile)
	if err != nil {
		return nil, err
	}

	lxjf, ok := lxjfile.(io.ReadSeekCloser)
	if !ok {
		return nil, errors.New("internal error: failed to load default lxj dictionary")
	}

	lxj, err := dictionary.NewLXJValidated(lxjf)
	if err != nil {
		return nil, err
	}

	return &EntropyLex8{
		input:  input,
		output: output,
		decode: decode,
		dict:   lxj,
	}, nil
}

func (e *EntropyLex8) Run(ctx context.Context) error {
	// TODO
	return nil
}
