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
	"fmt"
	"io"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/dictionary"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex/entropylex8"
)

const (
	el8LXJFile = "dict/entropylex-en-8-test-v1.lxj" // FIXME: Change to production version at some point.
)

func NewEntropyLex8Encoding() (*entropylex8.Encoding, error) {
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

	return entropylex8.NewEncoding(lxj), nil
}

/*
func (e EntropyLex8) Encoder() io.Writer {
	return e.encoding.Encoder(e.output)
}

func (e EntropyLex8) Decoder() io.Reader {
	// TODO
	return nil

}

func (e *EntropyLex8) Run(ctx context.Context) error {
	if e.decoding {
		return e.decode(ctx)
	}
	return e.encode(ctx)
}

func (e *EntropyLex8) decode(context.Context) error {
	// TODO
	return nil
}

func (e *EntropyLex8) encode(ctx context.Context) error {
	var (
		encoder = contextio.NewWriter(ctx, e.encoding.Encoder(e.output))
		buffer  = make([]byte, 1024)
	)

	for {
		bytesRead, err := os.Stdin.Read(buffer)

		if bytesRead > 0 {
			if _, err := encoder.Write(buffer); err != nil {
				return err
			}
		}

		if err != nil {
			if err == io.EOF {
				return err
			}
			break
		}
	}

	return nil
}
*/
