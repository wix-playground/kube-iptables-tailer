package json_logger


import (
"github.com/sirupsen/logrus"
"github.com/stretchr/testify/assert"
"testing"
)

func Test_LevelConfig(t *testing.T) {
	SetLogConfig(Config{
		Level:      "debug",
		Colors:     false,
		Properties: LogProperties{},
	})

	entry := GetLog("test")

	assert.Equal(t, entry.Logger.Level, logrus.DebugLevel)
}
