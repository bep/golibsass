// Copyright © 2018 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package tocss provides the common API for tranpilation of the source
// to CSS.
package tocss

import (
	"io"
)

type Result struct {
	// If source maps are configured.
	SourceMapFilename string
	SourceMapContent  string
}

type Transpiler interface {
	Execute(dst io.Writer, src io.Reader) (Result, error)
}
