package gophia

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
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
	ReadWrite Access = C.SPO_RDWR
	ReadOnly = C.SPO_RDONLY
	Create = C.SPO_CREAT
)

type Order C.uint32_t
const (
	GreaterThan Order = C.SPGT
	GT = GreaterThan
	GreaterThanEqual = C.SPGTE
	GTE = GreaterThanEqual
	LessThan = C.SPLT
	LT = LessThan
	LessThanEqual = C.SPLTE
	LTE = LessThanEqual
)

// ErrNotFound indicates that the key does not exist in the database.
var ErrNotFound = errors.New("Key not found")
// EnnUnimplemented indicates a method that isn't yet available.
var ErrUnimplemented = errors.New("Not yet implemented")

// Environment is used to configure the database before opening.
type Environment struct {
	unsafe.Pointer
}

// Database is used for accessing a database.
type Database struct {
	unsafe.Pointer
}

// Cursor iterates over key-values in a database.
type Cursor struct {
	unsafe.Pointer
}

// NewEnvironment creates a new environment for opening a database.
// Receivers must call Close() on the returned Environment.
func NewEnvironment() (*Environment, error) {
	env := &Environment{}
	env.Pointer = C.sp_env()
	if nil==env {
		return nil, errors.New("sp_env failed")
	}
	return env, nil
}

// Dir sets the access mode and the directory for the database.
func (env *Environment) Dir(access Access, directory string) error {
	cdir := C.CString(directory)
	defer C.free(unsafe.Pointer(cdir))
	if 0!= C.sp_ctl_dir(env.Pointer, C.uint32_t(access), cdir) {
		return env.Error()
	}
	return nil
}

// Cmp sets the database comparator function to use for 
// ordering keys.
//
// The function must return 0 if the keys are equal, -1 if
// the first key parameter is lower, and 1 if the second key
// parameter is lower.
func (env *Environment) Cmp(cmp Comparator) error {
	if 0!=C.sp_ctl_cmp(env.Pointer, unsafe.Pointer(&cmp)) {
		return env.Error()
	}
	return nil
}

// Page sets the max key count in a single page for the database.
// This option can be tweaked for performance.
func (env *Environment) Page(count int) error {
	if 0!=C.sp_ctl_page(env.Pointer, C.uint32_t(count)) {
		return env.Error()
	}
	return nil
}

// boolToCInt converts a go boolean to a C int value that has
// boolean meaning
func boolToCInt(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}

