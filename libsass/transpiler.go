// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package libsass a SCSS transpiler to CSS using LibSASS.
package libsass

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bep/golibsass/internal/libsass"
)

type libsassTranspiler struct {
	options Options
}

// New creates a new libsass transpiler configured with the given options.
func New(options Options) (Transpiler, error) {
	return &libsassTranspiler{options: options}, nil
}

// Execute transpiles the SCSS or SASS from src into dst.
func (t *libsassTranspiler) Execute(dst io.Writer, src io.Reader) (Result, error) {
	var result Result
	var sourceBytes []byte

	if t.options.SassSyntax {
		// LibSass does not support this directly, so have to handle the main SASS content
		// special.
		var buf bytes.Buffer
		err := libsass.SassToScss(&buf, src)
		if err != nil {
			return result, err
		}
		sourceBytes = buf.Bytes()
	} else {
		b, err := ioutil.ReadAll(src)
		if err != nil {
			return result, err
		}
		sourceBytes = b
	}

	dataCtx := libsass.SassMakeDataContext(string(sourceBytes))

	opts := libsass.SassDataContextGetOptions(dataCtx)
	defer libsass.SassDeleteOptions(opts)

	{
		// Set options

		if t.options.ImportResolver != nil {
			idx := libsass.AddImportResolver(opts, t.options.ImportResolver)
			defer libsass.DeleteImportResolver(idx)
		}

		if t.options.Precision != 0 {
			libsass.SassOptionSetPrecision(opts, t.options.Precision)
		}

		if t.options.SourceMapFilename != "" {
			libsass.SassOptionSetSourceMapFile(opts, t.options.SourceMapFilename)
		}

		if t.options.SourceMapRoot != "" {
			libsass.SassOptionSetSourceMapRoot(opts, t.options.SourceMapRoot)
		}

		if t.options.OutputPath != "" {
			libsass.SassOptionSetOutputPath(opts, t.options.OutputPath)
		}
		if t.options.InputPath != "" {
			libsass.SassOptionSetInputPath(opts, t.options.InputPath)
		}

		libsass.SassOptionSetSourceMapContents(opts, t.options.SourceMapContents)
		libsass.SassOptionSetOmitSourceMapURL(opts, t.options.OmitSourceMapURL)
		libsass.SassOptionSetSourceMapEmbed(opts, t.options.EnableEmbeddedSourceMap)
		libsass.SassOptionSetIncludePath(opts, strings.Join(t.options.IncludePaths, string(os.PathListSeparator)))
		libsass.SassOptionSetOutputStyle(opts, int(t.options.OutputStyle))
		libsass.SassOptionSetSourceComments(opts, false)
		libsass.SassDataContextSetOptions(dataCtx, opts)
	}

	ctx := libsass.SassDataContextGetContext(dataCtx)
	compiler := libsass.SassMakeDataCompiler(dataCtx)
	defer libsass.SassDeleteCompiler(compiler)

	libsass.SassCompilerParse(compiler)
	libsass.SassCompilerExecute(compiler)

	outputString := libsass.SassContextGetOutputString(ctx)

	io.WriteString(dst, outputString)

	if status := libsass.SassContextGetErrorStatus(ctx); status != 0 {
		return result, jsonToError(libsass.SassContextGetErrorJSON(ctx))
	}

	result.SourceMapFilename = libsass.SassOptionGetSourceMapFile(opts)
	result.SourceMapContent = libsass.SassContextGetSourceMapString(ctx)

	return result, nil
}
