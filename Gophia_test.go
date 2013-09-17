package gophia

import (
	"fmt"
	"testing"
)

func TestBasicSanity(t *testing.T) {
	// TODO: Testing to be added	
}

type person struct {
	Id int
	Name string
}

func TestGob(t *testing.T) {
	db, err := Open(ReadWrite | Create, "testdb")
	if nil!=err {
		t.Fatal(err)
	}
	defer db.Close()

	g := person{1, "Craig"}
	err = db.SetObjectString("craig", &g)
	if nil!=err {
		t.Fatal(err)
	}

	g = person{2, "Fred"}
	err = db.SetObjectString("fred", &g)
	if nil!=err {
		t.Fatal(err)
	}
	
	err = db.GetObjectString("craig", &g)
	if nil!=err {
		t.Fatal(err)
	}
	fmt.Println(g)
	err = db.GetObjectString("fred", &g)
	if nil!=err {
		t.Fatal(err)
	}
	fmt.Println(g)

}
