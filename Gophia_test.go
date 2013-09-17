package gophia

import (
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

	err = db.SetStrings("one", "ichi")
	checkErr(err)
	err = db.SetStrings("two", "nichi")
	checkErr(err)
	one, err := db.GetStrings("one")
	checkErr(err)
	if one != "ichi" {
		t.Fatal("one not saved / restored right")
	}
	two, err := db.GetStrings("two")
	checkErr(err)
	if two != "nichi" {
		t.Fatal("two not saved / restored right")
	}
	if !db.MustHasString("one") {
		t.Fatal("HasString failed for key in database")
	}
	checkErr(db.DeleteString("one"))
	if db.MustHasString("one") {
		t.Fatal("HasString succeeded for key in database")
	}
	_, err = db.GetStrings("one")
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
