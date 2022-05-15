// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package libsass

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/bep/golibsass/libsass/libsasserrors"
	qt "github.com/frankban/quicktest"
)

const (
	sassSample = `nav {
  ul {
    margin: 0;
    padding: 0;
    list-style: none;
  }

  li { display: inline-block; }

  a {
    display: block;
    padding: 6px 12px;
    text-decoration: none;
  }
}`
	sassSampleTranspiled = "nav ul {\n  margin: 0;\n  padding: 0;\n  list-style: none; }\n\nnav li {\n  display: inline-block; }\n\nnav a {\n  display: block;\n  padding: 6px 12px;\n  text-decoration: none; }\n"
)

func TestTranspiler(t *testing.T) {
	c := qt.New(t)

	importResolver := func(url string, prev string) (string, string, bool) {
		// This will make every import the same, which is probably not a common use
		// case.
		return url, `$white:    #fff`, true
	}

	for _, test := range []struct {
		name   string
		opts   Options
		src    string
		expect interface{}
	}{
		{"Output style compressed", Options{OutputStyle: CompressedStyle}, "div { color: #ccc; }", "div{color:#ccc}\n"},
		{"Invalid syntax", Options{OutputStyle: CompressedStyle}, "div { color: $white; }", false},
		{"Import not found", Options{OutputStyle: CompressedStyle}, "@import \"foo\"", false},
		{"Sass syntax", Options{OutputStyle: CompressedStyle, SassSyntax: true}, "$color: #ccc\ndiv { p { color: $color; } }", "div p{color:#ccc}\n"},
		{"Import resolver", Options{ImportResolver: importResolver}, "@import \"colors\";\ndiv { p { color: $white; } }", "div p {\n  color: #fff; }\n"},
		{"Precision", Options{Precision: 3}, "div { width: percentage(1 / 3); }", "div {\n  width: 33.333%; }\n"},
	} {

		test := test
		c.Run(test.name, func(c *qt.C) {
			b, ok := test.expect.(bool)
			shouldFail := ok && !b

			transpiler, err := New(test.opts)
			c.Assert(err, qt.IsNil)
			result, err := transpiler.Execute(test.src)
			if shouldFail {
				c.Assert(err, qt.Not(qt.IsNil))
			} else {
				c.Assert(err, qt.IsNil)
				c.Assert(result.CSS, qt.Equals, test.expect)
			}
		})

	}
}

func TestError(t *testing.T) {
	c := qt.New(t)
	transpiler, err := New(Options{OutputStyle: CompressedStyle})
	c.Assert(err, qt.IsNil)
	_, err = transpiler.Execute("\n\ndiv { color: $blue; }")
	c.Assert(err, qt.Not(qt.IsNil))

	lerr := err.(libsasserrors.Error)
	c.Assert(lerr.Line, qt.Equals, 3)
	c.Assert(lerr.Column, qt.Equals, 14)
	c.Assert(lerr.Message, qt.Equals, `Undefined variable: "$blue".`)
	c.Assert(lerr.Error(), qt.Equals, `file "stdin", line 3, col 14: Undefined variable: "$blue". `)
}

func TestSourceMapSettings(t *testing.T) {
	c := qt.New(t)
	src := `div { p { color: blue; } }`

	transpiler, err := New(Options{
		SourceMapOptions: SourceMapOptions{
			EnableEmbedded: false,
			Contents:       true,
			OmitURL:        false,
			Filename:       "source.map",
			OutputPath:     "outout.css",
			InputPath:      "input.scss",
			Root:           "/my/root",
		},
	})
	c.Assert(err, qt.IsNil)

	result, err := transpiler.Execute(src)
	c.Assert(err, qt.IsNil)
	c.Assert(result.CSS, qt.Equals, "div p {\n  color: blue; }\n\n/*# sourceMappingURL=source.map */")
	c.Assert(result.SourceMapFilename, qt.Equals, "source.map")
	c.Assert(`"sourceRoot": "/my/root",`, qt.Contains, `"sourceRoot": "/my/root",`)
	c.Assert(`"file": "outout.css",`, qt.Contains, `"file": "outout.css",`)
	c.Assert(`"input.scss",`, qt.Contains, `"input.scss",`)
	c.Assert(`mappings": "AAGA,AAAM,GAAH,CAAG,CAAC,CAAC;EAAE,KAAK,ECFH,OAAO,GDEM"`, qt.Contains, `mappings": "AAGA,AAAM,GAAH,CAAG,CAAC,CAAC;EAAE,KAAK,ECFH,OAAO,GDEM"`)
}

