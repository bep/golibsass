// Copyright © 2018 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package scss provides options for SCSS transpilers. Note that there are no
// current pure Go SASS implementation, so the only option is CGO and LibSASS.
// But hopefully, fingers crossed, this will happen.
package scss

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOutputStyle(t *testing.T) {
	assert := require.New(t)
	assert.Equal(NestedStyle, OutputStyleFromString("nested"))
	assert.Equal(ExpandedStyle, OutputStyleFromString("expanded"))
	assert.Equal(CompactStyle, OutputStyleFromString("compact"))
	assert.Equal(CompressedStyle, OutputStyleFromString("compressed"))
	assert.Equal(NestedStyle, OutputStyleFromString("moo"))

	assert.Equal("nested", OutputStyleToString(NestedStyle))
	assert.Equal("expanded", OutputStyleToString(ExpandedStyle))
	assert.Equal("compact", OutputStyleToString(CompactStyle))
	assert.Equal("compressed", OutputStyleToString(CompressedStyle))
	assert.Equal("nested", OutputStyleToString(43))
}
