// TODO(bep)

package libsass

/*
#cgo CFLAGS: -O2 -fPIC
#cgo CPPFLAGS: -I../libsass_src/include
#cgo CXXFLAGS: -g -std=c++0x -O2 -fPIC
#cgo LDFLAGS: -lstdc++ -lm
#cgo darwin linux LDFLAGS: -ldl
#include "stdint.h"
#include "sass/context.h"
#include <stdlib.h>
*/
import "C"

// SassCompilerState as declared in sass/context.h:25
type SassCompilerState int32

// SassCompilerState enumeration from sass/context.h:25
const ()

// SassOutputStyle as declared in sass/base.h:64
type SassOutputStyle int32

// SassOutputStyle enumeration from sass/base.h:64
const ()

// SassTag as declared in sass/values.h:17
type SassTag int32

// SassTag enumeration from sass/values.h:17
const ()

// SassSeparator as declared in sass/values.h:30
type SassSeparator int32

// SassSeparator enumeration from sass/values.h:30
const ()

// SassOp as declared in sass/values.h:39
type SassOp int32

// SassOp enumeration from sass/values.h:39
const ()

// SassCalleeType as declared in sass/functions.h:44
type SassCalleeType int32

// SassCalleeType enumeration from sass/functions.h:44
const ()
