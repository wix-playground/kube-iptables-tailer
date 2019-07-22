package json_logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func GetLogger(obj string)(logger *logrus.Entry){
	var log = logrus.New()
	log.SetLevel(logrus.InfoLevel)

	log.SetFormatter(&log.JSONFormatter{})
	dc_name := os.Getenv("DC_NAME")
	app_name := os.Getenv("APP_NAME")
	hostname := viper.BindEnv("podname", "HOSTNAME")
	logger =log.WithFields(log.Fields{"dc_name": dc_name, "hostname": hostname,"app_name": app_name, "obj": obj })
	return logger
}





type JsonLogHook struct {
	levels       []logrus.Level
	fileLogEntry *logrus.Entry
}

func NewJsonLogFileHook(fileName string, levelToSet logrus.Level, properties LogProperties) (retVal *JsonLogHook) {
	fileLG := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     30,
		Compress:   true,
	}

	return NewJsonLogHook(levelToSet, properties, fileLG)
}

func NewJsonLogHook(levelToSet logrus.Level, properties LogProperties, writer io.Writer) (retVal *JsonLogHook) {
	logrusLogger := logrus.New()
	logrusLogger.Level = levelToSet
	logrusLogger.Out = writer
	logrusLogger.Formatter = NewLogJsonFormatter()

	newFileLogEntry := newLogEntry(logrusLogger, &properties)

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


