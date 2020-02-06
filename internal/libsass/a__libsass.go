// TODO(bep)

package libsass

// #include "stdint.h"

// extern struct Sass_Import** HeaderBridge(uintptr_t idx);

// #include "stdlib.h"
// #include "sass/context.h"
// #include "sass2scss.h"
// #include <stdlib.h>
//
import "C"
import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"unsafe"
)

// SassMakeDataContext function as declared in sass/context.h:35
func SassMakeDataContext(s string) SassDataContext {
	ctx := C.sass_make_data_context(C.CString(s))
	return (SassDataContext)(ctx)
}

// SassDataContextGetOptions function as declared in sass/context.h:66
func SassDataContextGetOptions(ctx SassDataContext) SassOptions {
	return (SassOptions)(C.sass_data_context_get_options(ctx))
}

// SassOptionGetPrecision function as declared in sass/context.h:91
func SassOptionSetPrecision(o SassOptions, i int) {
	C.sass_option_set_precision(o, C.int(i))
}

// SassOptionSetSourceMapRoot function as declared in sass/context.h:106
func SassOptionSetSourceMapRoot(o SassOptions, s string) {
	C.sass_option_set_source_map_root(o, C.CString(s))

}

// SassOptionSetOutputPath function as declared in sass/context.h:102
func SassOptionSetOutputPath(o SassOptions, s string) {
	C.sass_option_set_output_path(o, C.CString(s))
}

func SassOptionSetOmitSourceMapURL(o SassOptions, b bool) {
	C.sass_option_set_omit_source_map_url(o, C.bool(b))
}

// SassOptionSetInputPath function as declared in sass/context.h:101
func SassOptionSetInputPath(o SassOptions, s string) {
	C.sass_option_set_input_path(o, C.CString(s))
}

// SassOptionSetSourceMapContents function as declared in sass/context.h:95
func SassOptionSetSourceMapContents(o SassOptions, b bool) {
	C.sass_option_set_source_map_contents(o, C.bool(b))
}

func SassOptionSetSourceMapFile(o SassOptions, s string) {
	C.sass_option_set_source_map_file(o, C.CString(s))
}

// SassOptionSetOmitSourceMapUrl function as declared in sass/context.h:97
func SassOptionSetOmitSourceMapUrl(o SassOptions, b bool) {
	C.sass_option_set_omit_source_map_url(o, C.bool(b))
}

// SassOptionSetSourceMapEmbed function as declared in sass/context.h:94
func SassOptionSetSourceMapEmbed(o SassOptions, b bool) {
	C.sass_option_set_source_map_embed(o, C.bool(b))
}

// SassOptionSetIncludePath function as declared in sass/context.h:104
func SassOptionSetIncludePath(o SassOptions, s string) {
	C.sass_option_set_include_path(o, C.CString(s))
}

// SassOptionSetOutputStyle function as declared in sass/context.h:92
func SassOptionSetOutputStyle(o SassOptions, i int) {
	C.sass_option_set_output_style(o, uint32(i))
}

// SassOptionSetSourceComments function as declared in sass/context.h:93
func SassOptionSetSourceComments(o SassOptions, b bool) {
	C.sass_option_set_source_comments(o, C.bool(b))

}

// SassDataContextSetOptions function as declared in sass/context.h:68
func SassDataContextSetOptions(ctx SassDataContext, opt SassOptions) {
	C.sass_data_context_set_options(ctx, opt)
}

// SassDataContextGetContext function as declared in sass/context.h:61
func SassDataContextGetContext(ctx SassDataContext) SassContext {
	return (SassContext)(C.sass_data_context_get_context(ctx))

}

// SassMakeDataCompiler function as declared in sass/context.h:43
func SassMakeDataCompiler(ctx SassDataContext) SassCompiler {
	return (SassCompiler)(C.sass_make_data_compiler(ctx))
}

