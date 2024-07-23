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

type StructA struct {
	IntField    int
	StringField string
	BoolField   bool
}

type StructB struct {
	IntField     int
	StringField  string
	AnotherField float64
}

func TestCopySameFields(t *testing.T) {
	structA := StructA{
		IntField:    42,
		StringField: "Hello",
		BoolField:   true,
	}

	structB := StructB{
		IntField:     0,
		StringField:  "",
		AnotherField: 3.14,
	}

	// Copy fields from structA to structB
	err := CopySameFields(&structA, &structB)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("StructA:", structA)
	fmt.Println("StructB:", structB)
}
