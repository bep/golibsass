//+build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func Gen() error {
	dstDir := filepath.FromSlash("internal/libsass")
	srcDir := filepath.Join("libsass_src", "src")

	// The Go and the Libsass C++ source must live side-by-side in the same
	// directory.
	//
	// The custom bindings are named with a "a__" prefix. Keep those.
	fis, err := ioutil.ReadDir(dstDir)
	if err != nil {
		return err
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
		return err
	}

	csourceRe := regexp.MustCompile(`\.[ch](pp)?$`)

	for _, fi := range fis {
		if fi.IsDir() || !csourceRe.MatchString(fi.Name()) {
			continue
		}

		target := filepath.Join(dstDir, fi.Name())

		if err := ioutil.WriteFile(target, []byte(fmt.Sprintf(`#ifndef USE_LIBSASS_SRC
#include "src/%s"
#endif
`, fi.Name())), 0644); err != nil {
			return err
		}
	}

	return nil
}