// GC turns the garbage collector on or off.
func (env *Environment) GC(enabled bool) error {
	if 0!=C.sp_ctl_gc(env.Pointer, boolToCInt(enabled)) {
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
	if 0!=C.sp_ctl_gcf(env.Pointer, C.double(factor)) {
		return env.Error()
	}
	return nil
}

// Grow sets the initial new size and resize factor for new database files.
// The values are used while the database extends during a merge.
//
// This option can be tweaked for performance. 
func (env *Environment) Grow(newsize uint32, newFactor float64) error {
	if 0!=C.sp_ctl_grow(env.Pointer, C.uint32_t(newsize), C.double(newFactor)) {
		return env.Error()
	}
	return nil
}

// Merge sets whether to launch a merger thread during Open().
func (env *Environment) Merge(merge bool) error {
	if 0!=C.sp_ctl_merge(env.Pointer, boolToCInt(merge)) {
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
	if 0!=C.sp_ctl_mergewm(env.Pointer, C.uint32_t(watermark)) {
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

// Open opens the database with the given access permissions in the given directory.
func Open(access Access, directory string) (*Database, error) {
	env, err := NewEnvironment()
	if nil!=err {
		return nil, err
	}
	defer env.Close()

	err = env.Dir(access, directory)
	if nil!=err {
		return nil, err
	}
	return env.Open()
}

// Open() opens the database that has been configured in the Environment.
// At a minimum, it should be necessary to call Dir() on the Environment to
// specify the directory for the database.
func (env *Environment) Open() (*Database,error) {
	db := &Database {}
	db.Pointer = C.sp_open(env.Pointer)
	if nil==db.Pointer {
		return nil, env.Error()
	}
	return db, nil
}

// Close closes the enviroment and frees its associated memory. You must call 
// Close on any Environment created with NewEnvironment.
func (env *Environment) Close() error {
	return sp_close(&env.Pointer)
}

// Close closes the database and frees its associated memory. You must
// call Close on any database opened with Open()
func (db *Database) Close() error {
	return sp_close(&db.Pointer)
}

// Error returns any error on the database. It should not be
// necessary to call this method, since most methods return errors
// automatically.
func (db *Database) Error() error {
	return sp_error(db.Pointer)
}

// Set sets the value of the key.
func (db *Database) Set(key, value []byte) error {
	e := C.sp_set(db.Pointer, unsafe.Pointer(&key[0]), C.size_t(len(key)), unsafe.Pointer(&value[0]), C.size_t(len(value)))
	if 0!=e {
		return db.Error()
	}
	return nil
}

// SetString sets the value of the key. It is a convenience function for working with strings
// rather than byte slices.
func (db *Database) SetString(key, value string) error {
	return db.Set([]byte(key), []byte(value))
}

// SetObject will gob encode the object and store it with the key
func (db *Database) SetObject(key []byte, value interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&buf)
	if nil!=err {
		return err
	}
	return db.Set(key, buf.Bytes())
}

// SetObjectString will gob encode the object and store it with the key. This
// is a convenience method to facilitate working with string keys.
func (db *Database) SetObjectString(key string, value interface{}) error {
	return db.SetObject([]byte(key), value)
}

// Has returns true if the database has a value for the key.
func (db *Database) Has(key []byte) (bool, error) {
	e := C.sp_get(db.Pointer, unsafe.Pointer(&key[0]), C.size_t(len(key)), nil, nil)
	switch int(e) {
	case -1:
		return false, db.Error()
	case 0:
		return false, nil
	case 1:
		return true, nil
	}
	return false, fmt.Errorf("ERROR: unexpected return value from sp_get: %v", e)
}

// HasString returns true if the database has a value for the key. It is a convenience
// function for working with strings rather than byte slices.
func (db *Database) HasString(key string) (bool, error) {
	return db.Has([]byte(key))
}

// Get retrieves the value for the key.
func (db *Database) Get(key []byte) ([]byte, error) {
	// return nil, ErrUnimplemented
	var vptr unsafe.Pointer
	var size C.size_t
	// var vptr unsafe.Pointer

	e := C.sp_get(db.Pointer, unsafe.Pointer(&key[0]), C.size_t(len(key)), &vptr, (*C.size_t)(&size))
	switch int(e) {
	case -1:
		return nil, db.Error()
	case 0:
		return nil, ErrNotFound
	case 1:
		// Continue after the switch
	default:
		return nil, fmt.Errorf("ERROR: unexpected return value from sp_get: %v", e)
	}
	value := C.GoBytes(vptr, C.int(size))
	C.sp_destroy(vptr)
	return value, nil
}

// MustGet returns the value of the key, or panics on an error.
func (db *Database) MustGet(key []byte) []byte {
	val, err := db.Get(key)
	if nil!=err {
		panic(err)
	}
	return val
}

// GetObject retrieves a gob encoded object into the out object.
func (db *Database) GetObject(key []byte, out interface{}) error {
	buf, err := db.Get(key)
	if nil!=err {
		return err
	}
	dec := gob.NewDecoder(bytes.NewReader(buf))
	return dec.Decode(out)
}

// GetObjectString retrieves a gob encoded object into the out object. It is a 
// convenience method to facilitate working with string keys.
func (db *Database) GetObjectString(key string, out interface{}) error {
	return db.GetObject([]byte(key), out)
}

// GetString retrieves the string value for the string key. It is a convenience function
// for working with strings rather than byte slices.
func (db *Database) GetString(key string) (string, error) {
	v, err := db.Get([]byte(key))
	if nil!=err {
		return "", err
	}
	return string(v), nil
}

// MustGetString returns the stirng value for a string key. It panics
// on an error.
func (db *Database) MustGetString(key string) string {
	value, err := db.Get([]byte(key))
	if nil!=err {
		panic(err)
	}
	return string(value)
}

// Delete deletes the key from the database.
func (db *Database) Delete(key []byte) error {
	if 0 != C.sp_delete(db.Pointer, unsafe.Pointer(&key[0]), C.size_t(len(key))) {
		return db.Error()
	}
	return nil
}

// DeleteString deletes the key from the database.
func (db *Database) DeleteString(key string) error {
	return db.Delete([]byte(key))
}

// Cursor returns a Cursor for iterating over rows in the database.
//
// If no key is provided, the Cursor will iterate over all rows.
//
// The order flag decides the direction of the iteration, and whether
// the key is included or excluded.
//
// Iterate over values with Fetch or Next methods.
func (db *Database) Cursor(order Order, key []byte) (*Cursor, error) {
	cur := &Cursor{}
	if nil==key {
		cur.Pointer = C.sp_cursor(db.Pointer, C.sporder(order), unsafe.Pointer(nil), C.size_t(0))
	} else {
		cur.Pointer = C.sp_cursor(db.Pointer, C.sporder(order), unsafe.Pointer(&key[0]), C.size_t(len(key)))
	}
	if nil==cur.Pointer {
		return nil, db.Error()
	}
	return cur, nil
}

// Each iterates through the key-values in the database, passing each to the each function.
// It is a convenience wrapper around a Cursor iteration.
func (db *Database) Each(order Order, key []byte, each func(key []byte, value []byte)) error {
	cur, err := db.Cursor(order, key)
	defer cur.Close()
	if nil!=err {
		return err
	}
	for cur.Fetch() {
		each(cur.Key(), cur.Value())
	}
	return nil
}

// CursorString returns a Cursor that fetches rows from the database
// from the given key, passed as a string.
// Callers must call Close() on the received Cursor.
func (db *Database) CursorString(order Order, key string) (*Cursor, error) {
	return db.Cursor(order, []byte(key))
}

// Close closes the cursor. If a cursor is not closed, future operations
// on the database can hang indefinitely.
func (cur *Cursor) Close() error {
	return sp_close(&cur.Pointer)
}

// Fetch fetches the next row for the cursor, and returns
// true if there is a next row, false if the cursor has reached the
// end of the rows.
func (cur *Cursor) Fetch() bool {
	return C.int(1)==C.sp_fetch(cur.Pointer)
}

// Next is identical to Fetch. It exists because it
// seems that Next() is more go-idiomatic.
func (cur *Cursor) Next() bool {
	return C.int(1)==C.sp_fetch(cur.Pointer)
}

// KeySize returns the size of the current key.
func (cur *Cursor) KeySize() int {
	return int(C.sp_keysize(cur.Pointer))
}

// KeyLen returns the length of the current key. It is 
// a synonym for KeySize()
func (cur *Cursor) KeyLen() int {
	return cur.KeySize()
}

// ValueSize returns the length of the current value.
func (cur *Cursor) ValueSize() int {
	return int(C.sp_valuesize(cur.Pointer))
}

// ValueLen returns the length of the current value. It is
// a synonym for ValueSize()
func (cur *Cursor) ValueLen() int {
	return cur.ValueSize()
}

// Key returns the current key of the cursor.
func (cur *Cursor) Key() []byte {
	size := C.int(C.sp_keysize(cur.Pointer))
	if 0==size {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(C.sp_key(cur.Pointer)), size)
}

// KeyString returns the current key as a string.
func (cur *Cursor) KeyString() string {
	return string(cur.Key())
}

// Value returns the current value of the cursor.
func (cur *Cursor) Value() []byte {
	size := C.int(C.sp_valuesize(cur.Pointer))
	if 0==size {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(C.sp_value(cur.Pointer)), size)
}

// ValueString returns the current value as a string.
func (cur *Cursor) ValueString() string {
	return string(cur.Value())
}

// ValueObject returns the current object, by gob decoding the
// current value at the cursor.
func (cur *Cursor) Object(out interface{}) error {
	buf := cur.Value()
	if nil==buf {
		return errors.New("Value is nil")
	}
	dec := gob.NewDecoder(bytes.NewReader(buf))
	return dec.Decode(out)
}


//export go_sp_comparator
func go_sp_comparator(aptr unsafe.Pointer, asz C.size_t, bptr unsafe.Pointer, bsz C.size_t, arg unsafe.Pointer) C.int {
	a := C.GoBytes(aptr, C.int(asz))
	b := C.GoBytes(bptr, C.int(bsz))
	cmp := (*Comparator)(arg)
	return C.int((*cmp)(a,b))
}

// sp_close closes the pointer and sets it to nil
// to ensure it cannot be closed twice.
func sp_close(p *unsafe.Pointer) error {
	if nil==*p {
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
	if nil==cerror {
		return nil
	}
	return errors.New(C.GoString(cerror))
}