func TestIncludePaths(t *testing.T) {
	dir1, _ := ioutil.TempDir(os.TempDir(), "libsass-test-include-paths-dir1")
	defer os.RemoveAll(dir1)
	dir2, _ := ioutil.TempDir(os.TempDir(), "libsass-test-include-paths-dir2")
	defer os.RemoveAll(dir2)

	colors := filepath.Join(dir1, "_colors.scss")
	content := filepath.Join(dir2, "_content.scss")

	ioutil.WriteFile(colors, []byte(`
$moo:       #f442d1 !default;
`), 0o644)

	ioutil.WriteFile(content, []byte(`
content { color: #ccc; }
`), 0o644)

	c := qt.New(t)
	src := `
@import "colors";
@import "content";
div { p { color: $moo; } }`

	transpiler, err := New(Options{
		IncludePaths: []string{dir1, dir2},
		OutputStyle:  CompressedStyle,
		ImportResolver: func(url string, prev string) (newUrl string, body string, resolved bool) {
			// Let LibSass resolve the import.
			return "", "", false
		},
	})
	c.Assert(err, qt.IsNil)

	result, err := transpiler.Execute(src)
	c.Assert(err, qt.IsNil)
	c.Assert(result.CSS, qt.Equals, "content{color:#ccc}div p{color:#f442d1}\n")
}

func TestConcurrentTranspile(t *testing.T) {
	c := qt.New(t)

	importResolver := func(url string, prev string) (string, string, bool) {
		return url, `$white:    #fff`, true
	}

	transpiler, err := New(Options{
		OutputStyle:    CompressedStyle,
		ImportResolver: importResolver,
	})

	c.Assert(err, qt.IsNil)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				src := `
@import "colors";

div { p { color: $white; } }`
				result, err := transpiler.Execute(src)
				c.Check(err, qt.IsNil)
				c.Check(result.CSS, qt.Equals, "div p{color:#fff}\n")
				if c.Failed() {
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestImportResolverConcurrent(t *testing.T) {
	c := qt.New(t)

	createImportResolver := func(width int) func(url string, prev string) (string, string, bool) {
		return func(url string, prev string) (string, string, bool) {
			return url, fmt.Sprintf(`$width:  %d`, width), true
		}
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				transpiler, err := New(Options{
					OutputStyle:    CompressedStyle,
					ImportResolver: createImportResolver(j),
				})
				c.Check(err, qt.IsNil)

				src := `
@import "widths";

div { p { width: $width; } }`

				for k := 0; k < 10; k++ {
					result, err := transpiler.Execute(src)
					c.Check(err, qt.IsNil)
					c.Check(result.CSS, qt.Equals, fmt.Sprintf("div p{width:%d}\n", j))
					if c.Failed() {
						return
					}
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkTranspile(b *testing.B) {
	type tester struct {
		src        string
		expect     string
		transpiler Transpiler
	}

	newTester := func(b *testing.B, opts Options) tester {
		transpiler, err := New(opts)
		if err != nil {
			b.Fatal(err)
		}

		return tester{
			transpiler: transpiler,
		}
	}

	runBench := func(b *testing.B, t tester) {
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			result, err := t.transpiler.Execute(t.src)
			if err != nil {
				b.Fatal(err)
			}
			if result.CSS != t.expect {
				b.Fatal("Got:", result.CSS)
			}
		}
	}

	b.Run("SCSS", func(b *testing.B) {
		t := newTester(b, Options{})
		t.src = sassSample
		t.expect = sassSampleTranspiled
		runBench(b, t)
	})

	b.Run("SCSS Parallel", func(b *testing.B) {
		t := newTester(b, Options{})
		t.src = sassSample
		t.expect = sassSampleTranspiled

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				result, err := t.transpiler.Execute(t.src)
				if err != nil {
					b.Fatal(err)
				}
				if result.CSS != t.expect {
					b.Fatalf("Got: %q\n", result.CSS)
				}
			}
		})
	})

	b.Run("Sass", func(b *testing.B) {
		t := newTester(b, Options{OutputStyle: CompressedStyle, SassSyntax: true})
		t.src = `
$color: #333;

.content-navigation
  border-color: $color`

		t.expect = ".content-navigation{border-color:#333}\n"
		runBench(b, t)
	})
}

func TestParseOutputStyle(t *testing.T) {
	c := qt.New(t)
	c.Assert(ParseOutputStyle("nested"), qt.Equals, NestedStyle)
	c.Assert(ParseOutputStyle("expanded"), qt.Equals, ExpandedStyle)
	c.Assert(ParseOutputStyle("compact"), qt.Equals, CompactStyle)
	c.Assert(ParseOutputStyle("compressed"), qt.Equals, CompressedStyle)
	c.Assert(ParseOutputStyle("EXPANDED"), qt.Equals, ExpandedStyle)
	c.Assert(ParseOutputStyle("foo"), qt.Equals, NestedStyle)
}
