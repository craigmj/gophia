package gophia

/* DEPRECATED METHODS *************************************************************/

// CursorString returns a Cursor that fetches rows from the database
// from the given key, passed as a string.
// Callers must call Close() on the received Cursor.
// @deprecated Use CursorS instead.
func (db *Database) CursorString(order Order, key string) (*Cursor, error) {
	return db.Cursor(order, []byte(key))
}

// DeleteString deletes the key from the database.
// @deprecated Use DeleteS instead.
func (db *Database) DeleteString(key string) error {
	return db.DeleteS(key)
}

// GetObject fetches a gob encoded object into the out object.
// @deprecated Use GetAO instead.
func (db *Database) GetObject(key []byte, out interface{}) error {
	return db.GetAO(key, out)
}

// GetObjectString retrieves a gob encoded object into the out object. It is a
// convenience method to facilitate working with string keys.
// @deprecated Use GetSO instead.
func (db *Database) GetObjectString(key string, out interface{}) error {
	return db.GetSO(key, out)
}

// GetString returns a byte array value for a string key.
// @deprecated Use GetS instead.
func (db *Database) GetString(key string) ([]byte, error) {
	return db.GetSA(key)
}

// GetString retrieves the string value for the string key. It is a convenience function
// for working with strings rather than byte slices.
// @deprecated Use GetSS instead.
func (db *Database) GetStrings(key string) (string, error) {
	return db.GetSS(key)
}

// HasString returns true if the database has a value for the key. It is a convenience
// function for working with strings rather than byte slices.
// @deprecated Use HasS instead
func (db *Database) HasString(key string) (bool, error) {
	return db.Has([]byte(key))
}

// KeyString returns the current key as a string.
// @deprecated Use KeyS instead
func (cur *Cursor) KeyString() string {
	return string(cur.Key())
}

// MustGetStrings returns the string value for a string key. It panics on error.
// @deprecated Use MustGetSS instead.
func (db *Database) MustGetStrings(key string) string {
	return db.MustGetSS(key)
}

// MustHasString returns true if the string exists, or false if it does not. It panics
// in the event of error.
// @deprecated Use MustHasS instead
func (db *Database) MustHasString(key string) bool {
	return db.MustHas([]byte(key))
}

// SetObject will gob encode the object and store it with the key.
// @deprecated Use SetAO instead
func (db *Database) SetObject(key []byte, value interface{}) error {
	return db.SetAO(key, value)
}

// SetObjectString will gob encode the object and store it with the key. This
// is a convenience method to facilitate working with string keys.
// @deprecated Prefer SetSO
func (db *Database) SetObjectString(key string, value interface{}) error {
	return db.SetSO(key, value)
}

// SetString sets the byte-slice value of the string key. It is a convenience function for working with string keys
// rather than byte slices.
// @deprecated Use SetSA instead
func (db *Database) SetString(key string, value []byte) error {
	return db.SetSA(key, value)
}

// SetStrings sets the value of the key. It is a convenience function for working with strings
// rather than byte slices.
// @deprecated Use SetSS instead
func (db *Database) SetStrings(key, value string) error {
	return db.SetSS(key, value)
}

// ValueObject returns the current object, by gob decoding the
// current value at the cursor.
// @deprecated Use ValueO instead.
func (cur *Cursor) Object(out interface{}) error {
	return cur.ValueO(out)
}

// ValueString returns the current value as a string.
// @deprecated Use ValueS instead
func (cur *Cursor) ValueString() string {
	return string(cur.Value())
}
