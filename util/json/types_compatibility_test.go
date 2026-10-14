package json

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

func TestTimestampNullCompatibility(t *testing.T) {
	for _, unit := range []TimestampType{TimestampTypeMilli, TimestampTypeSecond} {
		for _, tc := range []struct {
			name, raw string
			zero      bool
		}{
			{"exact-null", "null", true},
			{"spaces-null", " null ", false},
			{"json-whitespace-null", "\t\r\nnull\n", false},
			{"numeric-zero", "0", false},
		} {
			t.Run(fmt.Sprintf("%s/unit=%d", tc.name, unit), func(t *testing.T) {
				got, err := Timestamp2TimeWithType([]byte(tc.raw), unit)
				if err != nil {
					t.Fatal(err)
				}
				if tc.zero {
					if !got.IsZero() {
						t.Fatalf("精确 null 应返回零时间: %v", got)
					}
				} else if !got.Equal(time.Unix(0, 0)) || got.IsZero() {
					t.Fatalf("现有行为应返回 Unix 零点: %v", got)
				}
			})
		}
		for _, raw := range []string{" null ", "\t\r\nnull\n"} {
			// Normalization is an explicit caller choice; the original entry point is unchanged.
			got, err := Timestamp2TimeWithType(bytes.Trim([]byte(raw), " \t\r\n"), unit)
			if err != nil || !got.IsZero() {
				t.Errorf("规范化 %q 后 = %v, %v", raw, got, err)
			}
		}
	}
	for _, raw := range []string{"\u00a0null", "nullx", "nul", ""} {
		if _, err := Timestamp2TimeWithType([]byte(raw), TimestampTypeMilli); err == nil {
			t.Errorf("无效 JSON %q 应返回错误", raw)
		}
	}
}

func TestTimestampTypeValidationCompatibility(t *testing.T) {
	for _, unit := range []TimestampType{-1, 99} {
		encoded, err := Time2TimestampWithType(time.Time{}, unit)
		if err != nil || string(encoded) != "0" {
			t.Errorf("零时间应保持类型检查前返回 0: %q, %v", encoded, err)
		}
		decoded, err := Timestamp2TimeWithType([]byte("null"), unit)
		if err != nil || !decoded.IsZero() {
			t.Errorf("精确 null 应保持类型检查前返回零时间: %v, %v", decoded, err)
		}
		if _, err := Time2TimestampWithType(time.Unix(1, 0), unit); !errors.Is(err, toolkitError.ErrUnsupportedTimestampType) {
			t.Errorf("非零时间应检查类型: %v", err)
		}
		for _, raw := range []string{"0", "1", " null "} {
			if _, err := Timestamp2TimeWithType([]byte(raw), unit); !errors.Is(err, toolkitError.ErrUnsupportedTimestampType) {
				t.Errorf("%q 应检查类型: %v", raw, err)
			}
		}
		// Preserve JSON parsing before unit validation and its error precedence.
		if _, err := Timestamp2TimeWithType([]byte("invalid"), unit); err == nil || errors.Is(err, toolkitError.ErrUnsupportedTimestampType) {
			t.Errorf("无效 JSON 应先返回解析错误: %v", err)
		}
	}
}

func TestTimestampZeroRoundTripCompatibility(t *testing.T) {
	for _, unit := range []TimestampType{TimestampTypeMilli, TimestampTypeSecond} {
		for _, value := range []time.Time{time.Time{}, time.Unix(0, 0)} {
			raw, err := Time2TimestampWithType(value, unit)
			if err != nil || string(raw) != "0" {
				t.Fatalf("零时间及 Unix 零点均应编码为 0: %q, %v", raw, err)
			}
			got, err := Timestamp2TimeWithType(raw, unit)
			if err != nil || !got.Equal(time.Unix(0, 0)) || got.IsZero() {
				t.Fatalf("数字 0 应解码为 Unix 零点: %v, %v", got, err)
			}
		}
	}
	var wrapper Timestamp
	raw, err := wrapper.MarshalJSON()
	if err != nil || string(raw) != "0" {
		t.Fatalf("Timestamp 零值编码 = %q, %v", raw, err)
	}
	if err := wrapper.UnmarshalJSON(raw); err != nil || wrapper.IsZero() || !wrapper.Equal(time.Unix(0, 0)) {
		t.Fatalf("Timestamp 应保持原往返语义: %v, %v", wrapper.Time, err)
	}
}

