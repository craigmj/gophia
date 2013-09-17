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

type Access C.uint32_t

// Comparator function is used to compare keys in the database.
//
// The function must return 0 if the keys are equal, -1 if
// the first key parameter is lower, and 1 if the second key
// parameter is lower.
//
// See Environment.Cmp()
type Comparator func(a []byte, b []byte) int

const (
	ComparesEqual       int = 0
	ComparesLessThan        = -1
	ComparesGreaterThan     = 1
)

const (
	ReadWrite Access = C.SPO_RDWR
	ReadOnly         = C.SPO_RDONLY
	Create           = C.SPO_CREAT
)

// Environment is used to configure the database before opening.
type Environment struct {
	unsafe.Pointer
}

// NewEnvironment creates a new environment for opening a database.
// Receivers must call Close() on the returned Environment.
func NewEnvironment() (*Environment, error) {
	env := &Environment{}
	env.Pointer = C.sp_env()
	if nil == env {
		return nil, errors.New("sp_env failed")
	}
	return env, nil
}

// Dir sets the access mode and the directory for the database.
func (env *Environment) Dir(access Access, directory string) error {
	cdir := C.CString(directory)
	defer C.free(unsafe.Pointer(cdir))
	if 0 != C.sp_ctl_dir(env.Pointer, C.uint32_t(access), cdir) {
		return env.Error()
	}
	return nil
}

// Close closes the enviroment and frees its associated memory. You must call
// Close on any Environment created with NewEnvironment.
func (env *Environment) Close() error {
	return sp_close(&env.Pointer)
}

// Cmp sets the database comparator function to use for
// ordering keys.
//
// The function must return 0 if the keys are equal, -1 if
// the first key parameter is lower, and 1 if the second key
// parameter is lower.
func (env *Environment) Cmp(cmp Comparator) error {
	if 0 != C.sp_ctl_cmp(env.Pointer, unsafe.Pointer(&cmp)) {
		return env.Error()
	}
	return nil
}

// Error returns any error on the Environment. It should not be
// necessary to call this method, since the Go methods all return
// with errors themselves.
func (env *Environment) Error() error {
	return sp_error(env.Pointer)
}

// Page sets the max key count in a single page for the database.
// This option can be tweaked for performance.
func (env *Environment) Page(count int) error {
	if 0 != C.sp_ctl_page(env.Pointer, C.uint32_t(count)) {
		return env.Error()
	}
	return nil
}

// GC turns the garbage collector on or off.
func (env *Environment) GC(enabled bool) error {
	if 0 != C.sp_ctl_gc(env.Pointer, boolToCInt(enabled)) {
		return env.Error()
	}
	return nil
}

// GCF sets database garbage collector factor value, which is
// used to determine when to start the GC.
//
// For example: factor 0.5 means that all 'live' pages from any db
// file will be copied to new db when half or fewer of them are left.
//
// This option can be tweaked for performance.
func (env *Environment) GCF(factor float64) error {
	if 0 != C.sp_ctl_gcf(env.Pointer, C.double(factor)) {
		return env.Error()
	}
	return nil
}

// Grow sets the initial new size and resize factor for new database files.
// The values are used while the database extends during a merge.
//
// This option can be tweaked for performance.
func (env *Environment) Grow(newsize uint32, newFactor float64) error {
	if 0 != C.sp_ctl_grow(env.Pointer, C.uint32_t(newsize), C.double(newFactor)) {
		return env.Error()
	}
	return nil
}

// Merge sets whether to launch a merger thread during Open().
func (env *Environment) Merge(merge bool) error {
	if 0 != C.sp_ctl_merge(env.Pointer, boolToCInt(merge)) {
		return env.Error()
	}
	return nil
}

// MergeWM sets the database merge watermark value.
//
// When the database update count reaches this value, it notifies
// the merger thread to create a new epoch and start merging
// in-memory keys.
//
// This option can be tweaked for performance.
func (env *Environment) MergeWM(watermark uint32) error {
	if 0 != C.sp_ctl_mergewm(env.Pointer, C.uint32_t(watermark)) {
		return env.Error()
	}
	return nil
}

// Open() opens the database that has been configured in the Environment.
// At a minimum, it should be necessary to call Dir() on the Environment to
// specify the directory for the database.
func (env *Environment) Open() (*Database, error) {
	db := &Database{}
	db.Pointer = C.sp_open(env.Pointer)
	if nil == db.Pointer {
		return nil, env.Error()
	}
	return db, nil
}

// boolToCInt converts a go boolean to a C int value that has
// boolean meaning
func boolToCInt(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}
