// Package gsxgrammar exposes the vendored tree-sitter-gsx parser as a
// tree-sitter Language. The C sources under this directory are copied from
// ../gsxhq/tree-sitter-gsx/src; see README.md to regenerate.
package gsxgrammar

// #cgo CFLAGS: -std=c11 -fPIC -I${SRCDIR}
// #include <stdint.h>
// typedef struct TSLanguage TSLanguage;
// const TSLanguage *tree_sitter_gsx(void);
import "C"

import "unsafe"

// Language returns a pointer to the gsx TSLanguage for use with
// tree_sitter.NewLanguage.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_gsx())
}
