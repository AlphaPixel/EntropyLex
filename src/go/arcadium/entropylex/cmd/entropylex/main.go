package main

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"os"
	"path/filepath"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium/build"
	"github.com/urfave/cli/v3"
)

var (
	Version string
	Branch  string
	Commit  string
	Date    string
)

type (
	runner interface {
		Run(ctx context.Context) error
	}
)

func Main() error {
	info := build.Info(filepath.Base(os.Args[0]), Version, Branch, Commit, Date)

	name := filepath.Base(os.Args[0])

	description := "The data are encoded and decoded as described in the EntropyLex specification[1].\n\n" +
		"With no FILE, or when FILE is -, read standard input.\n" +
		"If --output/-o is not specified the command output default to stdout.\n" +
		"[1] https://github.com/AlphaPixel/EntropyLex/blob/main/SPEC.md"

	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Println(info.String())
	}
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
					case 8:
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Setup the input.
			infile := os.Stdin
			switch {
			case cmd.NArg() > 1:
				return errors.New("usage error")
			case cmd.NArg() == 1:
				filename := cmd.Args().Get(0)
				if filename != "-" {
					f, err := os.Open(filename)
					if err != nil {
						return err
					}
					infile = f
				}
			}

			// Setup the output.
			outfile := os.Stdout
			output := cmd.String("output")
			if cmd.String("output") != "" {
				f, err := os.Open(output)
				if err != nil {
					return err
				}
				outfile = f
			}

			// Are we encoding or decoding?
			decode := cmd.Bool("decode")

			var (
				el  runner
				err error
			)
			switch cmd.Uint("bit-depth") {
			case 8:
				el, err = NewEntropyLex8(infile, outfile, decode)
				if err != nil {
					return err
				}
			default:
				return errors.New("usage error")
			}

			return el.Run(ctx)
		},
	}

	return cmd.Run(context.Background(), os.Args)
}

func main() {
	if err := Main(); err != nil {
		os.Exit(1)
	}
}
