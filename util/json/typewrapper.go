package json

import (
	"encoding/json"
	"sync"
	"time"
)

type TimeStampType int

const (
	// TimestampTypeMilli 时间戳类型 毫秒级别
	TimestampTypeMilli TimeStampType = 0
	// TimestampTypeSecond 时间戳类型 秒级别
	TimestampTypeSecond TimeStampType = 1
)

var wrapperOnce sync.Once
var defaultOptions *TypeWrapperOptions = &TypeWrapperOptions{

	// 默认毫秒
	TimeStampType: TimestampTypeMilli,
}

type TypeWrapperOptions struct {

	// 序列化/反序列化 时间戳类型
	TimeStampType TimeStampType
}

type SetWrapperOption func(*TypeWrapperOptions)

func GlobalWrapperSetting(opt ...SetWrapperOption) {
	wrapperOnce.Do(func() {
		for _, o := range opt {
			o(defaultOptions)
		}
	})
}

type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	var timestamp int64
	if defaultOptions.TimeStampType == TimestampTypeSecond {
		timestamp = t.Time.Unix()
	} else if defaultOptions.TimeStampType == TimestampTypeMilli {
		timestamp = t.Time.UnixMilli()
	}
	return ToJsonBytesError(timestamp)
}

func (t Timestamp) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	if defaultOptions.TimeStampType == TimestampTypeSecond {
		t.Time = time.Unix(timestamp, 0)
	} else if defaultOptions.TimeStampType == TimestampTypeMilli {
		t.Time = time.UnixMilli(timestamp)
	}
	return nil
}
