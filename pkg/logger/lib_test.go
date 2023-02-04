package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombinedLogger_GetLogger(t *testing.T) {
	combinedLogger := New()
	assert.NotNil(t, combinedLogger.GetLogger(0))
	assert.Nil(t, combinedLogger.GetLogger(1))
}
