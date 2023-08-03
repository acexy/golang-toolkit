package main

import (
	"github.com/acexy/golang-toolkit/log"
	"github.com/sirupsen/logrus"
)

func main() {

	consoleDefault()
	//console()
	//fileText()
	//fileJson()
	//consoleAndFile()
}

func console() {
	l := &log.LogrusConfig{}
	l.EnableConsole(logrus.DebugLevel, false) // 非tty模式即使未禁用color也不会生效，自动替换为json模式

	log.Logrus().Infof("%d %s\n", 1, "s")
	log.Logrus().Infoln("Logger Console")
	log.Logrus().WithField("field", "value").Error("error")
}

func consoleDefault() {
	log.Logrus().WithField("field", "value").Error("error")
}

func fileText() {
	l := &log.LogrusConfig{}
	l.EnableFileWithText(logrus.DebugLevel)
	log.Logrus().Infof("%d %s\n", 1, "s")
	log.Logrus().Infoln("Logger Console")
}

func fileJson() {
	l := &log.LogrusConfig{}
	l.EnableFileWithJson(logrus.DebugLevel)
	log.Logrus().Infof("%d %s\n", 1, "s")
	log.Logrus().Infoln("Logger Console")
}

func consoleAndFile() {
	l := &log.LogrusConfig{}
	l.EnableConsole(logrus.DebugLevel, false)
	l.EnableFileWithJson(logrus.DebugLevel)
	log.Logrus().WithField("field", "value").Error("error")
}
