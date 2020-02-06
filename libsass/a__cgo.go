package libsass

// #cgo CFLAGS: -O2 -fPIC
// #cgo CPPFLAGS: -I../libsass_src/include
// #cgo CXXFLAGS: -g -std=c++0x -O2 -fPIC
// #cgo LDFLAGS: -lstdc++ -lm
// #cgo darwin linux LDFLAGS: -ldl
//
import "C"
