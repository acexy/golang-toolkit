package date

import "time"

// DefaultTimeLayout is the default date and time formatting layout.
const DefaultTimeLayout = "2006-01-02 15:04:05"

//
// -------------------- Parse & Format --------------------
//

// Parse parses with time.Parse and its layout-defined timezone rules.
func Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// Format uses the first supplied layout, or DefaultTimeLayout when absent or empty.
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

// ParseUnixSec converts Unix seconds to a time.
func ParseUnixSec(ts int64) time.Time {
	return time.Unix(ts, 0)
}

// ParseUnixMilli converts Unix milliseconds to a time.
func ParseUnixMilli(ts int64) time.Time {
	return time.UnixMilli(ts)
}

// ParseUnixNano converts Unix nanoseconds to a time.
func ParseUnixNano(ts int64) time.Time {
	return time.Unix(0, ts)
}

//
// -------------------- Unix Timestamp -> Format --------------------
//

// FormatUnixSec formats Unix seconds with an optional layout.
func FormatUnixSec(ts int64, layout ...string) string {
	return Format(ParseUnixSec(ts), layout...)
}

// FormatUnixMilli formats Unix milliseconds with an optional layout.
func FormatUnixMilli(ts int64, layout ...string) string {
	return Format(ParseUnixMilli(ts), layout...)
}

// FormatUnixNano formats Unix nanoseconds with an optional layout.
func FormatUnixNano(ts int64, layout ...string) string {
	return Format(ParseUnixNano(ts), layout...)
}

//
// -------------------- Time -> Unix Timestamp --------------------
//

// ToUnixSec returns Unix seconds for the time.
func ToUnixSec(t time.Time) int64 {
	return t.Unix()
}

// ToUnixMilli returns Unix milliseconds for the time.
func ToUnixMilli(t time.Time) int64 {
	return t.UnixMilli()
}

// ToUnixNano returns Unix nanoseconds for the time.
func ToUnixNano(t time.Time) int64 {
	return t.UnixNano()
}

//
// -------------------- Current Time --------------------
//

// CurrentUnixSec returns the current Unix timestamp in seconds.
func CurrentUnixSec() int64 {
	return time.Now().Unix()
}

// CurrentUnixMilli returns the current Unix timestamp in milliseconds.
func CurrentUnixMilli() int64 {
	return time.Now().UnixMilli()
}

// CurrentUnixNano returns the current Unix timestamp in nanoseconds.
func CurrentUnixNano() int64 {
	return time.Now().UnixNano()
}

// CurrentYear returns the current local year.
func CurrentYear() int {
	return time.Now().Year()
}

// CurrentMonth returns the current local month.
func CurrentMonth() time.Month {
	return time.Now().Month()
}

// CurrentDay returns the current local day of the month.
func CurrentDay() int {
	return time.Now().Day()
}

// CurrentHour returns the current local hour.
func CurrentHour() int {
	return time.Now().Hour()
}

// CurrentMinute returns the current local minute.
func CurrentMinute() int {
	return time.Now().Minute()
}

// CurrentSecond returns the current local second.
func CurrentSecond() int {
	return time.Now().Second()
}

//
// -------------------- Time Add --------------------
//

// AddDays adds calendar days using the time location.
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddHours adds an hour-based duration.
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes adds a minute-based duration.
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds adds a second-based duration.
func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

//
// -------------------- Time Diff --------------------
//

// TimeDiffInSeconds returns t1 minus t2 in seconds, truncating fractional seconds.
func TimeDiffInSeconds(t1, t2 time.Time) int64 {
	return int64(t1.Sub(t2).Seconds())
}

// AbsTimeDiffInSeconds returns the absolute difference in seconds, truncating fractional seconds.
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

// StartOfYear returns the start of the calendar year in the time location.
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the last nanosecond of the calendar year in the time location.
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999_999_999, t.Location())
}

// StartOfMonth returns the start of the calendar month in the time location.
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the last nanosecond of the calendar month in the time location.
func EndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999_999_999, t.Location())
}

// StartOfDay returns the start of the calendar day in the time location.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the last nanosecond of the calendar day in the time location.
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999_999_999, t.Location())
}

//
// -------------------- Calendar Helpers --------------------
//

// WeekdayName returns the English weekday name.
func WeekdayName(t time.Time) string {
	return t.Weekday().String()
}

// Quarter returns the calendar quarter from 1 to 4.
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

// IsAfter reports whether the first time is after the second.
func IsAfter(t time.Time, date time.Time) bool {
	return t.After(date)
}

// IsBefore reports whether the first time is before the second.
func IsBefore(t time.Time, date time.Time) bool {
	return t.Before(date)
}

// IsSameDay compares each time calendar date without normalizing their locations.
func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}

// IsWeekday reports whether the time falls on Monday through Friday.
func IsWeekday(t time.Time) bool {
	w := t.Weekday()
	return w >= time.Monday && w <= time.Friday
}

// IsWeekend reports whether the time falls on Saturday or Sunday.
func IsWeekend(t time.Time) bool {
	w := t.Weekday()
	return w == time.Saturday || w == time.Sunday
}
