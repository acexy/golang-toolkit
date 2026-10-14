package logger

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestForwardedLoggerLevels(t *testing.T) {
	const marker = "logger-forwarding-scenario:"
	if scenario, child := strings.CutPrefix(os.Args[len(os.Args)-1], marker); child {
		consoleLevel, fileLevel := DebugLevel, ErrorLevel
		switch scenario {
		case "console-debug":
		case "file-debug":
			consoleLevel, fileLevel = ErrorLevel, DebugLevel
		case "equal":
			fileLevel = DebugLevel
		default:
			t.Fatal("未知场景")
		}
		format := NewFormatter(func(_ TraceIdSupplier, entry *logrus.Entry) ([]byte, error) {
			return []byte(fmt.Sprintf("%s:%v\n", entry.Message, entry.Data["request"])), nil
		})
		// Initialize once per process, configuring console before file output.
		EnableConsoleWithFormatter(consoleLevel, format)
		file := &lumberjack.Logger{Filename: filepath.Join(t.TempDir(), "logrus.log")}
		t.Cleanup(func() { _ = file.Close() })
		EnableFileWithFormatter(fileLevel, format, file)
		var consoleOutput, fileOutput bytes.Buffer
		RawConsoleLogger().SetOutput(&consoleOutput)
		RawFileLogger().SetOutput(&fileOutput)
		Logrus().WithField("request", "request-001").Debug("debug")
		Logrus().WithField("request", "request-001").Error("error")
		for _, output := range []struct {
			name  string
			level Level
			data  *bytes.Buffer
		}{
			{"console", consoleLevel, &consoleOutput},
			{"file", fileLevel, &fileOutput},
		} {
			want := "error:request-001\n"
			if output.level == DebugLevel {
				want = "debug:request-001\n" + want
			}
			if got := output.data.String(); got != want {
				t.Errorf("%s = %q, want %q", output.name, got, want)
			}
		}
		return
	}
	for _, scenario := range []string{"console-debug", "file-debug", "equal"} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestForwardedLoggerLevels$", "--", marker+scenario)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("隔离日志测试失败: %v, timeout=%v\n%s", err, ctx.Err(), output)
			}
		})
	}
}

type countingFormatter struct {
	calls int
	err   error
}

func (f *countingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	f.calls++
	return []byte(entry.Message), f.err
}

func TestForwardedLoggerFiltersBeforeFormatting(t *testing.T) {
	var output bytes.Buffer
	expectedError := errors.New("formatter failed")
	format := &countingFormatter{err: expectedError}
	target := logrus.New()
	target.SetLevel(logrus.ErrorLevel)
	target.SetFormatter(format)
	target.SetOutput(&output)
	hook := &autoConsole{log: target}
	for _, level := range []logrus.Level{logrus.TraceLevel, logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel} {
		if err := hook.Fire(&logrus.Entry{Level: level, Message: "filtered"}); err != nil {
			t.Errorf("被过滤日志不应执行 formatter: %v", err)
		}
	}
	if format.calls != 0 || output.Len() != 0 {
		t.Fatal("被过滤日志不应格式化或写入")
	}
	for _, level := range []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel} {
		if err := hook.Fire(&logrus.Entry{Level: level, Message: "accepted"}); !errors.Is(err, expectedError) {
			t.Errorf("允许级别的格式化错误应正常返回: %v", err)
		}
	}
	if format.calls != 3 || output.Len() != 0 {
		t.Fatal("应格式化三个允许级别的日志，并保留错误时不写入的行为")
	}
}
