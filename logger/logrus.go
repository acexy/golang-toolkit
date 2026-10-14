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

// Level identifies a log severity threshold.
type Level uint32

// Type selects the console or file logger.
type Type int8

const (
	// PanicLevel is the panic severity threshold.
	PanicLevel Level = iota
	// FatalLevel is the fatal severity threshold.
	FatalLevel
	// ErrorLevel is the error severity threshold.
	ErrorLevel
	// WarnLevel is the warn severity threshold.
	WarnLevel
	// InfoLevel is the info severity threshold.
	InfoLevel
	// DebugLevel is the debug severity threshold.
	DebugLevel
	// TraceLevel is the trace severity threshold.
	TraceLevel

	// ConsoleLogger selects console output.
	ConsoleLogger Type = 1
	// FileLogger selects file output.
	FileLogger Type = 2
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

// TraceIdSupplier provides the current trace identifier and a setter for application integration.
type TraceIdSupplier interface {
	// GetTraceId returns the current trace identifier.
	GetTraceId() string
	// SetTraceId sets the current trace identifier.
	SetTraceId(string)
}

type autoConsole struct {
	log *logrus.Logger
}

// Fire forwards an entry only if the destination level allows it.
func (h *autoConsole) Fire(entry *logrus.Entry) error {
	if h.log == nil || h.log.Formatter == nil || h.log.Out == nil {
		return nil
	}
	// Forwarding bypasses the target entry point, so apply its level filter explicitly.
	if !h.log.IsLevelEnabled(entry.Level) {
		return nil
	}
	dup := *entry
	dup.Logger = h.log
	formatted, err := h.log.Formatter.Format(&dup)
	if err != nil {
		return err
	}
	_, err = h.log.Out.Write(formatted)
	return err
}

// Levels subscribes the forwarding hook to every level.
func (h *autoConsole) Levels() []logrus.Level {
	return logrus.AllLevels
}

type formatter struct {
	formatterFn func(trace TraceIdSupplier, entry *logrus.Entry) ([]byte, error)
}

// Format formats an entry with the configured trace supplier.
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return f.formatterFn(traceIdSupplier, entry)
}

// NewFormatter creates a formatter callback using the current trace supplier.
func NewFormatter(formatterFn func(traceSupplier TraceIdSupplier, entry *logrus.Entry) ([]byte, error)) *formatter {
	return &formatter{
		formatterFn: formatterFn,
	}
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

// SetTraceIdSupplier accepts the first nonnil supplier; configure it during startup.
func SetTraceIdSupplier(supplier TraceIdSupplier) {
	if traceIdSupplier == nil {
		if supplier != nil {
			traceIdSupplier = supplier
		}
	}
}

// GetTraceId returns the current trace identifier, or an empty string without a supplier.
func GetTraceId() string {
	if traceIdSupplier != nil {
		return traceIdSupplier.GetTraceId()
	}
	return ""
}

// EnableConsole initializes console output once; repeated initialization panics.
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

// EnableConsoleWithFormatter initializes console output once with a custom formatter.
func EnableConsoleWithFormatter(level Level, formatter *formatter) {
	if consoleLogger != nil {
		panic("repeated initialization")
	}
	consoleLogger = enableConsole(level, formatter, false)
	setActiveLogger()
}

// EnableFileWithJson initializes rotating JSON file output once; configure console first for dual output.
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

// EnableFileWithText initializes rotating text file output once; configure console first for dual output.
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

// EnableFileWithFormatter initializes file output once with a custom formatter and optional rotation settings.
func EnableFileWithFormatter(level Level, formatter *formatter, fileConfig ...*lumberjack.Logger) {
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

// Logrus returns the active logger, initially a debug-level console logger.
func Logrus() *logrus.Logger {
	return activeLogger
}

// RawConsoleLogger returns the configured console logger, or nil before initialization.
func RawConsoleLogger() *logrus.Logger {
	return consoleLogger
}

// RawFileLogger returns the configured file logger, or nil before initialization.
func RawFileLogger() *logrus.Logger {
	return fileLogger
}

// IsLevelEnabled checks a selected logger level; without a type, configured file output takes precedence.
func IsLevelEnabled(level Level, logType ...Type) bool {
	if len(logType) > 0 && logType[0] != 0 {
		if logType[0] == FileLogger {
			if fileLogger == nil {
				return false
			}
			return fileLogger.IsLevelEnabled(logrus.Level(level))
		} else if logType[0] == ConsoleLogger {
			if consoleLogger == nil {
				return false
			}
			return consoleLogger.IsLevelEnabled(logrus.Level(level))
		}
		return false
	}
	if fileLogger != nil {
		return fileLogger.IsLevelEnabled(logrus.Level(level))
	}
	return Logrus().IsLevelEnabled(logrus.Level(level))
}
