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

package unicode

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/AlphaPixel/EntropyLex/src/go/arcadium"
)

var (
	reCodePoint = regexp.MustCompile(`^(?:U\+[0-9A-F]{4}|U\+[1-9A-F][0-9A-F]{4,5})$`)

	ErrInvalidCodePoint = fmt.Errorf("%w: invalid unicode code point", arcadium.ErrInternal)
)

type (
	CodePoint string
)

const (
	maxCodePoint      = 0x10FFFF
	surrogateRangeLow = 0xD800
	surrogateRangeHi  = 0xDFFF
)

func (cp CodePoint) Decode() (string, error) {
	s := string(cp)

	if s == "" {
		return "", fmt.Errorf("%w, \"\"", ErrInvalidCodePoint)
	}

	if !reCodePoint.MatchString(s) {
		return "", fmt.Errorf("%w, %s", ErrInvalidCodePoint, s)
	}

	// Yes, I know I am ignoring the error from the parse. Yes, I know this is a
	// smell that indicates that validation and parsing want to be one step.
	// However I am not going to write a json unmarhaller for this.
	i, _ := strconv.ParseUint(s[2:], 16, 32)

	switch {
	case i > maxCodePoint:
		return "", fmt.Errorf("%w, %s", ErrInvalidCodePoint, s)
	case i >= surrogateRangeLow && i <= surrogateRangeHi:
		return "", fmt.Errorf("%w, surrogate code point %s", ErrInvalidCodePoint, s)
	}
	return string(rune(i)), nil
}
