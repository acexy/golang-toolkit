package json

import (
	"bytes"
	"sync"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/acexy/golang-toolkit/logger"
)

const (
	// TimestampTypeMilli 时间戳类型 毫秒级别
	TimestampTypeMilli TimestampType = 0
	// TimestampTypeSecond 时间戳类型 秒级别
	TimestampTypeSecond TimestampType = 1
)

var wrapperMu sync.RWMutex
var wrapperOnce sync.Once
var defaultOptions = &TypeWrapperOptions{
	// 默认毫秒
	TimestampType: TimestampTypeMilli,
}

// TimestampType 表示时间戳序列化/反序列化类型
type TimestampType int

// TypeWrapperOptions 表示 JSON 类型包装器的全局配置
type TypeWrapperOptions struct {
	// 序列化/反序列化 时间戳类型
	TimestampType TimestampType
}

// SetWrapperOption 表示 JSON 类型包装器配置项
type SetWrapperOption func(*TypeWrapperOptions)

// Timestamp 用于按全局时间戳配置序列化/反序列化 time.Time
type Timestamp struct {
	time.Time
}

// GlobalWrapperSetting 设置全局 JSON 类型包装器配置，仅首次调用生效
func GlobalWrapperSetting(opt ...SetWrapperOption) {
	initialized := false
	wrapperOnce.Do(func() {
		initialized = true
		wrapperMu.Lock()
		defer wrapperMu.Unlock()
		for _, o := range opt {
			o(defaultOptions)
		}
	})
	if !initialized {
		logger.Logrus().Warningln("json global wrapper setting already initialized")
	}
}

// WithTimestampType 设置时间戳序列化/反序列化类型
func WithTimestampType(timestampType TimestampType) SetWrapperOption {
	return func(options *TypeWrapperOptions) {
		options.TimestampType = timestampType
	}
}

func currentTimestampType() TimestampType {
	wrapperMu.RLock()
	defer wrapperMu.RUnlock()
	return defaultOptions.TimestampType
}

// Time2Timestamp 将 time.Time 按全局配置转换为时间戳 JSON 字节
func Time2Timestamp(value time.Time) ([]byte, error) {
	return Time2TimestampWithType(value, currentTimestampType())
}

// Time2TimestampWithType 将 time.Time 按指定类型转换为时间戳 JSON 字节
func Time2TimestampWithType(value time.Time, timestampType TimestampType) ([]byte, error) {
	if value.IsZero() {
		return ToBytesError(0)
	}
	switch timestampType {
	case TimestampTypeSecond:
		return ToBytesError(value.Unix())
	case TimestampTypeMilli:
		return ToBytesError(value.UnixMilli())
	default:
		return nil, toolkitError.ErrUnsupportedTimestampType
	}
}

// Timestamp2Time 将时间戳 JSON 字节按全局配置转换为 time.Time
func Timestamp2Time(timestampBytes []byte) (time.Time, error) {
	return Timestamp2TimeWithType(timestampBytes, currentTimestampType())
}

// Timestamp2TimeWithType 将时间戳 JSON 字节按指定类型转换为 time.Time
func Timestamp2TimeWithType(timestampBytes []byte, timestampType TimestampType) (time.Time, error) {
	if bytes.Equal(timestampBytes, []byte("null")) {
		return time.Time{}, nil
	}
	var timestamp int64
	if err := ParseBytesError(timestampBytes, &timestamp); err != nil {
		return time.Time{}, err
	}
	switch timestampType {
	case TimestampTypeSecond:
		return time.Unix(timestamp, 0), nil
	case TimestampTypeMilli:
		return time.UnixMilli(timestamp), nil
	default:
		return time.Time{}, toolkitError.ErrUnsupportedTimestampType
	}
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
