package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
	"sync"
)

var (
	callerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		fileName := path.Base(frame.File)
		return frame.Function, fileName + fmt.Sprintf(":%v", frame.Line)
	}
	consoleLogger *logrus.Logger
	fileLogger    *logrus.Logger

	activeLogger *logrus.Logger

	fileSet bool

	logrusOnce sync.Once
)

type LogrusConfig struct {
}

func (l *LogrusConfig) EnableConsole(level logrus.Level, disableColor bool) {
	consoleLogger = logrus.New()
	consoleLogger.SetOutput(os.Stdout)
	consoleLogger.SetReportCaller(true)
	consoleLogger.SetLevel(level)
	format := &logrus.TextFormatter{
		FullTimestamp:    true,
		CallerPrettyfier: callerPrettyfier,
	}
	if disableColor {
		format.DisableColors = true
	}
	consoleLogger.SetFormatter(format)
}

func enableFile(level logrus.Level, formatter logrus.Formatter, fileConfig ...*lumberjack.Logger) {
	fileSet = true
	fileLogger = logrus.New()
	fileLogger.SetReportCaller(true)
	fileLogger.SetLevel(level)
	fileLogger.SetFormatter(formatter)
	if len(fileConfig) == 0 {
		fileLogger.SetOutput(&lumberjack.Logger{
			Filename:   "./logrus.log",
			MaxSize:    200,
			MaxBackups: 100,
			MaxAge:     365,
			Compress:   true,
		})
	} else {
		fileLogger.SetOutput(fileConfig[0])
	}
}

func (l *LogrusConfig) EnableFileWithJson(level logrus.Level, fileConfig ...*lumberjack.Logger) {
	if fileSet {
		panic("repeated initialization")
	}
	fileSet = true
	enableFile(level, &logrus.JSONFormatter{
		DataKey:          "data",
		CallerPrettyfier: callerPrettyfier,
	}, fileConfig...)
}

func (l *LogrusConfig) EnableFileWithText(level logrus.Level, fileConfig ...*lumberjack.Logger) {
	if fileSet {
		panic("repeated initialization")
	}
	fileSet = true
	enableFile(level, &logrus.TextFormatter{
		DisableColors:    true,
		FullTimestamp:    true,
		CallerPrettyfier: callerPrettyfier,
	}, fileConfig...)
}

func Logrus() *logrus.Logger {
	logrusOnce.Do(func() {
		if consoleLogger == nil && fileLogger == nil {
			// 如果未手动初始化，则执行默认初始化配置
			config := &LogrusConfig{}
			config.EnableConsole(logrus.TraceLevel, false)
		}
		if consoleLogger != nil && fileLogger != nil {
			fileLogger.AddHook(&autoConsole{})
			activeLogger = fileLogger
		} else {
			if consoleLogger != nil {
				activeLogger = consoleLogger
			} else {
				activeLogger = fileLogger
			}
		}
	})
	return activeLogger
}
