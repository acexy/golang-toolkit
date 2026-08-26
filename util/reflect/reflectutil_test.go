package reflect

import (
	"fmt"
	"testing"

	"github.com/acexy/golang-toolkit/util/json"
)

func TestNonZeroField(t *testing.T) {
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
		B: new(1),
		C: true,
	}
	fields, err := NonZeroFieldName(testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Println(fields)
}

func TestGetNonZeroFieldValue(t *testing.T) {
	testStruct := struct {
		A string
		B *int
		C bool
		D int
	}{
		A: "a",
		B: new(1),
		C: true,
	}
	value, err := NonZeroFieldValue(testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}

func TestAllField(t *testing.T) {
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
		B: new(1),
		C: true,
	}
	fields, err := AllFieldName(testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Println(fields)
}

func TestAllFieldValue(t *testing.T) {
	testStruct := struct {
		A string
		B *int
		C bool
		D int
	}{
		A: "a",
		B: new(1),
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
	testStruct := &TestStruct{
		IntField:    10,
		StringField: "Hello",
		BoolField:   true,
		FloatField:  3.14,
		PtrField:    new(10),
		SliceField:  []string{"a", "b", "c"},
		MapField:    map[string]int{"one": 1, "two": 2},
	}

	for b.Loop() {
		DeepCopy(testStruct)
	}
}

func TestSetFieldValue(t *testing.T) {
	var s = new(TestStruct)
	err := SetFieldValue(s, map[string]any{
		"IntField":    10,
		"StringField": "Hello",
		"BoolField":   true,
		"FloatField":  3.14,
		"PtrField":    new(10),
		"SliceField":  []string{"a", "b", "c"},
		"MapField":    map[string]int{"one": 1, "two": 2},
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(json.ToString(s))

	fmt.Println(GetFieldValue(s, "MapField1"))
}
