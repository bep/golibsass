// Copyright © 2018 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package libsass a SCSS transpiler to CSS using github.com/wellington/go-libsass/libs.
package libsass

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bep/go-tocss/scss"
	"github.com/bep/go-tocss/tocss"

	//"github.com/wellington/go-libsass/libs"
	libs "github.com/bep/go-tocss/libsass"
)

type libsassTranspiler struct {
	options scss.Options
}

// New creates a new libsass transpiler configured with the given options.
func New(options scss.Options) (tocss.Transpiler, error) {
	return &libsassTranspiler{options: options}, nil
}

// Execute transpiles the SCSS or SASS from src into dst.
func (t *libsassTranspiler) Execute(dst io.Writer, src io.Reader) (tocss.Result, error) {
	var result tocss.Result
	var sourceBytes []byte

	if t.options.SassSyntax {
		// TODO1
		// LibSass does not support this directly, so have to handle the main SASS content
		// special.
		/*var buf bytes.Buffer
		err := libs.ToScss(src, &buf)
		if err != nil {
			return result, err
		}
		sourceStr = buf.String()
		*/
	} else {
		b, err := ioutil.ReadAll(src)
		if err != nil {
			return result, err
		}
		sourceBytes = b
	}

	dataCtx := libs.SassMakeDataContext(string(sourceBytes))

	opts := libs.SassDataContextGetOptions(dataCtx)

	{
		// Set options

		if t.options.ImportResolver != nil {
			//idx := libs.BindImporter(opts, t.options.ImportResolver)
			//defer libs.RemoveImporter(idx)
		}

		if t.options.Precision != 0 {
			libs.SassOptionSetPrecision(opts, t.options.Precision)
		}

		if t.options.SourceMapFilename != "" {
			libs.SassOptionSetSourceMapFile(opts, t.options.SourceMapFilename)
		}

		if t.options.SourceMapRoot != "" {
			libs.SassOptionSetSourceMapRoot(opts, t.options.SourceMapRoot)
		}

		if t.options.OutputPath != "" {
			libs.SassOptionSetOutputPath(opts, t.options.OutputPath)
		}
		if t.options.InputPath != "" {
			libs.SassOptionSetInputPath(opts, t.options.InputPath)
		}

		libs.SassOptionSetSourceMapContents(opts, t.options.SourceMapContents)
		libs.SassOptionSetOmitSourceMapURL(opts, t.options.OmitSourceMapURL)
		libs.SassOptionSetSourceMapEmbed(opts, t.options.EnableEmbeddedSourceMap)
		libs.SassOptionSetIncludePath(opts, strings.Join(t.options.IncludePaths, string(os.PathListSeparator)))
		libs.SassOptionSetOutputStyle(opts, int(t.options.OutputStyle))
		libs.SassOptionSetSourceComments(opts, false)
		libs.SassDataContextSetOptions(dataCtx, opts)
	}

	ctx := libs.SassDataContextGetContext(dataCtx)
	compiler := libs.SassMakeDataCompiler(dataCtx)

	libs.SassCompilerParse(compiler)
	libs.SassCompilerExecute(compiler)

	// TODO1 delete options
	defer libs.SassDeleteCompiler(compiler)

	outputString := libs.SassContextGetOutputString(ctx)

	io.WriteString(dst, outputString)

	if status := libs.SassContextGetErrorStatus(ctx); status != 0 {
		return result, errors.New("TODO") // result, scss.JSONToError(libs.SassContextGetErrorJSON(ctx))
	}

	result.SourceMapFilename = libs.SassOptionGetSourceMapFile(opts)
	result.SourceMapContent = libs.SassContextGetSourceMapString(ctx)

	return result, nil
}
