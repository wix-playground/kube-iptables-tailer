package json_logger

import (
	"github.com/sirupsen/logrus"
)

type LogJsonFormatter struct {
	*logrus.JSONFormatter
}

const TimestampFormat = "2006-01-02T15:04:05.999Z07:00"

func NewLogJsonFormatter() *LogJsonFormatter {
	return &LogJsonFormatter{
		&logrus.JSONFormatter{
			TimestampFormat: TimestampFormat,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime: "@timestamp",
				logrus.FieldKeyMsg:  "message",
			},
		}}
}

func (f *LogJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return f.JSONFormatter.Format(entry)
}
