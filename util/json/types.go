package json

import (
	"bytes"
	"sync"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/acexy/golang-toolkit/logger"
)

const (
	// TimestampTypeMilli selects millisecond timestamps.
	TimestampTypeMilli TimestampType = 0
	// TimestampTypeSecond selects second timestamps.
	TimestampTypeSecond TimestampType = 1
)

var wrapperMu sync.RWMutex
var wrapperOnce sync.Once
var defaultOptions = &TypeWrapperOptions{
	// Default to milliseconds.
	TimestampType: TimestampTypeMilli,
}

// TimestampType identifies the JSON timestamp unit.
type TimestampType int

// TypeWrapperOptions holds process-wide JSON wrapper settings.
type TypeWrapperOptions struct {
	// TimestampType selects the timestamp serialization unit.
	TimestampType TimestampType
}

// SetWrapperOption synchronously updates the supplied settings. Do not retain or asynchronously mutate the pointer.
type SetWrapperOption func(*TypeWrapperOptions)

// Timestamp wraps time.Time using the global timestamp unit.
// Zero time encodes as 0, which decodes as the Unix epoch; use a nil *Timestamp for a missing value.
type Timestamp struct {
	time.Time
}

// GlobalWrapperSetting applies only the first configuration call.
// Options run under the write lock: only assign to their parameter and do not retain its pointer.
// Reading global settings or recursively configuring them inside an option deadlocks.
// Prepare option data before calling this function.
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

// WithTimestampType sets the timestamp unit in an option.
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

// Time2Timestamp encodes a time using the global timestamp unit.
func Time2Timestamp(value time.Time) ([]byte, error) {
	return Time2TimestampWithType(value, currentTimestampType())
}

// Time2TimestampWithType encodes a time using an explicit unit.
// Zero time returns 0 before validating the unit; validate it first when strict checking is required.
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

// Timestamp2Time decodes a timestamp using the global unit.
func Timestamp2Time(timestampBytes []byte) (time.Time, error) {
	return Timestamp2TimeWithType(timestampBytes, currentTimestampType())
}

// Timestamp2TimeWithType decodes a timestamp using an explicit unit.
// Exact null returns zero time before unit validation; whitespace-padded null is treated as numeric 0.
// Callers may normalize JSON whitespace first and validate the unit explicitly for strict behavior.
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

// MarshalJSON encodes the wrapped time using global settings.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return Time2Timestamp(t.Time)
}

// UnmarshalJSON decodes into the wrapped time using global settings.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	formatTime, err := Timestamp2Time(data)
	if err != nil {
		return err
	}
	t.Time = formatTime
	return nil
}
