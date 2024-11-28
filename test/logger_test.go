package test

import (
	"github.com/acexy/golang-toolkit/logger"
	"github.com/acexy/golang-toolkit/sys"
	"testing"
)

func TestWithTraceId(t *testing.T) {
	sys.EnableLocalTraceId(nil)
	logger.Logrus().Info("nice")
}
