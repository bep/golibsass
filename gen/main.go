//go:generate go run main.go

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime err")
	}

	rootDir := path.Join(path.Dir(filename), "..")

	dstDir := filepath.Join(rootDir, "internal/libsass")
	srcDir := filepath.Join(rootDir, "libsass_src", "src")

	// The Go and the Libsass C++ source must live side-by-side in the same
	// directory.
	//
	// The custom bindings are named with a "a__" prefix. Keep those.
	fis, err := ioutil.ReadDir(dstDir)
	if err != nil {
		log.Fatal(err)
	}

	keepRe := regexp.MustCompile(`^(a__|\.)`)

	for _, fi := range fis {
		if keepRe.MatchString(fi.Name()) {
			continue
		}
		os.Remove(filepath.Join(dstDir, fi.Name()))
	}

	fis, err = ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	csourceRe := regexp.MustCompile(`\.[ch](pp)?$`)

	for _, fi := range fis {
		if fi.IsDir() || !csourceRe.MatchString(fi.Name()) {
			continue
		}

		target := filepath.Join(dstDir, fi.Name())

		if err := ioutil.WriteFile(target, []byte(fmt.Sprintf(`#ifndef USE_LIBSASS_SRC
#include "../../libsass_src/src/%s"
#endif
`, fi.Name())), 0o644); err != nil {
			log.Fatal(err)
		}
	}
}
