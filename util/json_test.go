package util

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
	Sex  uint
}

type Student struct {
	Name   string
	Sex    uint
	School string
}

func TestJsonCopyStructPanic(t *testing.T) {
	s := Student{
		Name:   "acexy",
		Sex:    1,
		School: "Q",
	}

	var person Person
	JsonCopyStructPanic(s, &person)

	fmt.Printf("%+v\n", person)

}
