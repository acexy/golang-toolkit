package log

import (
	"github.com/sirupsen/logrus"
)

type autoConsole struct {
}

func (h *autoConsole) Fire(entry *logrus.Entry) error {
	fn, file := callerPrettyfier(entry.Caller)
	consoleLogger.WithFields(entry.Data).Logln(entry.Level, file, fn, entry.Message)
	return nil
}

func (h *autoConsole) Levels() []logrus.Level {
	return logrus.AllLevels
}
