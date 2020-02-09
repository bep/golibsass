// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
package libsass

// #include "stdlib.h"
// #include "sass/context.h"
// #include "sass2scss.h"
import "C"
import (
	"reflect"
	"unsafe"
)

// A bridge function to C to resolve imports.
//
//export BridgeImport
func BridgeImport(currPath, prevPath *C.char, ci C.int) C.Sass_Import_List {
	parent := C.GoString(prevPath)
	rel := C.GoString(currPath)
	clist := C.sass_make_import_list(1)
	h := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(clist)),
		Len:  1, Cap: 1,
	}
	golist := *(*[]C.Sass_Import_Entry)(unsafe.Pointer(&h))

	resolver, ok := importsStore.Get(int(ci)).(ImportResolver)
	if ok {
		npath, body, ok := resolver(rel, parent)
		if ok {
			var bodyv *C.char // nil signals loading from the path.
			if body != "" {
				bodyv = C.CString(body)
			}
			entry := C.sass_make_import_entry(C.CString(npath), bodyv, nil)
			centry := (C.Sass_Import_Entry)(entry)
			golist[0] = centry
			return clist
		}
	}

	ent := C.sass_make_import_entry(currPath, nil, nil)
	cent := (C.Sass_Import_Entry)(ent)
	golist[0] = cent
	return clist
}

// SassCompilerExecute function as declared in sass/context.h:48
func SassCompilerExecute(compiler SassCompiler) {
	C.sass_compiler_execute(compiler)
}

// SassCompilerParse function as declared in sass/context.h:47
func SassCompilerParse(compiler SassCompiler) {
	C.sass_compiler_parse(compiler)
}

// SassContextGetErrorJSON function as declared in sass/context.h:115
func SassContextGetErrorJSON(ctx SassContext) string {
	s := C.sass_context_get_error_json(ctx)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// SassContextGetErrorStatus function as declared in sass/context.h:114
func SassContextGetErrorStatus(ctx SassContext) int {
	return int(C.sass_context_get_error_status(ctx))
}

func SassContextGetOutputString(ctx SassContext) string {
	s := C.sass_context_get_output_string(ctx)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// SassContextGetSourceMapString function as declared in sass/context.h:122
func SassContextGetSourceMapString(ctx SassContext) string {
	s := C.sass_context_get_source_map_string(ctx)
	return C.GoString(s)
}

// SassDataContextGetContext function as declared in sass/context.h:61
func SassDataContextGetContext(ctx SassDataContext) SassContext {
	return (SassContext)(C.sass_data_context_get_context(ctx))

}

// SassDataContextGetOptions function as declared in sass/context.h:66
func SassDataContextGetOptions(ctx SassDataContext) SassOptions {
	return (SassOptions)(C.sass_data_context_get_options(ctx))
}

// SassDataContextSetOptions function as declared in sass/context.h:68
func SassDataContextSetOptions(ctx SassDataContext, opt SassOptions) {
	C.sass_data_context_set_options(ctx, opt)
}

// SassDeleteCompiler function as declared in sass/context.h:52
func SassDeleteCompiler(compiler SassCompiler) {
	C.sass_delete_compiler(compiler)
}

// SassDeleteDataContext function as declared in sass/context.h:57
func SassDeleteDataContext(ctx SassDataContext) {
	C.sass_delete_data_context(ctx)
}

// SassDeleteFileContext function as declared in sass/context.h:56
func SassDeleteFileContext(ctx SassFileContext) {
	C.sass_delete_file_context(ctx)
}

// SassDeleteOptions function as declared in sass/context.h:53
func SassDeleteOptions(options SassOptions) {
	C.sass_delete_options(options)
}

// SassMakeDataCompiler function as declared in sass/context.h:43
func SassMakeDataCompiler(ctx SassDataContext) SassCompiler {
	return (SassCompiler)(C.sass_make_data_compiler(ctx))
}

// SassMakeDataContext function as declared in sass/context.h:35
func SassMakeDataContext(s string) SassDataContext {
	ctx := C.sass_make_data_context(C.CString(s))
	return (SassDataContext)(ctx)
}

// SassOptionGetSourceMapFile function as declared in sass/context.h:84
func SassOptionGetSourceMapFile(opts SassOptions) string {
	p := C.sass_option_get_source_map_file(opts)
	return C.GoString(p)
}

// SassOptionSetIncludePath function as declared in sass/context.h:104
func SassOptionSetIncludePath(o SassOptions, s string) {
	C.sass_option_set_include_path(o, C.CString(s))
}

// SassOptionSetInputPath function as declared in sass/context.h:101
func SassOptionSetInputPath(o SassOptions, s string) {
	C.sass_option_set_input_path(o, C.CString(s))
}

func SassOptionSetOmitSourceMapURL(o SassOptions, b bool) {
	C.sass_option_set_omit_source_map_url(o, C.bool(b))
}

// SassOptionSetOmitSourceMapUrl function as declared in sass/context.h:97
func SassOptionSetOmitSourceMapUrl(o SassOptions, b bool) {
	C.sass_option_set_omit_source_map_url(o, C.bool(b))
}

// SassOptionSetOutputPath function as declared in sass/context.h:102
func SassOptionSetOutputPath(o SassOptions, s string) {
	C.sass_option_set_output_path(o, C.CString(s))
}

// SassOptionSetOutputStyle function as declared in sass/context.h:92
func SassOptionSetOutputStyle(o SassOptions, i int) {
	C.sass_option_set_output_style(o, uint32(i))
}

// SassOptionGetPrecision function as declared in sass/context.h:91
func SassOptionSetPrecision(o SassOptions, i int) {
	C.sass_option_set_precision(o, C.int(i))
}

// SassOptionSetSourceComments function as declared in sass/context.h:93
func SassOptionSetSourceComments(o SassOptions, b bool) {
	C.sass_option_set_source_comments(o, C.bool(b))

}

// SassOptionSetSourceMapContents function as declared in sass/context.h:95
func SassOptionSetSourceMapContents(o SassOptions, b bool) {
	C.sass_option_set_source_map_contents(o, C.bool(b))
}

// SassOptionSetSourceMapEmbed function as declared in sass/context.h:94
func SassOptionSetSourceMapEmbed(o SassOptions, b bool) {
	C.sass_option_set_source_map_embed(o, C.bool(b))
}

func SassOptionSetSourceMapFile(o SassOptions, s string) {
	C.sass_option_set_source_map_file(o, C.CString(s))
}

// SassOptionSetSourceMapRoot function as declared in sass/context.h:106
func SassOptionSetSourceMapRoot(o SassOptions, s string) {
	C.sass_option_set_source_map_root(o, C.CString(s))

}

// SassToScss converts Sass to Scss using sass2scss.
func SassToScss(src string) string {
	in := C.CString(src)
	defer C.free(unsafe.Pointer(in))

	chars := C.sass2scss(
		in,
		C.int(1),
	)

	return C.GoString(chars)

}