func TestTimestampNullableFieldCompatibility(t *testing.T) {
	type record struct {
		Time *Timestamp `json:"time"`
	}
	for _, tc := range []struct {
		name string
		data record
		want string
	}{
		{"missing", record{}, `{"time":null}`},
		{"epoch", record{Time: &Timestamp{Time: time.Unix(0, 0)}}, `{"time":0}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := ToStringError(tc.data)
			if err != nil || raw != tc.want {
				t.Fatalf("可空字段编码 = %q, %v", raw, err)
			}
			var decoded record
			if err := ParseStringError(raw, &decoded); err != nil {
				t.Fatal(err)
			}
			if tc.data.Time == nil {
				if decoded.Time != nil {
					t.Fatal("缺失时间应保持 nil")
				}
			} else if decoded.Time == nil || !decoded.Time.Equal(time.Unix(0, 0)) {
				t.Fatal("可空字段应区分 Unix 零点与缺失时间")
			}
		})
	}
}

func TestGlobalWrapperSettingCompatibility(t *testing.T) {
	const marker = "json-compatibility-scenario:"
	if scenario, child := strings.CutPrefix(os.Args[len(os.Args)-1], marker); child {
		switch scenario {
		case "default-first":
			if currentTimestampType() != TimestampTypeMilli {
				t.Fatal("默认配置应为毫秒")
			}
			GlobalWrapperSetting()
			var called bool
			GlobalWrapperSetting(func(options *TypeWrapperOptions) {
				called = true
				options.TimestampType = TimestampTypeSecond
			})
			if called || currentTimestampType() != TimestampTypeMilli {
				t.Fatal("首次无参数初始化仍应消耗一次性配置机会")
			}
		case "safe-options":
			// Prepare data first; the callback only assigns synchronously without global reads or retained pointers.
			prepared := TimestampTypeSecond
			GlobalWrapperSetting(WithTimestampType(TimestampTypeMilli), func(options *TypeWrapperOptions) {
				options.TimestampType = prepared
			})
			assertGlobalTimestampSecond(t)
			GlobalWrapperSetting(WithTimestampType(TimestampTypeMilli))
			assertGlobalTimestampSecond(t)
			explicit, err := Time2TimestampWithType(time.Unix(1, 0), TimestampTypeMilli)
			if err != nil || string(explicit) != "1000" || currentTimestampType() != TimestampTypeSecond {
				t.Fatalf("显式单位转换不应改变全局配置: %q, %v", explicit, err)
			}
		case "concurrent-once":
			var calls atomic.Int32
			var workers sync.WaitGroup
			start := make(chan struct{})
			for i := 0; i < 16; i++ {
				workers.Add(1)
				go func() {
					defer workers.Done()
					<-start
					GlobalWrapperSetting(func(options *TypeWrapperOptions) {
						calls.Add(1)
						options.TimestampType = TimestampTypeSecond
					})
					assertGlobalTimestampSecond(t)
				}()
			}
			close(start)
			workers.Wait()
			if got := calls.Load(); got != 1 {
				t.Fatalf("option 应只执行一次，实际 %d", got)
			}
		default:
			t.Fatal("未知场景")
		}
		return
	}
	for _, scenario := range []string{"default-first", "safe-options", "concurrent-once"} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGlobalWrapperSettingCompatibility$", "--", marker+scenario)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("隔离配置测试失败: %v, timeout=%v\n%s", err, ctx.Err(), output)
			}
		})
	}
}

func assertGlobalTimestampSecond(t *testing.T) {
	t.Helper()
	raw, err := Time2Timestamp(time.Unix(1, 0))
	if err != nil || string(raw) != "1" {
		t.Errorf("全局秒级转换 = %q, %v", raw, err)
	}
	got, err := Timestamp2Time([]byte("1"))
	if err != nil || !got.Equal(time.Unix(1, 0)) {
		t.Errorf("全局秒级解析 = %v, %v", got, err)
	}
	wrapper := Timestamp{Time: time.Unix(1, 0)}
	raw, err = wrapper.MarshalJSON()
	if err != nil || string(raw) != "1" {
		t.Errorf("包装器应遵循全局秒级配置: %q, %v", raw, err)
	}
}
