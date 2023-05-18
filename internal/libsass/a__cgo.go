// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package libsass

// #cgo CFLAGS: -O2 -fPIC
// #cgo CPPFLAGS: -I../../libsass_src/include
// #cgo CXXFLAGS: -g -std=c++0x -O2 -fPIC
// #cgo LDFLAGS: -lstdc++ -lm
// #cgo darwin linux LDFLAGS: -ldl
import "C"
