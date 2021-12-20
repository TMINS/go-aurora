package aurora

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
)

const (
	Info = iota
	Warning
	Debug
	Error
)
const infoFormat = "%c[%d;%d;%dm%s|Info|%c[0m -->"
const warningFormat = "%c[%d;%d;%dm%s|Warning|%c[0m -->"
const debugFormat = "%c[%d;%d;%dm%s|Debug|%c[0m -->"
const errorFormat = "%c[%d;%d;%dm%s|Error|%c[0m -->"
const messageFormat = "%c[%d;%d;%dm%s%s%c[1m\n"

type Logs interface {
	Info(info ...interface{})
	Warning(warning ...interface{})
	Error(err ...interface{})
	Debug(debug ...interface{})
	Infos(format string, info ...interface{})
	Warnings(format string, warning ...interface{})
	Errors(format string, err ...interface{})
	Debugs(format string, debug ...interface{})
}

type Log struct {
	mu         *sync.Mutex           //锁
	out        io.Writer             //控制台输出
	buffer     *bytes.Buffer         //缓冲区
	pool       *sync.Pool            //复用 buffer  *bytes.Buffer
	formats    map[int][]interface{} //日志格式参数
	logFormats map[int][]interface{} //消息格式参数
	head       map[string]string     //log 头
}

func NewLog() *Log {
	return &Log{
		mu:  &sync.Mutex{},
		out: os.Stdout,
		formats: map[int][]interface{}{
			Info:    []interface{}{0: 0x1B, 1: 33, 2: 46, 3: 1, 4: "", 5: 0x1B},
			Warning: []interface{}{0: 0x1B, 1: 34, 2: 43, 3: 1, 4: "", 5: 0x1B},
			Debug:   []interface{}{0: 0x1B, 1: 36, 2: 40, 3: 1, 4: "", 5: 0x1B},
			Error:   []interface{}{0: 0x1B, 1: 34, 2: 41, 3: 1, 4: "", 5: 0x1B},
		},
		//   [1]切换显示方式
		//   [2]切换背景色40-47
		//   [3]切换字体颜色 30-37
		logFormats: map[int][]interface{}{
			Info:    []interface{}{0: 0x1B, 1: 0, 2: 1, 3: 32, 4: "", 5: 0x1B},
			Warning: []interface{}{0: 0x1B, 1: 0, 2: 1, 3: 33, 4: "", 5: 0x1B},
			Debug:   []interface{}{0: 0x1B, 1: 0, 2: 1, 3: 36, 4: "", 5: 0x1B},
			Error:   []interface{}{0: 0x1B, 1: 4, 2: 1, 3: 31, 4: "", 5: 0x1B},
		},
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (l *Log) Info(info ...interface{}) {
	l.format(Info, infoFormat, info...)
}

func (l *Log) Warning(warning ...interface{}) {
	l.format(Warning, warningFormat, warning...)
}

func (l *Log) Error(err ...interface{}) {
	l.format(Error, errorFormat, err...)
}

func (l *Log) Debug(debug ...interface{}) {
	l.format(Debug, debugFormat, debug...)
}

func (l *Log) Infos(format string, info ...interface{}) {
	l.format(Info, format, info...)
}

func (l *Log) Warnings(format string, warning ...interface{}) {
	l.format(Warning, format, warning...)
}

func (l *Log) Errors(format string, err ...interface{}) {
	l.format(Error, format, err...)
}

func (l *Log) Debugs(format string, debug ...interface{}) {
	l.format(Debug, format, debug...)
}

func (l *Log) format(level int, format string, args ...interface{}) {
	color := l.formats[level] //获取对应日志消息类型及其色彩
	lcolor := l.logFormats[level]
	l.buffer = l.pool.Get().(*bytes.Buffer)
	i := fileInfo()                          //获取文件位置信息
	logType := fmt.Sprintf(format, color...) //解析日志类型及其颜色参数
	message := fmt.Sprint(args...)           //解析需要打印的日志为字符串
	message += " , filepath:" + i            //拼接文件位置信息
	sprintf := fmt.Sprintf(logType+messageFormat, lcolor[0], lcolor[1], lcolor[2], lcolor[3], lcolor[4], message, lcolor[5])
	_, err := l.buffer.Write([]byte(sprintf))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	l.print()
}

func (l *Log) print() {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.out.Write(l.buffer.Bytes())
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	l.buffer.Reset()     //刷新缓冲区
	l.pool.Put(l.buffer) //放入池
}

func fileInfo() string {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		itoa := strconv.Itoa(line)
		return file + ":" + itoa
	}
	return ""
}
