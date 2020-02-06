
[![Build Status](https://travis-ci.org/bep/go-tocss.svg?branch=master)](https://travis-ci.org/bep/go-tocss)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/go-tocss)](https://goreportcard.com/report/github.com/bep/go-tocss)


it's possible to link against system libsass and forego C compiling with go build -tags dev.

E.g.:

```bash
brew install --HEAD libsass
go test -v -run TestOutputStyle -tags dev
```


The primary motivation for this project is to add `SCSS` support to [Hugo](https://gohugo.io/). It is has some generic `tocss` package names hoping that there will be a solid native Go implementation that can replace `LibSass` in the near future.
