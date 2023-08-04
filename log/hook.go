package log

import (
	"github.com/sirupsen/logrus"
)

type autoConsole struct {
}

func (h *autoConsole) Fire(entry *logrus.Entry) error {
	consoleLogger.WithFields(entry.Data).Logln(entry.Level, entry.Message)
	return nil
}

func (h *autoConsole) Levels() []logrus.Level {
	return logrus.AllLevels
}

type LogrusConfig struct {
}
