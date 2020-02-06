// +build dev

package libsass

// #cgo CPPFLAGS: -DUSE_LIBSASS_SRC
// #cgo LDFLAGS: -lsass
import "C"
