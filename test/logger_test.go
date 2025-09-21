package test

import (
	"testing"

	"github.com/acexy/golang-toolkit/logger"
	"github.com/acexy/golang-toolkit/sys"
)

func TestWithTraceId(t *testing.T) {
	sys.EnableLocalTraceId(nil)
	logger.Logrus().Info("nice")
}
