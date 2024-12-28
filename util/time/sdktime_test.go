package time

import (
	"testing"
	"time"
)

// TestParse 测试 Parse 函数。
func TestParse(t *testing.T) {
	tests := []struct {
		layout string
		value  string
		want   time.Time
	}{
		{"2006-01-02", "2023-10-05", time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC)},
		{"2006-01-02T15:04:05Z", "2023-10-05T12:34:56Z", time.Date(2023, 10, 5, 12, 34, 56, 0, time.UTC)},
	}

	for _, test := range tests {
		got, err := Parse(test.layout, test.value)
		if err != nil {
			t.Errorf("Parse(%q, %q) returned error: %v", test.layout, test.value, err)
		}
		if !got.Equal(test.want) {
			t.Errorf("Parse(%q, %q) = %v, want %v", test.layout, test.value, got, test.want)
		}
	}
}

// TestFormat 测试 Format 函数。
func TestFormat(t *testing.T) {
	tests := []struct {
		t      time.Time
		layout string
		want   string
	}{
		{time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC), "2006-01-02", "2023-10-05"},
		{time.Date(0, 0, 0, 12, 34, 56, 0, time.UTC), "15:04:05", "12:34:56"},
		{time.Date(2023, 10, 5, 12, 34, 56, 0, time.UTC), "2006-01-02T15:04:05Z", "2023-10-05T12:34:56Z"},
	}

	for _, test := range tests {
		got := Format(test.t, test.layout)
		if got != test.want {
			t.Errorf("Format(%v, %q) = %q, want %q", test.t, test.layout, got, test.want)
		}
	}
}

// TestFormatUnixTimestamp 测试 FormatUnixTimestamp 函数。
func TestFormatUnixTimestamp(t *testing.T) {
	tests := []struct {
		ts     int64
		layout string
		want   string
	}{
		{1633411200, "2006-01-02", "2021-10-05"},
		{1633411200, "15:04:05", "13:20:00"},
	}

	for _, test := range tests {
		got := FormatUnixTimestamp(test.ts, test.layout)
		if got != test.want {
			t.Errorf("FormatUnixTimestamp(%d, %q) = %q, want %q", test.ts, test.layout, got, test.want)
		}
	}
}

// TestFormatUnixMilli 测试 FormatUnixMilli 函数。
func TestFormatUnixMilli(t *testing.T) {
	tests := []struct {
		ts     int64
		layout string
		want   string
	}{
		{1633411200000, "2006-01-02", "2021-10-05"},
		{1633411200000, "15:04:05", "13:20:00"},
	}

	for _, test := range tests {
		got := FormatUnixMilli(test.ts, test.layout)
		if got != test.want {
			t.Errorf("FormatUnixMilli(%d, %q) = %q, want %q", test.ts, test.layout, got, test.want)
		}
	}
}

// TestFormatUnixNano 测试 FormatUnixNano 函数。
func TestFormatUnixNano(t *testing.T) {
	tests := []struct {
		ts     int64
		layout string
		want   string
	}{
		{1633411200000000000, "2006-01-02", "2021-10-05"},
		{1633411200000000000, "15:04:05", "13:20:00"},
	}

	for _, test := range tests {
		got := FormatUnixNano(test.ts, test.layout)
		if got != test.want {
			t.Errorf("FormatUnixNano(%d, %q) = %q, want %q", test.ts, test.layout, got, test.want)
		}
	}
}

// TestFromUnixTimestamp 测试 FromUnixTimestamp 函数。
func TestFromUnixTimestamp(t *testing.T) {
	tests := []struct {
		ts   int64
		want time.Time
	}{
		{1633411200, time.Date(2021, 10, 5, 13, 20, 0, 0, time.Local)},
	}

	for _, test := range tests {
		got := FromUnixTimestamp(test.ts)
		if !got.Equal(test.want) {
			t.Errorf("FromUnixTimestamp(%d) = %v, want %v", test.ts, got, test.want)
		}
	}
}

// TestFromUnixNano 测试 FromUnixNano 函数。
func TestFromUnixNano(t *testing.T) {
	tests := []struct {
		ts   int64
		want time.Time
	}{
		{1633411200000000000, time.Date(2021, 10, 5, 13, 20, 0, 0, time.Local)},
	}

	for _, test := range tests {
		got := FromUnixNano(test.ts)
		if !got.Equal(test.want) {
			t.Errorf("FromUnixNano(%d) = %v, want %v", test.ts, got, test.want)
		}
	}
}

