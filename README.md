
[![Build Status](https://travis-ci.org/bep/golibsass.svg?branch=master)](https://travis-ci.org/bep/golibsass)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/golibsass)](https://goreportcard.com/report/github.com/bep/golibsass)


The primary motivation for this project is to provide `SCSS` support to [Hugo](https://gohugo.io/). I welcome PRs with bug fixes. I will also consider adding functionality, but please raise an issue discussing it first.

If you need functionality than this project can provide you may want to have a look at [go-libsass](https://github.com/wellington/go-libsass).

## Update LibSASS

This project embeds the [LibSASS](https://github.com/sass/libsass) source code as a Git subtree. To update:

1. Pull in the relevant LibSASS version, e.g. `./pull-libsass.sh 3.6.3`
2. Regenerate wrappers with `go generate ./gen`

## Local development

Compiling C++ code isn' particulary fast; if you install libsass on your PC you can link against that, useful during development.

On a Mac you may do something like:

```bash
brew install --HEAD libsass
go test ./libsass -tags dev
```
