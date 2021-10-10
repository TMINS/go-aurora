package logs

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
)

func NewLog() *logrus.Logger {
	log := logrus.New()
	log.Formatter = &WebLog{}
	return log
}

func NewRouteLog() *logrus.Logger {
	log := logrus.New()
	log.Formatter = &RoutLog{}
	return log
}

type WebLog struct {
}

func (wl *WebLog) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}
	//时间格式化
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	//日志格式化
	newLog := fmt.Sprintf("%s [%s] -----> %s\n", timestamp, entry.Level, entry.Message)
	buffer.WriteString(newLog)
	return buffer.Bytes(), nil
}

type RoutLog struct{}

func (rl RoutLog) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}
	//时间格式化
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	//日志格式化
	newLog := fmt.Sprintf("%s %s \n", timestamp, entry.Message)
	buffer.WriteString(newLog)
	return buffer.Bytes(), nil
}
