// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
package libsass

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCreateOutputStyle(t *testing.T) {
	c := qt.New(t)
	c.Assert(getOutputStyle("nested"), qt.Equals, NestedStyle)
	c.Assert(getOutputStyle("expanded"), qt.Equals, ExpandedStyle)
	c.Assert(getOutputStyle("compact"), qt.Equals, CompactStyle)
	c.Assert(getOutputStyle("compressed"), qt.Equals, CompressedStyle)
	c.Assert(getOutputStyle("moo"), qt.Equals, NestedStyle)

	c.Assert(getOutputStyleString(NestedStyle), qt.Equals, "nested")
	c.Assert(getOutputStyleString(ExpandedStyle), qt.Equals, "expanded")
	c.Assert(getOutputStyleString(CompactStyle), qt.Equals, "compact")
	c.Assert(getOutputStyleString(CompressedStyle), qt.Equals, "compressed")
	c.Assert(getOutputStyleString(43), qt.Equals, "nested")
}
