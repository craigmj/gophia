package gophia

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -lsophia
#include <sophia.h>

*/
import "C"

// Cursor iterates over key-values in a database.
type Cursor struct {
	unsafe.Pointer
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
	return C.int(1) == C.sp_fetch(cur.Pointer)
}

// Key returns the current key of the cursor.
func (cur *Cursor) Key() []byte {
	size := C.int(C.sp_keysize(cur.Pointer))
	if 0 == size {
		fmt.Println("Key is 0 len")
		return nil
	}
	return C.GoBytes(unsafe.Pointer(C.sp_key(cur.Pointer)), size)
}

// KeySize returns the size of the current key.
func (cur *Cursor) KeySize() int {
	return int(C.sp_keysize(cur.Pointer))
}

// Value returns the current value of the cursor.
func (cur *Cursor) Value() []byte {
	size := C.int(C.sp_valuesize(cur.Pointer))
	if 0 == size {
		fmt.Println("Value is 0 len")
		return nil
	}
	return C.GoBytes(unsafe.Pointer(C.sp_value(cur.Pointer)), size)
}

// ValueSize returns the length of the current value.
func (cur *Cursor) ValueSize() int {
	return int(C.sp_valuesize(cur.Pointer))
}
