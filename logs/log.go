package logs

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
)

// NewLog 生成运行时候log打印
func NewLog() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&Log{})
	return l
}

// NewServiceLog 业务日志
func NewServiceLog() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&ServiceLog{})
	return l
}

type Log struct {
}

type ServiceLog struct {
}

func (wl *Log) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}
	//时间格式化
	timestamp := entry.Time.Format("2006/01/02 15:04:05")
	//日志格式化
	newLog := fmt.Sprintf("%s [%-6s] ==> %s\n", timestamp, entry.Level, entry.Message)
	buffer.WriteString(newLog)
	return buffer.Bytes(), nil
}

func (wl *ServiceLog) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}
	//时间格式化
	timestamp := entry.Time.Format("2006/01/02 15:04:05")
	//日志格式化
	newLog := fmt.Sprintf("[%s] [%-6s]:%s\n\t", timestamp, entry.Level, entry.Message)
	buffer.WriteString(newLog)
	return buffer.Bytes(), nil
}
