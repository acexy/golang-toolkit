package logger

import (
	"errors"
	"testing"

	"gopkg.in/natefinch/lumberjack.v2"
)

/**
logrus在普通Test模式中由于Format环境变量的原因Console模式的output可能不按预期输出
*/

func TestConsole(t *testing.T) {
	EnableConsole(DebugLevel, false) // 非tty模式即使未禁用color也不会生效，自动替换为json模式
	Logrus().Debugf("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
}

func TestConsoleDefault(t *testing.T) {
	Logrus().Debugf("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
}

func TestFileText(t *testing.T) {
	EnableFileWithText(ErrorLevel)
	Logrus().Debugf("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithError(errors.New("ERROR")).WithField("field", "value").Error("error")
	Logrus().WithField("field", "value").Traceln("----------------------")
}

func TestFileJson(t *testing.T) {
	EnableFileWithJson(DebugLevel)
	Logrus().Debugf("%d %s\n", 1, "s")
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

func TestTriceId(t *testing.T) {
	SetTraceIdSupplier(&traceId{})
	Logrus().Infoln("info")
}
