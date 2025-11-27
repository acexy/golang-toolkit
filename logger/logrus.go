package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var traceIdSupplier TraceIdSupplier

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
	activeLogger  = enableConsole(DebugLevel, nil, false)
)

// TraceIdSupplier 生成日志跟踪id
type TraceIdSupplier interface {
	GetTraceId() string
	SetTraceId(string)
}

type autoConsole struct {
	log *logrus.Logger
}

func (h *autoConsole) Fire(entry *logrus.Entry) error {
	fn, file := callerPrettifier(entry.Caller)
	h.log.WithFields(entry.Data).Logln(entry.Level, file, fn, entry.Message)
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

// SetTraceIdSupplier 设置TraceIdSupplier
func SetTraceIdSupplier(supplier TraceIdSupplier) {
	if traceIdSupplier == nil {
		if supplier != nil {
			traceIdSupplier = supplier
		}
	}
}

// EnableConsole 启用该设置后，日志内容将向标准控台输出
func EnableConsole(level Level, disableColor ...bool) {
	if consoleLogger != nil {
		panic("repeated initialization")
	}
	var disColor bool
	if len(disableColor) > 0 {
		disColor = disableColor[0]
	}
	consoleLogger = enableConsole(level, nil, disColor)
	setActiveLogger()
}

// EnableConsoleWithFormatter 启用该设置后，日志内容将向标准控台输出
func EnableConsoleWithFormatter(level Level, formatter logrus.Formatter) {
	if consoleLogger != nil {
		panic("repeated initialization")
	}
	consoleLogger = enableConsole(level, formatter, false)
	setActiveLogger()
}

// EnableFileWithJson 启用该配置后将写入日志文件，并将日志输出json格式
// 如果要使用console+file需要先初始化Console配置
func EnableFileWithJson(level Level, fileConfig ...*lumberjack.Logger) {
	if fileLogger != nil {
		panic("repeated initialization")
	}
	enableFile(level, &logrus.JSONFormatter{
		DataKey:          "data",
		CallerPrettyfier: callerPrettifier,
	}, fileConfig...)
	setActiveLogger()
}

// EnableFileWithText 启用该配置后写入日志文件，将日志输出为text格式
// 如果要使用console+file需要先初始化Console配置
func EnableFileWithText(level Level, fileConfig ...*lumberjack.Logger) {
	if fileLogger != nil {
		panic("repeated initialization")
	}
	enableFile(level, &logrus.TextFormatter{
		DisableColors:    true,
		FullTimestamp:    true,
		CallerPrettyfier: callerPrettifier,
	}, fileConfig...)
	setActiveLogger()
}

// EnableFileWithFormatter 启用该配置后写入日志文件，将日志输出为指定格式
// 如果要使用console+file需要先初始化Console配置
func EnableFileWithFormatter(level Level, formatter logrus.Formatter, fileConfig ...*lumberjack.Logger) {
	if fileLogger != nil {
		panic("repeated initialization")
	}
	enableFile(level, formatter, fileConfig...)
	setActiveLogger()
}

func setActiveLogger() {
	if consoleLogger != nil && fileLogger == nil {
		activeLogger = consoleLogger
	} else if consoleLogger == nil && fileLogger != nil {
		activeLogger = fileLogger
	} else if consoleLogger != nil && fileLogger != nil {
		consoleLevel := consoleLogger.GetLevel()
		fileLevel := fileLogger.GetLevel()
		if consoleLevel < fileLevel {
			fileLogger.AddHook(&autoConsole{
				log: consoleLogger,
			})
			activeLogger = fileLogger
		} else {
			consoleLogger.AddHook(&autoConsole{
				log: fileLogger,
			})
			activeLogger = consoleLogger
		}
	}
}

// Logrus 获取logrus实例
func Logrus() *logrus.Logger {
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
func IsLevelEnabled(level Level, log ...*logrus.Logger) bool {
	if len(log) > 0 && log[0] != nil {
		return log[0].IsLevelEnabled(logrus.Level(level))
	}
	if fileLogger != nil {
		return fileLogger.IsLevelEnabled(logrus.Level(level))
	}
	return Logrus().IsLevelEnabled(logrus.Level(level))
}
