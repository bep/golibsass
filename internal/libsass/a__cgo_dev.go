// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
// +build dev

package libsass

// #cgo CPPFLAGS: -DUSE_LIBSASS_SRC
// #cgo LDFLAGS: -lsass
import "C"
