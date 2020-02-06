package libsass

// #cgo CPPFLAGS: -DUSE_LIBSASS_SRC
// #cgo CPPFLAGS: -I../libsass_src/include
// #cgo LDFLAGS: -lsass
// #cgo darwin linux LDFLAGS: -ldl
//
import "C"
