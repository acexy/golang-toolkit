package logger

import (
	"fmt"
	"github.com/acexy/golang-toolkit/sys"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
	"sync"
)

var (
	callerPrettifier = func(frame *runtime.Frame) (function string, file string) {
		fileName := path.Base(frame.File)
		funcName := path.Base(frame.Function)
		if sys.IsEnabledLocalTraceId() {
			return funcName, fmt.Sprintf("[%s] %s:%v", sys.GetLocalTraceId(), fileName, frame.Line)
		}
		return funcName, fmt.Sprintf("%s:%v", fileName, frame.Line)
	}
	consoleLogger *logrus.Logger
	fileLogger    *logrus.Logger
	activeLogger  *logrus.Logger
	fileSet       bool
	logrusOnce    sync.Once
)

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func enableConsole(level Level, disableColor bool) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.Level(level))
	format := &logrus.TextFormatter{
		FullTimestamp:    true,
		CallerPrettyfier: callerPrettifier,
	}
	if disableColor {
		format.DisableColors = true
	} else {
		format.ForceColors = true
		format.EnvironmentOverrideColors = true
	}
	logger.SetFormatter(format)
	return logger
}

// EnableConsole 启用该设置后，日志内容将向标准控台输出
func EnableConsole(level Level, disableColor bool) {
	consoleLogger = enableConsole(level, disableColor)
	activeLogger = consoleLogger
}

func enableFile(level Level, formatter logrus.Formatter, fileConfig ...*lumberjack.Logger) {
	fileLogger = logrus.New()
	fileLogger.SetReportCaller(true)
	fileLogger.SetLevel(logrus.Level(level))
	fileLogger.SetFormatter(formatter)
	if len(fileConfig) == 0 {
		fileLogger.SetOutput(&lumberjack.Logger{
			Filename:   "./logs/logrus.log",
			MaxSize:    200,
			MaxBackups: 100,
			MaxAge:     365,
			Compress:   true,
		})
	} else {
		fileLogger.SetOutput(fileConfig[0])
	}
}

// EnableFileWithJson 启用该配置后将写入日志文件，并将日志输出json格式
func EnableFileWithJson(level Level, fileConfig ...*lumberjack.Logger) {
	if fileSet {
		panic("repeated initialization")
	}
	fileSet = true
	enableFile(level, &logrus.JSONFormatter{
		DataKey:          "data",
		CallerPrettyfier: callerPrettifier,
	}, fileConfig...)
	if consoleLogger != nil {
		consoleLogger.ReportCaller = false
		fileLogger.AddHook(&autoConsole{})
		activeLogger = fileLogger
	} else {
		activeLogger = fileLogger
	}
}

// EnableFileWithText 启用该配置后写入日志文件，将日志输出为text格式
func EnableFileWithText(level Level, fileConfig ...*lumberjack.Logger) {
	if fileSet {
		panic("repeated initialization")
	}
	fileSet = true
	enableFile(level, &logrus.TextFormatter{
		DisableColors:    true,
		FullTimestamp:    true,
		CallerPrettyfier: callerPrettifier,
	}, fileConfig...)
	if consoleLogger != nil {
		consoleLogger.ReportCaller = false
		fileLogger.AddHook(&autoConsole{})
		activeLogger = fileLogger
	} else {
		activeLogger = fileLogger
	}
}

func Logrus() *logrus.Logger {
	logrusOnce.Do(func() {
		if consoleLogger == nil && fileLogger == nil {
			// 如果未手动初始化，则执行默认初始化配置
			activeLogger = enableConsole(TraceLevel, false)
		}
	})
	return activeLogger
}

type autoConsole struct {
}

func (h *autoConsole) Fire(entry *logrus.Entry) error {
	fn, file := callerPrettifier(entry.Caller)
	consoleLogger.WithFields(entry.Data).Logln(entry.Level, file, fn, entry.Message)
	return nil
}

func (h *autoConsole) Levels() []logrus.Level {
	return logrus.AllLevels
}
