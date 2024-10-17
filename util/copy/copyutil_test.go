package copy

import (
	"fmt"
	"github.com/acexy/golang-toolkit/util/json"
	"testing"
	"time"
)

type StructA struct {
	Time json.Timestamp
	Name string
}

type StructB struct {
	Time time.Time
	Name string
}

type StructA1 struct {
	Time *json.Timestamp
	Name string
}

type StructB1 struct {
	Time *time.Time
	Name string
}

func TestCopy1(t *testing.T) {
	// 添加初始化和调试信息
	var sa StructA
	var sa1 StructA1

	now := time.Now()
	fmt.Println("初始化时间:", now)
	fmt.Println()

	sb := StructB{
		Time: now,
		Name: "Acexy",
	}

	sb1 := StructB1{
		Time: &now,
		Name: "Acexy",
	}

	sa = StructA{}
	fmt.Printf("初始 sa: %+v\n", sa)
	_ = Copy(&sa, &sb)
	fmt.Printf("复制后 sa: %+v\n\n", sa)

	sa = StructA{}
	fmt.Printf("初始 sa: %+v\n", sa)
	_ = Copy(&sa, &sb1)
	fmt.Printf("复制后 sa: %+v\n\n", sa)

	sa1 = StructA1{}
	fmt.Printf("初始 sa1: %+v\n", sa1)
	_ = Copy(&sa1, &sb)
	fmt.Printf("复制后 sa1: %+v\n\n", sa1)

	sa1 = StructA1{}
	fmt.Printf("初始 sa1: %+v\n", sa1)
	_ = Copy(&sa1, &sb1)
	fmt.Printf("复制后 sa1: %+v\n\n", sa1)
}

func TestCopy2(t *testing.T) {
	var sb StructB
	var sb1 StructB1

	sa := StructA{
		Time: json.Timestamp{Time: time.Now()},
		Name: "Acexy",
	}

	sa1 := StructA1{
		Time: &json.Timestamp{Time: time.Now()},
		Name: "Acexy",
	}

	// 从 sa 复制到 sb
	sb = StructB{}
	_ = Copy(&sb, &sa)
	fmt.Println(sb) // 期望输出非零时间值

	// 从 sa1 复制到 sb
	sb = StructB{}
	_ = Copy(&sb, &sa1)
	fmt.Println(sb) // 期望输出非零时间值

	// 从 sa 复制到 sb1
	sb1 = StructB1{}
	_ = Copy(&sb1, &sa)
	fmt.Println(sb1) // 期望输出非零时间值

	// 从 sa1 复制到 sb1
	sb1 = StructB1{}
	_ = Copy(&sb1, &sa1)
	fmt.Println(sb1) // 期望输出非零时间值
}

func Test3(t *testing.T) {
	type StructC struct {
		Time  json.Timestamp
		Time1 time.Time
		Name  string
	}

	type StructD struct {
		Time  time.Time
		Time1 json.Timestamp
		Name  string
	}

	c := &StructC{
		Time:  json.Timestamp{Time: time.Now()},
		Time1: time.Now(),
		Name:  "Acexy",
	}
	var d StructD
	_ = Copy(&d, c)
	fmt.Println(d)

}
