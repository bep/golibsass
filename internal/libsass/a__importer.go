// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
package libsass

// #include <stdint.h>
// #include "sass/context.h"
//
// extern struct Sass_Import** BridgeImport(const char* currPath, const char* prevPath, int i);
//
// Sass_Import_List SassImport(const char* currPath, Sass_Importer_Entry imp, struct Sass_Compiler* comp)
// {
//   void* c = sass_importer_get_cookie(imp);
//   uintptr_t ci = (uintptr_t)c;
//   struct Sass_Import* prevPath = sass_compiler_get_last_import(comp);
//   const char* prev_path = sass_import_get_imp_path(prevPath);
//   return BridgeImport(currPath, prev_path, ci);
// }
import "C"
import (
	"sync"
	"unsafe"
)

const uintptrOffset = 4096 // the smallest legal pointer

var importsStore = &idMap{
	m: make(map[int]interface{}),
	i: uintptrOffset,
}

// AddImportResolver adds a function to resolve imports in LibSASS.
// Make sure to run call DeleteImportResolver when done.
//go:nocheckptr
func AddImportResolver(opts SassOptions, resolver ImportResolver) int {
	i := importsStore.Set(resolver)
	// This looks unsafe, but LibSass is using void* to store an int.
	// TODO(bep) this prevents us from "fail on go vet errors" in Travis.
	id := unsafe.Pointer(uintptr(i))

	importers := C.sass_make_importer_list(1)
	C.sass_importer_set_list_entry(
		importers,
		0,
		C.sass_make_importer(
			C.Sass_Importer_Fn(C.SassImport),
			C.double(0),
			id,
		),
	)

	C.sass_option_set_c_importers(
		(*C.struct_Sass_Options)(unsafe.Pointer(opts)),
		importers,
	)

	return i
}

func DeleteImportResolver(i int) error {
	importsStore.Delete(i)
	return nil
}

// ImportResolver can be used as a custom import resolver.
// Return an empty body to load the import body from the path.
// See AddImportResolver.
type ImportResolver func(currPath string, prevPath string) (newPath string, body string, resolved bool)

type idMap struct {
	sync.RWMutex
	m       map[int]interface{}
	i       int
	idStack []int
}

func (m *idMap) Delete(i int) {
	m.Lock()
	defer m.Unlock()
	m.idStack = append(m.idStack, i)
	delete(m.m, i)
}

func (m *idMap) Get(i int) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.m[i]
}

func (m *idMap) nextID() (id int) {
	if len(m.idStack) == 0 {
		for i := 1; i <= 50; i++ {
			m.idStack = append(m.idStack, m.i+i)
		}
		m.i += 50
	}
	id, m.idStack = m.idStack[len(m.idStack)-1], m.idStack[:len(m.idStack)-1]
	return
}

func (m *idMap) Set(v interface{}) int {
	m.Lock()
	defer m.Unlock()
	id := m.nextID()
	m.m[id] = v
	return id
}
