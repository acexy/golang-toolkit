package reflect

import (
	"fmt"
	"testing"
)

func TestNonZeroField(t *testing.T) {
	i := 1
	testStruct := struct {
		A string
		B *int
		C bool
		D int
		E []int
		F [1]int
		G map[string]int
	}{
		A: "a",
		B: &i,
		C: true,
	}
	fields, err := NonZeroField(testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Println(fields)
}

func TestGetNonZeroFieldValue(t *testing.T) {
	i := 1
	testStruct := struct {
		A string
		B *int
		C bool
		D int
	}{
		A: "a",
		B: &i,
		C: true,
	}
	value, err := NonZeroFieldValue(testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}

func TestAllField(t *testing.T) {
	i := 1
	testStruct := struct {
		A string
		B *int
		C bool
		D int
		E []int
		F [1]int
		G map[string]int
	}{
		A: "a",
		B: &i,
		C: true,
	}
	fields, err := AllField(testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Println(fields)
}

func TestAllFieldValue(t *testing.T) {
	i := 1
	testStruct := struct {
		A string
		B *int
		C bool
		D int
	}{
		A: "a",
		B: &i,
		C: true,
	}
	value, err := AllFieldValue(testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}

type TestStruct struct {
	IntField    int
	StringField string
	BoolField   bool
	FloatField  float64
	PtrField    *int
	SliceField  []string
	MapField    map[string]int
}

func BenchmarkDeepCopy(b *testing.B) {
	ptrVal := 10
	testStruct := &TestStruct{
		IntField:    10,
		StringField: "Hello",
		BoolField:   true,
		FloatField:  3.14,
		PtrField:    &ptrVal,
		SliceField:  []string{"a", "b", "c"},
		MapField:    map[string]int{"one": 1, "two": 2},
	}

	for i := 0; i < b.N; i++ {
		DeepCopy(testStruct)
	}
}
