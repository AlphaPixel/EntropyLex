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
	"errors"
	"fmt"
	"os"
)

// OutputFile creates an output file give the output filename and the force flag.
func OutputFile(filename string, force bool) (*os.File, error) {
	if filename == "" {
		return nil, nil
	}

	// See if the file exists. If not Create it.
	fs, err := os.Stat(filename)
	if err != nil {
		return os.Create(filename)
	}
	mode := fs.Mode()

	// If the file exists and it isn't a regular file, return an error.
	if !mode.IsRegular() {
		errmsg := fmt.Sprintf("output file \"%s\" is not a regular file", filename)
		if mode.IsDir() {
			errmsg = fmt.Sprintf("output file \"%s\" is a directory", filename)
		}
		return nil, errors.New(errmsg)
	}

	// If the file exists and it's a regular file, create it if the size is 0.
	if fs.Size() == 0 {
		return os.Create(filename)
	}

	// If the size is non-zero and the force flag isn't present, return an error.
	if !force {
		return nil, errors.New("a non-empty output file exists, to overwrite use the --force option")
	}

	return os.Create(filename)
}
