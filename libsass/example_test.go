// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package libsass_test

import (
	"fmt"
	"log"

	"github.com/bep/golibsass/libsass"
)

func ExampleTranspiler() {
	transpiler, err := libsass.New(libsass.Options{OutputStyle: libsass.CompressedStyle})
	if err != nil {
		log.Fatal(err)
	}

	result, err := transpiler.Execute(`
$font-stack:    Helvetica, sans-serif;
$primary-color: #333;

body {
  font: 100% $font-stack;
  color: $primary-color;
}
`)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.CSS)
	// Output: body{font:100% Helvetica,sans-serif;color:#333}
}
