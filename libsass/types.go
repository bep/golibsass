// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
package libsass

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Result struct {
	CSS string

	// If source maps are configured.
	SourceMapFilename string
	SourceMapContent  string
}

type Transpiler interface {
	Execute(src string) (Result, error)
}

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

var outputStyles = map[string]OutputStyle{
	nestedStyleStr:     NestedStyle,
	expandedStyleStr:   ExpandedStyle,
	compactStyleStr:    CompactStyle,
	compressedStyleStr: CompressedStyle,
}

var outputStylesString = map[OutputStyle]string{
	NestedStyle:     nestedStyleStr,
	ExpandedStyle:   expandedStyleStr,
	CompactStyle:    compactStyleStr,
	CompressedStyle: compressedStyleStr,
}

// TODO1 add a method printing the LibSass version.

func getOutputStyle(style string) OutputStyle {
	os, found := outputStyles[strings.ToLower(style)]
	if found {
		return os
	}
	return NestedStyle
}

func getOutputStyleString(style OutputStyle) string {
	os, found := outputStylesString[style]
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

func jsonToError(jsonstr string) (e Error) {
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
