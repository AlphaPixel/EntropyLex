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
	"errors"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/dolmen-go/contextio"
	"github.com/urfave/cli/v3"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/build"
	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/entropylex"
)

func NewCommand(info build.Information) *cli.Command {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Println(info.String())
	}

	name := filepath.Base(os.Args[0])

	description := `The data are encoded and decoded as described in the EntropyLex
specification[1].

With no FILE, or when FILE is -, standard input is read.

If an output file option (--output/-o) is not specified, the command output
defaults to stdout. If an output file option is specfied, the file may either 
be an empty existing file, or not exist. When creating an output file, if the
file name is part of a directory path, the directory must exist. If an
non-empty output file exists, it can be overwritten using the -force option.

[1] https://github.com/AlphaPixel/EntropyLex/blob/main/SPEC.md
`

	cmd := &cli.Command{
		Name:        name,
		Usage:       "encode/decode binary data to/from a lexical encoding",
		UsageText:   fmt.Sprintf("%s [OPTION]... [FILE]", name),
		Description: description,
		Version:     info.Version,
		Authors: []any{
			&mail.Address{Name: "Ian Cahoon", Address: "ian@arcadium.dev"},
		},

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "decode",
				Aliases: []string{"d"},
				Usage:   "decode data",
			},
			&cli.UintFlag{
				Name:    "bit-depth",
				Aliases: []string{"b"},
				Usage:   "encoder bit depth, possible values are 8, 12, 14 or 16",
				Value:   8,
				Validator: func(i uint) error {
					switch i {
					case 8, 12, 14, 16:
						return nil
					}
					return errors.New("possible values are 8, 12, 14 or 16")
				},
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "output file",
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "force override of existing output file",
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Setup the output.
			outfile := os.Stdout
			f, err := OutputFile(cmd.String("output"), cmd.Bool("force"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Incorrect Usage: %v\n\n", err)
				_ = cli.DefaultShowRootCommandHelp(cmd)
				return err
			}
			if f != nil {
				outfile = f
				defer f.Close()
			}

			// Setup the input.
			infile := os.Stdin
			switch {
			case cmd.NArg() > 1:
				fmt.Fprintf(os.Stderr, "Incorrect Usage: extra input file \"%s\"\n\n", cmd.Args().Get(1))
				_ = cli.DefaultShowRootCommandHelp(cmd)
				return errors.New("usage error")
			case cmd.NArg() == 1:
				filename := cmd.Args().Get(0)
				f, err := InputFile(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Incorrect Usage: %v\n\n", err)
					_ = cli.DefaultShowRootCommandHelp(cmd)
					return err
				}
				if f != nil {
					infile = f
					defer f.Close()
				}
			}

			var enc entropylex.Encoding
			bitDepth := cmd.Uint("bit-depth")
			switch bitDepth {
			case 8:
				enc, err = NewEntropyLex8Encoding()
				if err != nil {
					return err
				}
			case 12, 14, 16:
				return fmt.Errorf("bit depth %d unimplemented", bitDepth)
			default:
				return errors.New("usage error")
			}

			// Are we encoding or decoding?
			if cmd.Bool("decode") {
				return decode(ctx, enc, infile, outfile)
			}
			return encode(ctx, enc, infile, outfile)
		},
	}

	return cmd
}

func decode(context.Context, entropylex.Encoding, io.Reader, io.Writer) error {
	return nil
}

func encode(ctx context.Context, enc entropylex.Encoding, r io.Reader, w io.Writer) error {
	var (
		output = contextio.NewWriter(ctx, entropylex.NewEncoder(enc, w))
		input  = contextio.NewReader(ctx, r)

		buffer = make([]byte, 1024)
	)

	for {
		bytesRead, err := input.Read(buffer)

		if bytesRead > 0 {
			if _, err := output.Write(buffer); err != nil {
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
