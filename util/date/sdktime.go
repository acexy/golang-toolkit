package date

import "time"

// Parse 解析时间字符串为 time.Time
func Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// Format 格式化时间为指定格式字符串
func Format(t time.Time, layout string) string {
	return t.Format(layout)
}

// FormatUnixTimestamp 将 Unix(秒) 时间戳转为指定格式字符串
func FormatUnixTimestamp(ts int64, layout string) string {
	return Format(FromUnixTimestamp(ts), layout)
}

// FormatUnixMilli 将 Unix 毫秒时间戳转为指定格式字符串
func FormatUnixMilli(ts int64, layout string) string {
	return Format(FromUnixMilli(ts), layout)
}

// FormatUnixNano 将 Unix 纳秒时间戳转为指定格式字符串
func FormatUnixNano(ts int64, layout string) string {
	return Format(FromUnixNano(ts), layout)
}

// UnixTimestamp 获取当前 Unix 时间戳（秒）
func UnixTimestamp() int64 {
	return time.Now().Unix()
}

// FromUnixTimestamp 将 Unix 时间戳转为 time.Time
func FromUnixTimestamp(ts int64) time.Time {
	return time.Unix(ts, 0)
}

// UnixNano 获取当前 Unix 纳秒时间戳
func UnixNano() int64 {
	return time.Now().UnixNano()
}

// FromUnixNano 从 Unix 纳秒时间戳转为 time.Time
func FromUnixNano(ts int64) time.Time {
	sec := ts / int64(time.Second)
	nsec := ts % int64(time.Second)
	return time.Unix(sec, nsec)
}

// UnixMilli 获取当前 Unix 毫秒时间戳
func UnixMilli() int64 {
	return time.Now().Unix() * 1000
}

// FromUnixMilli 从 Unix 毫秒时间戳转为 time.Time
func FromUnixMilli(ts int64) time.Time {
	sec := ts / 1000
	nsec := (ts % 1000) * 1000000
	return time.Unix(sec, nsec)
}

// CurrentYear 获取当前时间的年份
func CurrentYear() int {
	return time.Now().Year()
}

// CurrentMonth 获取当前月份
func CurrentMonth() time.Month {
	return time.Now().Month()
}

// CurrentDay 获取当前日期
func CurrentDay() int {
	return time.Now().Day()
}

// AddDays 将时间加上指定的天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddHours 将时间加上指定的小时数
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes 将时间加上指定的分钟数
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds 将时间加上指定的秒数
func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// TimeDiffInSeconds 获取两者时间差（秒）
func TimeDiffInSeconds(t1, t2 time.Time) int64 {
	return int64(t1.Sub(t2).Seconds())
}

// StartOfYear 获取时间的年份开始日期
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 获取时间的年份结束日期
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}

// StartOfMonth 获取时间的月份开始日期
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 获取时间的月份结束日期
func EndOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999999999, t.Location())
}

// WeekdayName 获取时间的星期几（字符串）
func WeekdayName(t time.Time) string {
	return t.Weekday().String()
}

// Quarter 获取时间的季度（1~4）
func Quarter(t time.Time) int {
	month := t.Month()
	switch {
	case month >= 1 && month <= 3:
		return 1
	case month >= 4 && month <= 6:
		return 2
	case month >= 7 && month <= 9:
		return 3
	default:
		return 4
	}
}

// IsAfter 判断当前时间是否在指定日期之后
func IsAfter(t time.Time, date time.Time) bool {
	return t.After(date)
}

// IsBefore 判断当前时间是否在指定日期之前
func IsBefore(t time.Time, date time.Time) bool {
	return t.Before(date)
}

// IsSameDay 判断当前时间是否是同一天
func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}

// CurrentHour 获取当前时间的小时
func CurrentHour() int {
	return time.Now().Hour()
}

// CurrentMinute 获取当前时间的分钟
func CurrentMinute() int {
	return time.Now().Minute()
}

// CurrentSecond 获取当前时间的秒钟
func CurrentSecond() int {
	return time.Now().Second()
}

// StartOfDay 获取当前日期的开始时间（00:00:00）
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay 获取当前日期的结束时间（23:59:59）
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// IsWeekday 获取当前时间是否在工作日
func IsWeekday(t time.Time) bool {
	weekday := t.Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
}

// IsWeekend 获取当前时间是否在周末
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}
