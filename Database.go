package gophia

import (
	"errors"
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -lsophia
#include <sophia.h>

*/
import "C"

// Order of items in an iteration
type Order C.uint32_t

const (
	GreaterThan      Order = C.SPGT
	GT                     = GreaterThan
	GreaterThanEqual       = C.SPGTE
	GTE                    = GreaterThanEqual
	LessThan               = C.SPLT
	LT                     = LessThan
	LessThanEqual          = C.SPLTE
	LTE                    = LessThanEqual
)

// ErrNotFound indicates that the key does not exist in the database.
var ErrNotFound = errors.New("Key not found")
// ErrTransactionInProgress returned when attempt to begin a transaction while there is already
// a transaction in progress.
var ErrTransactionInProgress = errors.New("Transaction already in progress")

// Database is used for accessing a database.
type Database struct {
	unsafe.Pointer
	env *Environment
}

// Begin starts a multi-statement transaction.
//
// During transaction, all updates are not written to the database files
// until a sp_commit(3) is called. All updates made during transaction are
// available through sp_get(3) or by using cursor.

// No nested transactions are supported.
func (db *Database) Begin() error {
	switch C.sp_begin(db.Pointer) {
	case 0:
		return nil
	case 1:
		return ErrTransactionInProgress
	case -1:
		return db.Error()
	}
	// All cases should be handled by switch
	return errors.New("Unexpected return from sp_begin, not in [-1,1]")
}

// Close closes the database and frees its associated memory. You must
// call Close on any database opened with Open()
func (db *Database) Close() error {
	err := sp_close(&db.Pointer)
	if nil != err {
		return err
	}
	if nil != db.env {
		return db.env.Close()
	}
	return nil
}

// Commit applies changes to a multi-statement
// transaction. All modifications made during the transaction are written to 
// the log file in a single batch.
//
// If commit failed, transaction modifications are discarded.
func (db *Database) Commit() error {
	e := C.sp_commit(db.Pointer)
	if 0!=e {
		return db.Error()
	}
	return nil
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
	if nil == key {
		cur.Pointer = C.sp_cursor(db.Pointer, C.sporder(order), unsafe.Pointer(nil), C.size_t(0))
	} else {
		cur.Pointer = C.sp_cursor(db.Pointer, C.sporder(order), unsafe.Pointer(&key[0]), C.size_t(len(key)))
	}
	if nil == cur.Pointer {
		return nil, db.Error()
	}
	return cur, nil
}

// Delete deletes the key from the database.
func (db *Database) Delete(key []byte) error {
	if 0 != C.sp_delete(db.Pointer, unsafe.Pointer(&key[0]), C.size_t(len(key))) {
		return db.Error()
	}
	return nil
}

// Error returns any error on the database. It should not be
// necessary to call this method, since most methods return errors
// automatically.
func (db *Database) Error() error {
	return sp_error(db.Pointer)
}

// Get retrieves the value for the key.
func (db *Database) Get(key []byte) ([]byte, error) {
	var vptr unsafe.Pointer
	var size C.size_t

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

// Rollback discards the changes of a multi-statement
// transaction. All modifications made during the transaction are not written to 
// the log file.
func (db *Database) Rollback() error {
	e := C.sp_rollback(db.Pointer)
	if 0!=e {
		return db.Error()
	}
	return nil
}

// Set sets the value of the key.
func (db *Database) Set(key, value []byte) error {
	e := C.sp_set(db.Pointer, unsafe.Pointer(&key[0]), C.size_t(len(key)), unsafe.Pointer(&value[0]), C.size_t(len(value)))
	if 0 != e {
		return db.Error()
	}
	return nil
}
