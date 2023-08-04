package log

import (
	"github.com/sirupsen/logrus"
	"testing"
)

/**
logrus在普通Test模式中由于Format环境变量的原因Console模式的output可能不按预期输出

> Test Ind Debug Model
*/

func TestConsole(t *testing.T) {
	l := &LogrusConfig{}
	l.EnableConsole(logrus.DebugLevel, false) // 非tty模式即使未禁用color也不会生效，自动替换为json模式
	Logrus().Infof("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
	Logrus().WithField("field", "value").Error("error")
}

func TestConsoleDefault(t *testing.T) {
	Logrus().WithField("field", "value").Error("error")
}

func TestFileText(t *testing.T) {
	l := &LogrusConfig{}
	l.EnableFileWithText(logrus.DebugLevel)
	Logrus().Infof("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
}

func TestFileJson(t *testing.T) {
	l := &LogrusConfig{}
	l.EnableFileWithJson(logrus.DebugLevel)
	Logrus().Infof("%d %s\n", 1, "s")
	Logrus().Infoln("Logger Console")
}

func TestConsoleAndFile(t *testing.T) {
	l := &LogrusConfig{}
	l.EnableConsole(logrus.DebugLevel, false)
	l.EnableFileWithJson(logrus.DebugLevel)
	Logrus().WithField("field", "value").Error("error")
}
