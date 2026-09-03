// MIT License
//
// Copyright 2026 arcadium.dev <info@arcadium.dev>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"context"
	"fmt"
	"io"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/dictionary"
)

const (
	el8LXJFile = "dict/entropylex-en-8-test-v1.lxj" // FIXME: Change to production version at some point.
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
		return nil, fmt.Errorf("failed to open dictionary \"%s\", %w: %w", el8LXJFile, err, ErrInternal)
	}

	lxjf, ok := lxjfile.(io.ReadSeekCloser)
	if !ok {
		return nil, fmt.Errorf("failed to load default lxj dictionary: %w", ErrInternal)
	}

	lxj, err := dictionary.NewLXJValidated(lxjf)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", err, ErrInternal)
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
