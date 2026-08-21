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

package dictionary

import (
	"embed"
	"log"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schema/*
var schemaFS embed.FS

const (
	lxjSchemaFile = "schema/lxj-v1.schema.json"
)

type schemaType int

const (
	lxjSchema schemaType = iota
)

var (
	schemas = map[schemaType]*jsonschema.Schema{}
)

func init() {
	schemas = make(map[schemaType]*jsonschema.Schema)

	var err error

	rawLXJSchema, err := schemaFS.ReadFile(lxjSchemaFile)
	if err != nil {
		log.Fatalf("Failed to read lxj schema file: %v", err)
	}

	lxjJSON, err := jsonschema.UnmarshalJSON(strings.NewReader(string(rawLXJSchema)))
	if err != nil {
		log.Fatalf("Failed to unmarshal lxj schema: %v", err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource(lxjSchemaFile, lxjJSON); err != nil {
		log.Fatal(err)
	}

	lxj, err := c.Compile(lxjSchemaFile)
	if err != nil {
		log.Fatal(err)
	}

	schemas[lxjSchema] = lxj
}
