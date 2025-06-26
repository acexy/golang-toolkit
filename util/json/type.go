package json

import (
	"sync"
	"time"
)

const (
	// TimestampTypeMilli 时间戳类型 毫秒级别
	TimestampTypeMilli TimeStampType = 0
	// TimestampTypeSecond 时间戳类型 秒级别
	TimestampTypeSecond TimeStampType = 1
)

var wrapperOnce sync.Once
var defaultOptions = &TypeWrapperOptions{
	// 默认毫秒
	TimestampType: TimestampTypeMilli,
}

type TimeStampType int

type TypeWrapperOptions struct {
	// 序列化/反序列化 时间戳类型
	TimestampType TimeStampType
}

type SetWrapperOption func(*TypeWrapperOptions)

type Timestamp struct {
	time.Time
}

func GlobalWrapperSetting(opt ...SetWrapperOption) {
	wrapperOnce.Do(func() {
		for _, o := range opt {
			o(defaultOptions)
		}
	})
}

func Time2Timestamp(time time.Time) ([]byte, error) {
	if time.IsZero() {
		return ToJsonBytesError(0)
	}
	var timestamp int64
	if defaultOptions.TimestampType == TimestampTypeSecond {
		timestamp = time.Unix()
	} else if defaultOptions.TimestampType == TimestampTypeMilli {
		timestamp = time.UnixMilli()
	}
	return ToJsonBytesError(timestamp)
}

func Timestamp2Time(timestampBytes []byte) (time.Time, error) {
	var timestamp int64
	if err := ParseBytesError(timestampBytes, &timestamp); err != nil {
		return time.Time{}, err
	}
	if defaultOptions.TimestampType == TimestampTypeSecond {
		return time.Unix(timestamp, 0), nil
	} else if defaultOptions.TimestampType == TimestampTypeMilli {
		return time.UnixMilli(timestamp), nil
	}
	return time.Time{}, nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return Time2Timestamp(t.Time)
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	formatTime, err := Timestamp2Time(data)
	if err != nil {
		return err
	}
	t.Time = formatTime
	return nil
}
