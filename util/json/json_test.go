package json

import (
	"errors"
	"testing"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

type person struct {
	Name string `json:"name"`
	Sex  uint   `json:"sex"`
}

type student struct {
	Name   string `json:"name"`
	Sex    uint   `json:"sex"`
	School string `json:"school"`
}

func TestToStringAndParseString(t *testing.T) {
	raw := person{Name: "acexy", Sex: 1}
	jsonString := ToString(raw)
	if jsonString != `{"name":"acexy","sex":1}` {
		t.Fatalf("unexpected json string: %s", jsonString)
	}

	var parsed person
	if err := ParseStringError(jsonString, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed != raw {
		t.Fatalf("unexpected parsed value: %+v", parsed)
	}
}

func TestToStringFormat(t *testing.T) {
	jsonString, err := ToStringFormatError(person{Name: "acexy", Sex: 1})
	if err != nil {
		t.Fatal(err)
	}
	expected := "{\n  \"name\": \"acexy\",\n  \"sex\": 1\n}"
	if jsonString != expected {
		t.Fatalf("unexpected formatted json:\n%s", jsonString)
	}
}

func TestCopyStruct(t *testing.T) {
	source := student{Name: "acexy", Sex: 1, School: "Q"}
	var target person
	if err := CopyStructError(source, &target); err != nil {
		t.Fatal(err)
	}
	if target != (person{Name: "acexy", Sex: 1}) {
		t.Fatalf("unexpected copied value: %+v", target)
	}
}

type timestampUser struct {
	Name string     `json:"name"`
	Time *Timestamp `json:"time"`
}

func TestTimestampMilli(t *testing.T) {
	baseTime := time.UnixMilli(1729136314000)
	user := timestampUser{
		Name: "acexy",
		Time: &Timestamp{baseTime},
	}
	jsonString := ToString(user)
	if jsonString != `{"name":"acexy","time":1729136314000}` {
		t.Fatalf("unexpected timestamp json: %s", jsonString)
	}

	var parsed timestampUser
	if err := ParseStringError(jsonString, &parsed); err != nil {
		t.Fatal(err)
	}
	if !parsed.Time.Equal(baseTime) {
		t.Fatalf("unexpected parsed time: %s", parsed.Time)
	}
}

func TestTimestampSecond(t *testing.T) {
	baseTime := time.Unix(1729136314, 0)
	jsonBytes, err := Time2TimestampWithType(baseTime, TimestampTypeSecond)
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonBytes) != "1729136314" {
		t.Fatalf("unexpected timestamp: %s", string(jsonBytes))
	}

	parsed, err := Timestamp2TimeWithType([]byte("1729136314"), TimestampTypeSecond)
	if err != nil {
		t.Fatal(err)
	}
	if !parsed.Equal(baseTime) {
		t.Fatalf("unexpected parsed time: %s", parsed)
	}
}

func TestTimestampNull(t *testing.T) {
	parsed, err := Timestamp2Time([]byte("null"))
	if err != nil {
		t.Fatal(err)
	}
	if !parsed.IsZero() {
		t.Fatalf("expected zero time, got %s", parsed)
	}
}

func TestUnsupportedTimestampType(t *testing.T) {
	if _, err := Time2TimestampWithType(time.Now(), TimestampType(99)); !errors.Is(err, toolkitError.ErrUnsupportedTimestampType) {
		t.Fatalf("expected ErrUnsupportedTimestampType, got %v", err)
	}
	if _, err := Timestamp2TimeWithType([]byte("1"), TimestampType(99)); !errors.Is(err, toolkitError.ErrUnsupportedTimestampType) {
		t.Fatalf("expected ErrUnsupportedTimestampType, got %v", err)
	}
}
