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

Create an Environment:

    env, err := gophia.NewEnvironment()
    // check for errors
    defer env.Close()

Set the database access type and directory:

	err = env.Dir(gophia.ReadWrite | gophia.Create, "testdb")
	// check for errors

Open the database:

    db, err := env.Open()
    // check for errors
    defer db.Close()

All of the above steps can be wrapped into:

    db, err := gophia.Open(gophia.ReadWrite | gophia.Create, "testdb")
    // check for errors
    defer db.Close()

CAUTION: When using the shorthand gophia.Open(), the database's environment is destroyed after the
database is opened. At this point, I'm unsure whether this might make Sophia unstable. If it does,
we'll do a sensible workaround. For now, use this form with a little caution.

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


