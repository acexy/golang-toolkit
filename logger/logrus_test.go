package logger

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Console colors can differ under non-terminal test output.

func TestConsole(t *testing.T) {
	EnableConsole(DebugLevel) // Exercise console output in the test runner.
	println(IsLevelEnabled(DebugLevel))
	Logrus().Traceln("trace")
	Logrus().Debugf("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
	Logrus().Traceln("trace")
	Logrus().Debugf("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
}

func TestConsoleWithFormatter(t *testing.T) {
	SetTraceIdSupplier(&traceId{})
	EnableConsoleWithFormatter(TraceLevel, NewFormatter(func(trace TraceIdSupplier, entry *logrus.Entry) ([]byte, error) {
		// Format the timestamp with millisecond precision.
		timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
		// Format an uppercase, aligned severity label.
		level := strings.ToUpper(entry.Level.String())
		if len(level) > 5 {
			level = level[:5]
		}
		// Read the original caller file and line.
		file := "unknown:0"
		if entry.HasCaller() {
			file = fmt.Sprintf("%s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line)
		}
		log := fmt.Sprintf("%s %s %-5s [%s] - %s\n", trace.GetTraceId(), timestamp, level, file, entry.Message)
		return []byte(log), nil
	}))
	Logrus().Traceln("trace")
	Logrus().Debugf("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
}
func TestConsoleDefault(t *testing.T) {
	SetTraceIdSupplier(&traceId{})
	Logrus().Traceln("trace")
	Logrus().Debugf("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
}

func TestFileText(t *testing.T) {
	EnableFileWithText(ErrorLevel)
	Logrus().Traceln("trace")
	Logrus().Debugf("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
}

func TestFileJson(t *testing.T) {
	SetTraceIdSupplier(&traceId{})
	EnableConsole(DebugLevel)
	EnableFileWithJson(TraceLevel)
	Logrus().Traceln("trace")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
}

func TestConsoleAndFile(t *testing.T) {
	EnableConsole(DebugLevel, false)
	EnableFileWithJson(DebugLevel, &lumberjack.Logger{
		Filename: "logrus.log",
	})
	Logrus().Debugf("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
}

func TestFileWrite(t *testing.T) {
	EnableConsole(TraceLevel, false)
	EnableFileWithText(TraceLevel)
	Logrus().Info("应用启动中")
}

type traceId struct {
}

func (t *traceId) GetTraceId() string {
	return "traceId"
}

func (t *traceId) SetTraceId(traceId string) {
}

func TestTriceId(t *testing.T) {
	SetTraceIdSupplier(&traceId{})
	Logrus().Infoln("info")
}
