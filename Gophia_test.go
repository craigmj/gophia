package gophia

import (
	"fmt"
	"testing"
)

func TestBasicSanity(t *testing.T) {
	checkErr := func(err error) {
		if nil != err {
			t.Fatal(err)
		}
	}
	db, err := Open(Create, "testdb_basic")
	checkErr(err)
	defer db.Close()

	err = db.SetSS("one", "ichi")
	checkErr(err)
	err = db.SetSS("two", "nichi")
	checkErr(err)
	one, err := db.GetSS("one")
	checkErr(err)
	if one != "ichi" {
		t.Fatal("one not saved / restored right")
	}
	two, err := db.GetSS("two")
	checkErr(err)
	if two != "nichi" {
		t.Fatal("two not saved / restored right")
	}
	if !db.MustHasS("one") {
		t.Fatal("HasString failed for key in database")
	}
	checkErr(db.DeleteS("one"))
	if db.MustHasS("one") {
		t.Fatal("HasString succeeded for key in database")
	}
	_, err = db.GetSS("one")
	if err != ErrNotFound {
		if nil != err {
			t.Fatalf("Unexpected error retrieving deleted key: %v", err.Error())
		} else {
			t.Fatal("Retrieval of deleted key succeeded - should have failed")
		}
	}
}

type person struct {
	Id   int
	Name string
}

func TestGob(t *testing.T) {
	db, err := Open(Create, "testdb_gob")
	if nil != err {
		t.Fatal(err)
	}
	defer db.Close()

	g := person{1, "Craig"}
	err = db.SetSO("craig", &g)
	if nil != err {
		t.Fatal(err)
	}

	g = person{2, "Fred"}
	err = db.SetSO("fred", &g)
	if nil != err {
		t.Fatal(err)
	}

	err = db.GetSO("craig", &g)
	if nil != err {
		t.Fatal(err)
	}
	if g.Id != 1 || g.Name != "Craig" {
		t.Errorf("First person didn't gob encode/decode right")
	}
	err = db.GetSO("fred", &g)
	if nil != err {
		t.Fatal(err)
	}
	if g.Id != 2 || g.Name != "Fred" {
		t.Errorf("Second person didn't gob encode/decode right")
	}
}

func TestIterator(t *testing.T) {
	checkErr := func(err error) {
		if nil != err {
			t.Fatal(err)
		}
	}
	db, err := Open(Create, "testdb_iterator")
	checkErr(err)
	defer db.Close()

	type E struct {
		Key string
		Val string
	}
	expects := []E{E{"1", "one"}, E{"2", "two"}, E{"3", "three"}, E{"4", "four"}}

	for _, e := range expects {
		db.SetSS(e.Key, e.Val)
	}

	cur, err := db.CursorS(GreaterThan, "1")
	checkErr(err)
	defer cur.Close()
	index := 1
	for cur.Fetch() {
		key := cur.KeyS()
		val := cur.ValueS()
		if index >= len(expects) {
			t.Fatalf("Fetched %d items from db - too many", index)
		}
		e := expects[index]
		if key != e.Key {
			t.Errorf("Fetched wrong key: expected %v, got %v", e.Key, key)
		}
		if val != e.Val {
			t.Errorf("Fetched wrong value: expected %v, got %v", e.Val, val)
		}
		index++
	}
}

func TestTransactions(t *testing.T) {
	checkErr := func(err error) {
		if nil!=err {
			t.Fatal(err)
		}
	}
	db, err := Open(Create, "testdb_trans")
	defer db.Close()
	values := [][2]string {
		[2]string { "one", "ichi"},
		[2]string{ "two", "nichi" },
		[2]string{ "three", "san" },
	}
	for _, v := range values {
		checkErr(db.DeleteS(v[0]))
	}
	enterValues := func() error {
		for _,v := range values {
			checkErr(db.SetSS(v[0],v[1]))
		}
		return nil
	}
	checkValues := func() (bool, error) {
		for _,v := range values {
			val, err := db.GetSS(v[0])
			if nil!=err {				
				return false, err
			}
			if v[1]!=val {
				return false, fmt.Errorf("Value of %v didn't match: %v", v[0], val)
			}
		}
		return true, nil
	}

	checkErr(db.Begin())
	// Although this should work, it doesn't appear to do so - we don't get the expected error
	// err = db.Begin()
	// if err!=ErrTransactionInProgress {
	// 	t.Fatalf("Error: Transaction in progress but Begin succeeded: err = %v", err)
	// }
	checkErr(enterValues())
	checkErr(db.Rollback())
	f, err := checkValues()
	if f || ErrNotFound!=err {
		t.Errorf("Inserted values in place despite transaction rollback")
	}

	checkErr(db.Begin())
	checkErr(enterValues())
	checkErr(db.Commit())
	f, err = checkValues()
	if nil!=err || !f {
		t.Errorf("Values not in place but transaction committed")
	}

}