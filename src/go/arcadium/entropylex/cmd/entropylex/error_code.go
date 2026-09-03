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

type (
	ErrorCode struct {
		msg  string
		code int
	}
)

func (err ErrorCode) Error() string {
	return err.msg
}

func (err ErrorCode) Code() int {
	return err.code
}

const (
	UsageErrorCode         int = 1
	InternalErrorCode      int = 4
	UnknownErrorCode       int = 9
	UnimplementedErrorCode int = 111
)

var (
	ErrUsage         = ErrorCode{msg: "usage error", code: UsageErrorCode}
	ErrInternal      = ErrorCode{msg: "internal error", code: InternalErrorCode}
	ErrUnimplemented = ErrorCode{msg: "unimplemented", code: UnimplementedErrorCode}
)
