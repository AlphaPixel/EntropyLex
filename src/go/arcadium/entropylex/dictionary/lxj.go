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

// Package dictionary provides dictionary implementations for EntropyLex.
// Currently, it supports the LXJ dictionary, which is a JSON-based format.
package dictionary

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type (
	LXJ struct {
		Format      Format      `json:"format"`
		Dictionary  Dictionary  `json:"dictionary"`
		Profile     Profile     `json:"profile"`
		Recognition Recognition `json:"recognition"`
		Indexing    string      `json:"indexing,omitempty"`
		Tokens      Tokens      `json:"tokens"`
		Fingerprint Fingerprint `json:"fingerprint"`
		Provenance  Provenance  `json:"provenance"`
	}

	Format struct {
		Name    string `json:"name,omitempty"`
		Version int    `json:"version,omitempty"`
	}

	Dictionary struct {
		ID       string `json:"id,omitempty"`
		Version  int    `json:"version,omitempty"`
		Purpose  string `json:"purpose,omitempty"`
		Name     string `json:"name,omitempty"`
		Language string `json:"language,omitempty"`
		Script   string `json:"script,omitempty"`
	}

	Profile struct {
		Name          string `json:"name,omitempty"`
		NormalBits    int    `json:"normal_bits,omitempty"`
		RemainderBits []int  `json:"remainder_bits,omitempty"`
		NormalCount   int    `json:"normal_count,omitempty"`
		TrimCount     int    `json:"trim_count,omitempty"`
		TotalCount    int    `json:"total_count,omitempty"`
	}

	Recognition struct {
		UnicodeVersion string       `json:"unicode_version,omitempty"`
		Normalization  string       `json:"normalization,omitempty"`
		Case           string       `json:"case_folding,omitempty"`
		TokenText      string       `json:"token_text,omitempty"`
		Tokenization   Tokenization `json:"tokenization"`
	}

	Tokenization struct {
		Kind               string   `json:"kind,omitempty"`
		CanonicalSeparator string   `json:"canonical_separator,omitempty"`
		Separators         []string `json:"separators,omitempty"`
	}

	Tokens []string

	Fingerprint struct {
		Recipe string `json:"recipe,omitempty"`
		SHA256 string `json:"sha256,omitempty"`
	}

	Provenance struct {
		Sources   []Source  `json:"sources,omitempty"`
		Selection Selection `json:"selection"`
	}

	Source struct {
		Name        string `json:"name,omitempty"`
		Version     string `json:"version,omitempty"`
		Role        string `json:"role,omitempty"`
		Location    string `json:"location,omitempty"`
		License     string `json:"license,omitempty"`
		SHA256      string `json:"sha256,omitempty"`
		RecordCount int    `json:"record_count,omitempty"`
	}

	Selection struct {
		Method           string `json:"method,omitempty"`
		SettingsLocation string `json:"settings_location,omitempty"`
		SettingsSHA256   string `json:"settings_sha256,omitempty"`
	}
)

// NewLXJ creates a new LXJ dictionary from an io.Reader.
func NewLXJ(r io.ReadSeeker) (*LXJ, error) {
	if r == nil {
		return nil, errors.New("invalid dictionary file")
	}

	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	var dict LXJ
	if err := json.NewDecoder(r).Decode(&dict); err != nil {
		return nil, err
	}

	return &dict, nil
}

// NewLXJValidated creates a new LXJ dictionary from an io.Reader and validates
// it against the LXJ schema.
func NewLXJValidated(r io.ReadSeeker) (*LXJ, error) {
	if r == nil {
		return nil, errors.New("invalid dictionary file")
	}

	// Validate the JSON file against the LXJ schema.
	dict, err := jsonschema.UnmarshalJSON(r)
	if err != nil {
		return nil, err
	}
	if err := schemas[lxjSchema].Validate(dict); err != nil {
		return nil, err
	}

	return NewLXJ(r)
}
