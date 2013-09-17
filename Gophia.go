package gophia

import (
	"errors"
	"unsafe"
)

/*
#cgo LDFLAGS: -lsophia
#include <sophia.h>

extern int sp_ctl_dir(void *p, uint32_t access, char *dir);
extern int sp_ctl_cmp(void *p, void *cmp);
extern int sp_ctl_page(void *p, uint32_t count);
extern int sp_ctl_gc(void *p, int active);
extern int sp_ctl_gcf(void *p, double factor);
extern int sp_ctl_grow(void *p, uint32_t newsize, double factor);
extern int sp_ctl_merge(void *p, int merge);
extern int sp_ctl_mergewm(void *p, uint32_t watermark);

*/
import "C"

//export go_sp_comparator
func go_sp_comparator(aptr unsafe.Pointer, asz C.size_t, bptr unsafe.Pointer, bsz C.size_t, arg unsafe.Pointer) C.int {
	a := C.GoBytes(aptr, C.int(asz))
	b := C.GoBytes(bptr, C.int(bsz))
	cmp := (*Comparator)(arg)
	return C.int((*cmp)(a, b))
}

// sp_close closes the pointer and sets it to nil
// to ensure it cannot be closed twice.
func sp_close(p *unsafe.Pointer) error {
	if nil == *p {
		return nil
	}
	if 0 != C.sp_destroy(*p) {
		return sp_error(*p)
	}
	*p = nil
	return nil
}

// sp_error returns the error for the given
// Sophia pointer as a golang error
func sp_error(p unsafe.Pointer) error {
	cerror := C.sp_error(p)
	if nil == cerror {
		return nil
	}
	return errors.New(C.GoString(cerror))
}
