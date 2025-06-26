package json

import (
	"fmt"
	"testing"
	"time"
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

func TestCopyStructPanic(t *testing.T) {
	s := Student{
		Name:   "acexy",
		Sex:    1,
		School: "Q",
	}

	var person Person
	CopyStructPanic(s, &person)
	fmt.Printf("%+v\n", person)

	ss := []*Student{{Name: "acexy", Sex: 1, School: "Q"}, {Name: "acexy", Sex: 1, School: "H"}}

	fmt.Println(ToJson(ss))
	fmt.Println(ToJsonFormat(ss))
}

type User struct {
	Name string     `json:"name"`
	Time *Timestamp `json:"time"`
}

type Model[T any] struct {
	T    T `json:"data"`
	Code int
}

func TestTimestamp(t *testing.T) {
	user := User{
		Name: "acexy",
		Time: &Timestamp{time.Now()},
	}
	//GlobalWrapperSetting(func(options *TypeWrapperOptions) {
	//	options.TimestampType = TimestampTypeSecond
	//})
	fmt.Println(ToJson(user))
	ParseJson("{\"name\":\"acexy\",\"time\":1729136314000}", &user)
	fmt.Println(user)

	model := Model[User]{
		T:    user,
		Code: 200,
	}
	json := ToJson(model)
	fmt.Println(json)

	var m1 Model[User]
	ParseJson(json, &m1)
	fmt.Println(m1)
}
