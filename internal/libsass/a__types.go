// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package libsass

// #include "stdint.h"
// #include "stdlib.h"
// #include "sass/context.h"
import "C"
import "unsafe"

// SassValue as declared in sass/values.h:14
const sizeofSassValue = unsafe.Sizeof(C.union_Sass_Value{})

// SassCalleeEntry as declared in sass/functions.h:25
type SassCalleeEntry C.Sass_Callee_Entry

// SassCompiler as declared in sass/functions.h:18
type SassCompiler *C.struct_Sass_Compiler

// SassContext as declared in sass/context.h:20
type SassContext *C.struct_Sass_Context

// SassDataContext as declared in sass/context.h:22
type SassDataContext *C.struct_Sass_Data_Context

// SassEnvFrame as declared in sass/functions.h:23
type SassEnvFrame C.Sass_Env_Frame

// SassFileContext as declared in sass/context.h:21
type SassFileContext *C.struct_Sass_File_Context

// SassFunctionEntry as declared in sass/functions.h:37
type SassFunctionEntry C.Sass_Function_Entry

// SassFunctionList as declared in sass/functions.h:38
type SassFunctionList C.Sass_Function_List

// SassImportEntry as declared in sass/functions.h:27
type SassImportEntry C.Sass_Import_Entry

// SassImportList as declared in sass/functions.h:28
type SassImportList C.Sass_Import_List

// SassImporterEntry as declared in sass/functions.h:30
type SassImporterEntry C.Sass_Importer_Entry

// SassImporterFn type as declared in sass/functions.h:33
type SassImporterFn func(url string, cb SassImporterEntry, compiler SassCompiler) SassImportList

// SassImporterList as declared in sass/functions.h:31
type SassImporterList C.Sass_Importer_List

// SassOptions as declared in sass/functions.h:17
type SassOptions *C.struct_Sass_Options

type SassValue [sizeofSassValue]byte
