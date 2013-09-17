package gophia

import (
	"bytes"
	"encoding/gob"
	"errors"
)

// CursorS returns a Cursor that fetches rows from the database
// from the given key, passed as a string.
// Callers must call Close() on the received Cursor.
func (db *Database) CursorS(order Order, key string) (*Cursor, error) {
	return db.Cursor(order, []byte(key))
}

// DeleteS deletes the key from the database.
func (db *Database) DeleteS(key string) error {
	return db.Delete([]byte(key))
}

// Each iterates through the key-values in the database, passing each to the each function.
// It is a convenience wrapper around a Cursor iteration.
func (db *Database) Each(order Order, key []byte, each func(key []byte, value []byte)) error {
	cur, err := db.Cursor(order, key)
	defer cur.Close()
	if nil != err {
		return err
	}
	for cur.Fetch() {
		each(cur.Key(), cur.Value())
	}
	return nil
}

// GetAO returns on object value for a byte-array key.
func (db *Database) GetAO(key []byte, out interface{}) error {
	buf, err := db.Get(key)
	if nil != err {
		return err
	}
	dec := gob.NewDecoder(bytes.NewReader(buf))
	return dec.Decode(out)
}

// GetS retrieves an array value for a string key. It is a convenience
// method to simplify working with string keys.
func (db *Database) GetSA(key string) ([]byte, error) {
	return db.Get([]byte(key))
}

// GetSO fetches a gob encoded object for a string key.
func (db *Database) GetSO(key string, out interface{}) error {
	return db.GetAO([]byte(key), out)
}

// GetSS returns a string value for a string key.
func (db *Database) GetSS(key string) (string, error) {
	v, err := db.Get([]byte(key))
	if nil != err {
		return "", err
	}
	return string(v), nil
}

// HasS returns true if the database has a value for the string key.
func (db *Database) HasS(key string) (bool, error) {
	return db.Has([]byte(key))
}

// KeyLen returns the length of the current key. It is
// a synonym for KeySize()
func (cur *Cursor) KeyLen() int {
	return cur.KeySize()
}

// KeyS returns the current key as a string.
func (cur *Cursor) KeyS() string {
	return string(cur.Key())
}

// MustHas returns true if the key exists, false otherwise. It panics
// in the even of error.
func (db *Database) MustHas(key []byte) bool {
	has, err := db.Has(key)
	if nil != err {
		panic(err)
	}
	return has
}

// MustGet returns the value of the key, or panics on an error.
func (db *Database) MustGet(key []byte) []byte {
	val, err := db.Get(key)
	if nil != err {
		panic(err)
	}
	return val
}

// MustGetSA returns the byte array value for the string key. It panics on
// error.
func (db *Database) MustGetSA(key string) []byte {
	value, err := db.Get([]byte(key))
	if nil != err {
		panic(err)
	}
	return value
}

// MustGetString returns the byte array value for the string key. It panics
// on error.
// @deprecated Use MustGetSA instead.
func (db *Database) MustGetString(key string) []byte {
	return db.MustGetSA(key)
}

// MustGetSS returns the string value for a string key. It panics
// on an error.
func (db *Database) MustGetSS(key string) string {
	value, err := db.Get([]byte(key))
	if nil != err {
		panic(err)
	}
	return string(value)
}

// MustHasS returns true if the string exists, or false if it does not. It panics
// in the event of error.
func (db *Database) MustHasS(key string) bool {
	return db.MustHas([]byte(key))
}

// Next is identical to Fetch. It exists because it
// seems that Next() is more go-idiomatic.
func (cur *Cursor) Next() bool {
	return cur.Fetch()
}

// Open opens the database with the given access permissions in the given directory.
func Open(access Access, directory string) (*Database, error) {
	env, err := NewEnvironment()
	if nil != err {
		return nil, err
	}

	err = env.Dir(access, directory)
	if nil != err {
		env.Close()
		return nil, err
	}
	db, err := env.Open()
	if nil != err {
		env.Close()
		return nil, err
	}
	db.env = env
	return db, nil
}

// SetAO sets a byte array key to an object value.
func (db *Database) SetAO(key []byte, value interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if nil != err {
		return err
	}
	return db.Set(key, buf.Bytes())
}

// SetSA sets a string key to a byte array value.
func (db *Database) SetSA(key string, value []byte) error {
	return db.Set([]byte(key), value)
}

// SetSS sets a string key to a string value
func (db *Database) SetSS(key, value string) error {
	return db.Set([]byte(key), []byte(value))
}

// SetSO sets a string key to an object value.
func (db *Database) SetSO(key string, value interface{}) error {
	return db.SetAO([]byte(key), value)
}

// ValueLen returns the length of the current value. It is
// a synonym for ValueSize()
func (cur *Cursor) ValueLen() int {
	return cur.ValueSize()
}

// ValueO returns the current value as an object, by gob decoding
// the current value at the cursor.
func (cur *Cursor) ValueO(out interface{}) error {
	buf := cur.Value()
	if nil == buf {
		return errors.New("Value is nil")
	}
	dec := gob.NewDecoder(bytes.NewReader(buf))
	return dec.Decode(out)
}

// ValueS returns the current value as a string.
func (cur *Cursor) ValueS() string {
	return string(cur.Value())
}