// SassCompilerParse function as declared in sass/context.h:47
func SassCompilerParse(compiler SassCompiler) {
	C.sass_compiler_parse(compiler)
}

// SassCompilerExecute function as declared in sass/context.h:48
func SassCompilerExecute(compiler SassCompiler) {
	C.sass_compiler_execute(compiler)
}

// SassDeleteCompiler function as declared in sass/context.h:52
func SassDeleteCompiler(compiler SassCompiler) {
	C.sass_delete_compiler(compiler)
}

// SassDeleteOptions function as declared in sass/context.h:53
func SassDeleteOptions(options SassOptions) {
	C.sass_delete_options(options)
}

// SassDeleteFileContext function as declared in sass/context.h:56
func SassDeleteFileContext(ctx SassFileContext) {
	C.sass_delete_file_context(ctx)
}

// SassDeleteDataContext function as declared in sass/context.h:57
func SassDeleteDataContext(ctx SassDataContext) {
	C.sass_delete_data_context(ctx)
}

func SassContextGetOutputString(ctx SassContext) string {
	s := C.sass_context_get_output_string(ctx)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// SassContextGetErrorStatus function as declared in sass/context.h:114
func SassContextGetErrorStatus(ctx SassContext) int {
	return int(C.sass_context_get_error_status(ctx))
}

// SassContextGetErrorJSON function as declared in sass/context.h:115
func SassContextGetErrorJSON(ctx SassContext) string {
	s := C.sass_context_get_error_json(ctx)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// SassOptionGetSourceMapFile function as declared in sass/context.h:84
func SassOptionGetSourceMapFile(opts SassOptions) string {
	p := C.sass_option_get_source_map_file(opts)
	return C.GoString(p)
}

// SassContextGetSourceMapString function as declared in sass/context.h:122
func SassContextGetSourceMapString(ctx SassContext) string {
	s := C.sass_context_get_source_map_string(ctx)
	return C.GoString(s)
}

// SassToScss converts Sass to Scss using sass2scss.
func SassToScss(dst io.Writer, src io.Reader) error {
	b, _ := ioutil.ReadAll(src)
	in := C.CString(string(b))
	defer C.free(unsafe.Pointer(in))

	chars := C.sass2scss(
		in,
		C.int(1),
	)
	_, err := io.WriteString(dst, C.GoString(chars))
	return err
}

// A bridge function to C to resolve imports.
//
//export ResolveImportBridge
func ResolveImportBridge(url *C.char, prev *C.char, cidx C.uintptr_t) C.Sass_Import_List {
	var importResolver ImportResolver

	// Retrieve the index
	idx := int(cidx)
	importResolver, ok := imports.Get(idx).(ImportResolver)

	if !ok {
		fmt.Printf("failed to resolve import handler: %d\n", idx)
	}

	parent := C.GoString(prev)
	rel := C.GoString(url)
	list := C.sass_make_import_list(1)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(list)),
		Len:  1, Cap: 1,
	}

	golist := *(*[]C.Sass_Import_Entry)(unsafe.Pointer(&hdr))

	if importResolver != nil {
		newURL, body, resolved := importResolver(rel, parent)
		if resolved {
			// Passing a nil as body is a signal to load the import from the URL.
			var bodyv *C.char
			if body != "" {
				bodyv = C.CString(body)
			}
			ent := C.sass_make_import_entry(C.CString(newURL), bodyv, nil)
			cent := (C.Sass_Import_Entry)(ent)
			golist[0] = cent

			return list
		}
	}
	// TODO1
	if strings.HasPrefix(rel, "compass") {
		ent := C.sass_make_import_entry(url, C.CString(""), nil)
		cent := (C.Sass_Import_Entry)(ent)
		golist[0] = cent
	} else {
		ent := C.sass_make_import_entry(url, nil, nil)
		cent := (C.Sass_Import_Entry)(ent)
		golist[0] = cent
	}

	return list
}