// TestFromUnixMilli 测试 FromUnixMilli 函数。
func TestFromUnixMilli(t *testing.T) {
	tests := []struct {
		ts   int64
		want time.Time
	}{
		{1633411200000, time.Date(2021, 10, 5, 13, 20, 0, 0, time.Local)},
	}

	for _, test := range tests {
		got := FromUnixMilli(test.ts)
		if !got.Equal(test.want) {
			t.Errorf("FromUnixMilli(%d) = %v, want %v", test.ts, got, test.want)
		}
	}
}

// TestQuarter 测试 Quarter 函数。
func TestQuarter(t *testing.T) {
	tests := []struct {
		t    time.Time
		want int
	}{
		{time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), 1},
		{time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC), 2},
		{time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC), 3},
		{time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC), 4},
	}

	for _, test := range tests {
		got := Quarter(test.t)
		if got != test.want {
			t.Errorf("Quarter(%v) = %d, want %d", test.t, got, test.want)
		}
	}
}

// TestIsAfter 测试 IsAfter 函数。
func TestIsAfter(t *testing.T) {
	tests := []struct {
		t1   time.Time
		t2   time.Time
		want bool
	}{
		{time.Date(2023, 10, 5, 12, 0, 0, 0, time.UTC), time.Date(2023, 10, 5, 11, 0, 0, 0, time.UTC), true},
		{time.Date(2023, 10, 5, 11, 0, 0, 0, time.UTC), time.Date(2023, 10, 5, 12, 0, 0, 0, time.UTC), false},
	}

	for _, test := range tests {
		got := IsAfter(test.t1, test.t2)
		if got != test.want {
			t.Errorf("IsAfter(%v, %v) = %v, want %v", test.t1, test.t2, got, test.want)
		}
	}
}

// TestIsBefore 测试 IsBefore 函数。
func TestIsBefore(t *testing.T) {
	tests := []struct {
		t1   time.Time
		t2   time.Time
		want bool
	}{
		{time.Date(2023, 10, 5, 11, 0, 0, 0, time.UTC), time.Date(2023, 10, 5, 12, 0, 0, 0, time.UTC), true},
		{time.Date(2023, 10, 5, 12, 0, 0, 0, time.UTC), time.Date(2023, 10, 5, 11, 0, 0, 0, time.UTC), false},
	}

	for _, test := range tests {
		got := IsBefore(test.t1, test.t2)
		if got != test.want {
			t.Errorf("IsBefore(%v, %v) = %v, want %v", test.t1, test.t2, got, test.want)
		}
	}
}

// TestIsSameDay 测试 IsSameDay 函数。
func TestIsSameDay(t *testing.T) {
	tests := []struct {
		t1   time.Time
		t2   time.Time
		want bool
	}{
		{time.Date(2023, 10, 5, 12, 0, 0, 0, time.UTC), time.Date(2023, 10, 5, 11, 0, 0, 0, time.UTC), true},
		{time.Date(2023, 10, 5, 12, 0, 0, 0, time.UTC), time.Date(2023, 10, 6, 12, 0, 0, 0, time.UTC), false},
	}

	for _, test := range tests {
		got := IsSameDay(test.t1, test.t2)
		if got != test.want {
			t.Errorf("IsSameDay(%v, %v) = %v, want %v", test.t1, test.t2, got, test.want)
		}
	}
}

// TestIsWeekday 测试 IsWeekday 函数。
func TestIsWeekday(t *testing.T) {
	tests := []struct {
		t    time.Time
		want bool
	}{
		{time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC), true},  // 星期四
		{time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC), false}, // 星期六
	}

	for _, test := range tests {
		got := IsWeekday(test.t)
		if got != test.want {
			t.Errorf("IsWeekday(%v) = %v, want %v", test.t, got, test.want)
		}
	}
}

// TestIsWeekend 测试 IsWeekend 函数。
func TestIsWeekend(t *testing.T) {
	tests := []struct {
		t    time.Time
		want bool
	}{
		{time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC), false}, // 星期四
		{time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC), true},  // 星期六
	}

	for _, test := range tests {
		got := IsWeekend(test.t)
		if got != test.want {
			t.Errorf("IsWeekend(%v) = %v, want %v", test.t, got, test.want)
		}
	}
}
