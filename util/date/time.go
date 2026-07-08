package date

import "time"

const DefaultTimeLayout = "2006-01-02 15:04:05"

//
// -------------------- Parse & Format --------------------
//

// Parse 解析时间字符串为 time.Time（不带时区，使用 layout 自身语义）
func Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// Format 格式化时间为字符串（默认使用 DefaultTimeLayout）
func Format(t time.Time, layout ...string) string {
	l := DefaultTimeLayout
	if len(layout) > 0 && layout[0] != "" {
		l = layout[0]
	}
	return t.Format(l)
}

//
// -------------------- Unix Timestamp -> Time --------------------
//

// ParseUnixSec 从 Unix 秒级时间戳转为 time.Time
func ParseUnixSec(ts int64) time.Time {
	return time.Unix(ts, 0)
}

// ParseUnixMilli 从 Unix 毫秒级时间戳转为 time.Time
func ParseUnixMilli(ts int64) time.Time {
	return time.UnixMilli(ts)
}

// ParseUnixNano 从 Unix 纳秒级时间戳转为 time.Time
// ts 必须来自 time.Time.UnixNano()
func ParseUnixNano(ts int64) time.Time {
	return time.Unix(0, ts)
}

//
// -------------------- Unix Timestamp -> Format --------------------
//

// FormatUnixSec 将 Unix 秒级时间戳格式化为字符串
func FormatUnixSec(ts int64, layout ...string) string {
	return Format(ParseUnixSec(ts), layout...)
}

// FormatUnixMilli 将 Unix 毫秒级时间戳格式化为字符串
func FormatUnixMilli(ts int64, layout ...string) string {
	return Format(ParseUnixMilli(ts), layout...)
}

// FormatUnixNano 将 Unix 纳秒级时间戳格式化为字符串
func FormatUnixNano(ts int64, layout ...string) string {
	return Format(ParseUnixNano(ts), layout...)
}

//
// -------------------- Time -> Unix Timestamp --------------------
//

func ToUnixSec(t time.Time) int64 {
	return t.Unix()
}

func ToUnixMilli(t time.Time) int64 {
	return t.UnixMilli()
}

func ToUnixNano(t time.Time) int64 {
	return t.UnixNano()
}

//
// -------------------- Current Time --------------------
//

func CurrentUnixSec() int64 {
	return time.Now().Unix()
}

func CurrentUnixMilli() int64 {
	return time.Now().UnixMilli()
}

func CurrentUnixNano() int64 {
	return time.Now().UnixNano()
}

func CurrentYear() int {
	return time.Now().Year()
}

func CurrentMonth() time.Month {
	return time.Now().Month()
}

func CurrentDay() int {
	return time.Now().Day()
}

func CurrentHour() int {
	return time.Now().Hour()
}

func CurrentMinute() int {
	return time.Now().Minute()
}

func CurrentSecond() int {
	return time.Now().Second()
}

//
// -------------------- Time Add --------------------
//

func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

//
// -------------------- Time Diff --------------------
//

// TimeDiffInSeconds 返回 t1 - t2 的秒差（可能为负）
func TimeDiffInSeconds(t1, t2 time.Time) int64 {
	return int64(t1.Sub(t2).Seconds())
}

// AbsTimeDiffInSeconds 返回绝对时间差（秒）
func AbsTimeDiffInSeconds(t1, t2 time.Time) int64 {
	diff := t1.Sub(t2)
	if diff < 0 {
		diff = -diff
	}
	return int64(diff.Seconds())
}

//
// -------------------- Time Range --------------------
//

func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999_999_999, t.Location())
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999_999_999, t.Location())
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999_999_999, t.Location())
}

//
// -------------------- Calendar Helpers --------------------
//

// WeekdayName 星期名
func WeekdayName(t time.Time) string {
	return t.Weekday().String()
}

// Quarter 季度
func Quarter(t time.Time) int {
	switch t.Month() {
	case time.January, time.February, time.March:
		return 1
	case time.April, time.May, time.June:
		return 2
	case time.July, time.August, time.September:
		return 3
	default:
		return 4
	}
}

func IsAfter(t time.Time, date time.Time) bool {
	return t.After(date)
}

func IsBefore(t time.Time, date time.Time) bool {
	return t.Before(date)
}

func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}

func IsWeekday(t time.Time) bool {
	w := t.Weekday()
	return w >= time.Monday && w <= time.Friday
}

func IsWeekend(t time.Time) bool {
	w := t.Weekday()
	return w == time.Saturday || w == time.Sunday
}
