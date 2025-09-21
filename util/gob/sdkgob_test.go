package gob

import (
	"fmt"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func TestEncode(t *testing.T) {

	user := User{
		Name: "zhangsan",
		Age:  18,
	}

	fmt.Println(Encode(user))
	fmt.Println(Encode(&user))

}
