// Copyright © 2018 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package scss provides options for SCSS transpilers. Note that there are no
// current pure Go SASS implementation, so the only option is CGO and LibSASS.
// But hopefully, fingers crossed, this will happen.
package scss

import (
	"encoding/json"
	"fmt"
	"strings"
)

type (
	OutputStyle int
)

const (
	NestedStyle OutputStyle = iota
	ExpandedStyle
	CompactStyle
	CompressedStyle
)

const (
	nestedStyleStr     = "nested"
	expandedStyleStr   = "expanded"
	compactStyleStr    = "compact"
	compressedStyleStr = "compressed"
)

var outputStyleFromString = map[string]OutputStyle{
	nestedStyleStr:     NestedStyle,
	expandedStyleStr:   ExpandedStyle,
	compactStyleStr:    CompactStyle,
	compressedStyleStr: CompressedStyle,
}

var outputStyleToString = map[OutputStyle]string{
	NestedStyle:     nestedStyleStr,
	ExpandedStyle:   expandedStyleStr,
	CompactStyle:    compactStyleStr,
	CompressedStyle: compressedStyleStr,
}

func OutputStyleFromString(style string) OutputStyle {
	os, found := outputStyleFromString[strings.ToLower(style)]
	if found {
		return os
	}
	return NestedStyle
}

func OutputStyleToString(style OutputStyle) string {
	os, found := outputStyleToString[style]
	if found {
		return os
	}
	return nestedStyleStr
}

type Options struct {
	// Default is nested.
	OutputStyle OutputStyle

	// Precision of floating point math.
	Precision int

	// File paths to use to resolve imports.
	IncludePaths []string

	// ImportResolver can be used to supply a custom import resolver, both to redirect
	// to another URL or to return the body.
	ImportResolver func(url string, prev string) (newURL string, body string, resolved bool)

	// Used to indicate "old style" SASS for the input stream.
	SassSyntax bool

	// Source map settings
	SourceMapFilename       string
	SourceMapRoot           string
	InputPath               string
	OutputPath              string
	SourceMapContents       bool
	OmitSourceMapURL        bool
	EnableEmbeddedSourceMap bool
}

func JSONToError(jsonstr string) (e Error) {
	if err := json.Unmarshal([]byte(jsonstr), &e); err != nil {
		e.Message = "unknown error"
	}
	return
}

type Error struct {
	Status  int    `json:"status"`
	Column  int    `json:"column"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("file %q, line %d, col %d: %s ", e.File, e.Line, e.Column, e.Message)
}
