gophia
======

Gophia is a Go (golang) binding to the Sophia key-value database (http://sphia.org/)

Installation
============

Before installing Gophia, you need to install Sophia. Download it from http://sphia.org, and build.

The makefiles don't include an 'install', so you will need to manually install somewhere where Go can find the headers and the libs. Alternatively, you can set your CGO LDFLAGS and CFLAGS like so:

	# My Sophia files installed to /usr/local/src/sophia
	export CGO_LDFLAGS="-L/usr/local/src/sophia/db"
	export CGO_CFLAGS="-I/usr/local/src/sophia/db"	

Once Sophia is installed, Gophia can be installed with `go get github.com/craigmj/gophia`

Usage
=====

Open the database:

    db, err := gophia.Open(gophia.ReadWrite | gophia.Create, "testdb")
    // check for errors
    defer db.Close()

You're ready to go:

	db.SetString("one","ichi")
	db.SetString("two","nichi")

	fmt.Println("one is ", db.MustGetString("one"))
	fmt.Println("two is ", db.MustGetString("two"))

You can also use a cursor:

	// Without a starting key, every key-value will be returned
	cur, err := db.Cursor(gophia.GTE, nil)
	for cur.Fetch() {
		fmt.Println(cur.KeyString(), "=", cur.ValueString())
	}
	cur.Close()

Of course it's easy to delete a key-value:

	db.DeleteString("one")

See the documentation for more.

Important
=========

When a Cursor is open, no other access to the database is possible: a Cursor locks the entire db, even from other Cursors.

Therefore, you cannot do anything (Set, Delete, etc) while processing a Cursor. Also, you cannot Close the database util the Cursor is itself closed. ***NOT CLOSING THE DATABASE CORRUPTS THE DATABASE***

In Gophia, this is simplified because you can always Close a Cursor (or Database or Environment) even if it has been previously Closed. This means that you can use the form:

    cur, _ := db.Cursor(gophia.GTE, nil)
    defer cur.Close()
    for cur.Fetch() {
    	// ...
    }
    cur.Close()

If for some reason you exit during the loop, your cursor will still Close, and hence the Database as well. If you continue, your Cursor is manually closed.

Gophia also offers the Database.Each method, which iterates over the key-value pairs passing each to a given function. Each takes care of closing the Cursor when it returns.

***MOST IMPORTANTLY*** attempting to change the database while in a Cursor loop will hang the program. ***DO NOT DO THIS:***

    cur, _ := db.Cursor(gophia.GTE, nil)
    defer cur.Close()
    for cur.Fetch() {
    	db.Delete(cur.Key())
    }
    cur.Close()
