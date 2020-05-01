
[![Build Status](https://travis-ci.org/bep/golibsass.svg?branch=master)](https://travis-ci.org/bep/golibsass)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/golibsass)](https://goreportcard.com/report/github.com/bep/golibsass)
[![LibSass Version](https://img.shields.io/badge/LibSass-v3.6.4-blue)](https://github.com/sass/libsass)
[![codecov](https://codecov.io/gh/bep/golibsass/branch/master/graph/badge.svg)](https://codecov.io/gh/bep/golibsass)
[![GoDoc](https://godoc.org/github.com/bep/golibsass/libsass?status.svg)](https://godoc.org/github.com/bep/golibsass/libsass)

The primary motivation for this project is to provide `SCSS` support to [Hugo](https://gohugo.io/). I welcome PRs with bug fixes. I will also consider adding functionality, but please raise an issue discussing it first.

If you need more functionality than this project can provide you may want to have a look at [go-libsass](https://github.com/wellington/go-libsass).

## Usage

A basic example (error handling omitted):

```go
transpiler, _ := libsass.New(libsass.Options{OutputStyle: libsass.CompressedStyle})
	
result, _ := transpiler.Execute(`
$font-stack:    Helvetica, sans-serif;
$primary-color: #333;

body {
  font: 100% $font-stack;
  color: $primary-color;
}
`)

fmt.Println(result.CSS)
// Output: body{font:100% Helvetica,sans-serif;color:#333}
```

See the [GoDoc](https://godoc.org/github.com/bep/golibsass/libsass) for more options.

## Update LibSass version

This project embeds the [LibSASS](https://github.com/sass/libsass) source code as a Git subtree. To update:

1. Pull in the relevant LibSASS version, e.g. `./pull-libsass.sh 3.6.3`
2. Regenerate wrappers with `go generate ./gen`
3. Update the LibSass version badge above.

## Local development

Compiling C++ code isn' particulary fast; if you install libsass on your PC you can link against that, useful during development.

On a Mac you may do something like:

```bash
brew install --HEAD libsass
go test ./libsass -tags dev
```
