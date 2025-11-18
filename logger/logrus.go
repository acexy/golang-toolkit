package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var traceIdSupplier TraceIdSupplier
var autoConsoleHook = new(autoConsole)

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

var (
	callerPrettifier = func(frame *runtime.Frame) (function string, file string) {
		fileName := path.Base(frame.File)
		funcName := path.Base(frame.Function)
		if traceIdSupplier != nil {
			return funcName, fmt.Sprintf("[%s] %s:%v", traceIdSupplier.GetTraceId(), fileName, frame.Line)
		}
		return funcName, fmt.Sprintf("%s:%v", fileName, frame.Line)
	}
	consoleLogger *logrus.Logger
	fileLogger    *logrus.Logger
	activeLogger  *logrus.Logger
	fileSet       bool
	logrusOnce    sync.Once
)

// TraceIdSupplier 生成日志跟踪id
type TraceIdSupplier interface {
	GetTraceId() string
	SetTraceId(string)
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

func enableConsole(level Level, formatter logrus.Formatter, disableColor bool) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.Level(level))
	if formatter != nil {
		logger.SetFormatter(formatter)
	} else {
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
	}
	return logger
}

// EnableConsole 启用该设置后，日志内容将向标准控台输出
func EnableConsole(level Level, disableColor ...bool) {
	var disColor bool
	if len(disableColor) > 0 {
		disColor = disableColor[0]
	}
	activeLogger = enableConsole(level, nil, disColor)
	consoleLogger = activeLogger
}

// EnableConsoleWithFormatter 启用该设置后，日志内容将向标准控台输出
func EnableConsoleWithFormatter(level Level, formatter logrus.Formatter) {
	activeLogger = enableConsole(level, formatter, false)
	consoleLogger = activeLogger
}

// SetTraceIdSupplier 设置TraceIdSupplier
func SetTraceIdSupplier(supplier TraceIdSupplier) {
	if traceIdSupplier == nil {
		if supplier != nil {
			traceIdSupplier = supplier
		}
	}
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
// 如果要使用console+file需要先初始化Console配置
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
		fileLogger.AddHook(autoConsoleHook)
	}
	activeLogger = fileLogger
}

// EnableFileWithText 启用该配置后写入日志文件，将日志输出为text格式
// 如果要使用console+file需要先初始化Console配置
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
		fileLogger.AddHook(autoConsoleHook)
	}
	activeLogger = fileLogger
}

// EnableFileWithFormatter 启用该配置后写入日志文件，将日志输出为指定格式
// 如果要使用console+file需要先初始化Console配置
func EnableFileWithFormatter(level Level, formatter logrus.Formatter, fileConfig ...*lumberjack.Logger) {
	if fileSet {
		panic("repeated initialization")
	}
	fileSet = true
	enableFile(level, formatter, fileConfig...)
	if consoleLogger != nil {
		fileLogger.AddHook(autoConsoleHook)
	}
	activeLogger = fileLogger
}

// Logrus 获取logrus实例
func Logrus() *logrus.Logger {
	logrusOnce.Do(func() {
		if activeLogger == nil {
			// 如果未手动初始化，则执行默认初始化配置
			activeLogger = enableConsole(DebugLevel, nil, false)
		}
	})
	return activeLogger
}

// RawConsoleLogger 获取原始控制台日志实例
func RawConsoleLogger() *logrus.Logger {
	return consoleLogger
}

// RawFileLogger 获取原始文件日志实例
func RawFileLogger() *logrus.Logger {
	return fileLogger
}

// IsLevelEnabled 判断指定级别是否启用 优先以fileLogger实例的状态判断
func IsLevelEnabled(level Level) bool {
	return Logrus().IsLevelEnabled(logrus.Level(level))
}
