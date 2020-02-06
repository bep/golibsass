// Copyright © 2018 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package libsass

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/bep/go-tocss/scss"

	"github.com/stretchr/testify/require"
)

func _TestWithImportResolver(t *testing.T) {
	assert := require.New(t)
	src := bytes.NewBufferString(`
@import "colors";

div { p { color: $white; } }`)

	var dst bytes.Buffer

	importResolver := func(url string, prev string) (string, string, bool) {
		// This will make every import the same, which is probably not a common use
		// case.
		return url, `$white:    #fff`, true
	}

	transpiler, err := New(scss.Options{ImportResolver: importResolver})
	assert.NoError(err)

	_, err = transpiler.Execute(&dst, src)
	assert.NoError(err)
	assert.Equal("div p {\n  color: #fff; }\n", dst.String())
}

func TestSassSyntax(t *testing.T) {
	assert := require.New(t)
	src := bytes.NewBufferString(`
$color: #333;

.content-navigation
  border-color: $color
`)

	var dst bytes.Buffer

	transpiler, err := New(scss.Options{OutputStyle: scss.CompressedStyle, SassSyntax: true})
	assert.NoError(err)

	_, err = transpiler.Execute(&dst, src)
	assert.NoError(err)
	assert.Equal(".content-navigation{border-color:#333}\n", dst.String())
}

func TestOutputStyle(t *testing.T) {
	assert := require.New(t)
	src := bytes.NewBufferString(`
div { p { color: #ccc; } }`)

	var dst bytes.Buffer

	transpiler, err := New(scss.Options{OutputStyle: scss.CompressedStyle})
	assert.NoError(err)

	_, err = transpiler.Execute(&dst, src)
	assert.NoError(err)
	assert.Equal("div p{color:#ccc}\n", dst.String())
}

func TestSourceMapSettings(t *testing.T) {
	dir, _ := ioutil.TempDir(os.TempDir(), "tocss")
	defer os.RemoveAll(dir)

	colors := filepath.Join(dir, "_colors.scss")

	ioutil.WriteFile(colors, []byte(`
$moo:       #f442d1 !default;
`), 0755)

	assert := require.New(t)
	src := bytes.NewBufferString(`
@import "colors";

div { p { color: $moo; } }`)

	var dst bytes.Buffer

	transpiler, err := New(scss.Options{
		IncludePaths:            []string{dir},
		EnableEmbeddedSourceMap: false,
		SourceMapContents:       true,
		OmitSourceMapURL:        false,
		SourceMapFilename:       "source.map",
		OutputPath:              "outout.css",
		InputPath:               "input.scss",
		SourceMapRoot:           "/my/root",
	})
	assert.NoError(err)

	result, err := transpiler.Execute(&dst, src)
	assert.NoError(err)
	assert.Equal("div p {\n  color: #f442d1; }\n\n/*# sourceMappingURL=source.map */", dst.String())
	assert.Equal("source.map", result.SourceMapFilename)
	assert.Contains(result.SourceMapContent, `"sourceRoot": "/my/root",`)
	assert.Contains(result.SourceMapContent, `"file": "outout.css",`)
	assert.Contains(result.SourceMapContent, `"input.scss",`)
	assert.Contains(result.SourceMapContent, `mappings": "AAGA,AAAM,GAAH,CAAG,CAAC,CAAC;EAAE,KAAK,ECFH,OAAO,GDEM"`)
}

func TestConcurrentTranspile(t *testing.T) {

	assert := require.New(t)

	importResolver := func(url string, prev string) (string, string, bool) {
		return url, `$white:    #fff`, true
	}

	transpiler, err := New(scss.Options{
		OutputStyle:    scss.CompressedStyle,
		ImportResolver: importResolver})

	assert.NoError(err)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				src := bytes.NewBufferString(`
@import "colors";

div { p { color: $white; } }`)
				var dst bytes.Buffer
				_, err := transpiler.Execute(&dst, src)
				assert.NoError(err)
				assert.Equal("div p{color:#fff}\n", dst.String())
			}
		}()
	}
	wg.Wait()
}

//  3000	    397942 ns/op	    2192 B/op	       4 allocs/op
func BenchmarkTranspile(b *testing.B) {
	srcs := `div { p { color: #ccc; } }`

	var src, dst bytes.Buffer

	transpiler, err := New(scss.Options{OutputStyle: scss.CompressedStyle})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		src.Reset()
		dst.Reset()
		src.WriteString(srcs)
		if _, err := transpiler.Execute(&dst, &src); err != nil {
			b.Fatal(err)
		}
		if dst.String() != "div p{color:#ccc}\n" {
			b.Fatal("Got:", dst.String())
		}
	}
}
