package json_logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
)

var loggerEntry *logrus.Entry

func GetLog(obj string) (retval *logrus.Entry) {
	retval = logrus.New().WithField("obj", obj)
	if loggerEntry != nil {
		retval = loggerEntry.WithField("obj", obj)
	}
	return retval
}

func newLogEntry(logger *logrus.Logger) (retval *logrus.Entry) {
	dc_name := os.Getenv("DC_NAME")
	app_name := os.Getenv("APP_NAME")
	hostname := os.Getenv("HOST_NAME")
	retval = logrus.NewEntry(logger).WithFields(logrus.Fields{"dc_name": dc_name,
		"hostname": hostname,
		"app_name": app_name})
	return retval
}

type JsonLogHook struct {
	levels       []logrus.Level
	fileLogEntry *logrus.Entry
}

func NewJsonLogFileHook(fileName string, levelToSet logrus.Level) (retVal *JsonLogHook) {
	fileLG := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     30,
		Compress:   true,
	}

	return NewJsonLogHook(levelToSet, fileLG)
}

func NewJsonLogHook(levelToSet logrus.Level, writer io.Writer) (retVal *JsonLogHook) {
	logrusLogger := logrus.New()
	logrusLogger.Level = levelToSet
	logrusLogger.Out = writer
	logrusLogger.Formatter = NewLogJsonFormatter()

	newFileLogEntry := newLogEntry(logrusLogger)

	levels := make([]logrus.Level, 0)
	for _, nextLevel := range logrus.AllLevels {
		levels = append(levels, nextLevel)
		if int32(nextLevel) >= int32(levelToSet) {
			break
		}
	}

	retVal = &JsonLogHook{
		levels:       levels,
		fileLogEntry: newFileLogEntry,
	}
	return retVal
}

// Fire is required to implement Logrus hook
func (this *JsonLogHook) Fire(entry *logrus.Entry) error {
	type printMethod func(args ...interface{})
	var funcToCallForPrint printMethod

	entryTolog := this.fileLogEntry.WithFields(entry.Data)
	switch entry.Level {
	case logrus.DebugLevel:
		funcToCallForPrint = entryTolog.Debug
	case logrus.InfoLevel:
		funcToCallForPrint = entryTolog.Info
	case logrus.WarnLevel:
		funcToCallForPrint = entryTolog.Warn
	case logrus.ErrorLevel:
		funcToCallForPrint = entryTolog.Error
	case logrus.FatalLevel:
		funcToCallForPrint = entryTolog.Fatal
	case logrus.PanicLevel:
		funcToCallForPrint = entryTolog.Panic
	}
	funcToCallForPrint(entry.Message)
	return nil
}

// Levels Required for logrus hook implementation
func (hook *JsonLogHook) Levels() []logrus.Level {
	return hook.levels
}

func ConfigureLogger(fileName string, logLevel string) {
	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	levelStr := logLevel
	logger := logrus.New()
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(level)

	fileHook := NewJsonLogFileHook(fileName, level)
	logger.Hooks.Add(fileHook)

	loggerEntry = logrus.NewEntry(logger)
	GetLog("logging").Info("Logging module configured successfully ", fileName, logLevel)
}
