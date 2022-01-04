package aurora

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	/*
		levelFormat   = "[%c[%d;%d;%dm%s%s%c[0m] ==> "
		用于定义彩色日志输出格式

		messageFormat = "%c[%d;%d;%dm%s%s%c[0m\n"
		定义彩色消息体

		defaultFormat = "|%s| --> "
		用于定义普通format的输出格式，与彩色模式不兼容
	*/
	levelFormat   = "[%c[%d;%d;%dm%s%s%c[0m] ==> "
	messageFormat = "%c[%d;%d;%dm%s%s%c[0m\n"
	defaultFormat = "[%s] ==> "

	/*
		Info          = iota
		Warning
		Debug
		Error
		枚举定义了日志等级
	*/

	Info = iota
	Warning
	Debug
	Error

	/*
		日志颜色
	*/

	Black     = 30 //黑色
	Rea       = 31 //红色
	Green     = 32 //绿色
	Yellow    = 33 //黄色
	Blue      = 34 //蓝色
	Pink      = 35 //粉色
	DarkGreen = 36 //墨绿
	Grey      = 37 //灰色

	/*
		Log头索引
	*/

	Time = iota
	Level
)

var level = map[int]string{
	Info:    "Infos",
	Warning: "Warns",
	Debug:   "Debug",
	Error:   "Error",
}

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

// Entry 日志条目
type Entry struct {
	Time     string        `json:"TIME"`    //时间
	Head     []interface{} `json:"HEAD"`    //头字段
	Level    string        `json:"LEVEL"`   //日志等级信息
	Message  string        `json:"MESSAGE"` //消息体
	FileInfo string        `json:"CODE"`    //日志调用位置信息
}
type Log struct {
	mu         *sync.Mutex //锁
	out        io.Writer   //控制台输出
	level      map[int]string
	pool       *sync.Pool            //复用 buffer  *bytes.Buffer
	head       []interface{}         //构建日志头字段
	formats    map[int][]interface{} //彩色日志颜色参数
	logFormats map[int][]interface{} //消息格式参数
	path       string
	length     int
}

// NewLog 生成一个日志实例
func NewLog() *Log {
	getwd, _ := os.Getwd()
	return &Log{
		mu:    &sync.Mutex{},
		out:   os.Stdout,
		head:  nil,
		level: level,
		formats: map[int][]interface{}{
			Info:    []interface{}{0: 0x1B, 1: 0, 2: 1, 3: Yellow, 4: "", 5: 0x1B},
			Warning: []interface{}{0: 0x1B, 1: 0, 2: 1, 3: Blue, 4: "", 5: 0x1B},
			Debug:   []interface{}{0: 0x1B, 1: 0, 2: 1, 3: Pink, 4: "", 5: 0x1B},
			Error:   []interface{}{0: 0x1B, 1: 0, 2: 1, 3: Rea, 4: "", 5: 0x1B},
		},
		//   [1]切换显示方式
		//   [2]切换背景色40-47
		//   [3]切换字体颜色 30-37
		logFormats: map[int][]interface{}{
			Info:    []interface{}{0: 0x1B, 1: 0, 2: 1, 3: Green, 4: "", 5: 0x1B},
			Warning: []interface{}{0: 0x1B, 1: 0, 2: 1, 3: Yellow, 4: "", 5: 0x1B},
			Debug:   []interface{}{0: 0x1B, 1: 0, 2: 1, 3: DarkGreen, 4: "", 5: 0x1B},
			Error:   []interface{}{0: 0x1B, 1: 4, 2: 1, 3: Rea, 4: "", 5: 0x1B},
		},
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		path:   getwd,
		length: len(getwd),
	}
}

func (l *Log) Info(info ...interface{}) {
	l.format(Info, levelFormat, info...)
}

func (l *Log) Warning(warning ...interface{}) {
	l.format(Warning, levelFormat, warning...)
}

func (l *Log) Error(err ...interface{}) {
	l.format(Error, levelFormat, err...)
}

func (l *Log) Debug(debug ...interface{}) {
	l.format(Debug, levelFormat, debug...)
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
	sprintf := "" //预备log
	//var entry Entry
	//根据 format 生成 sprintf
	switch format {
	case levelFormat:
		_, sprintf = l.colorFormat(level, format, args...)
	default:
		_, sprintf = l.defaultFormats(format, level, args...)
	}
	if sprintf == "" {
		return
	}
	buffer := l.pool.Get().(*bytes.Buffer)
	_, err := buffer.Write([]byte(sprintf))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	_, err = l.out.Write(buffer.Bytes())
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	//fmt.Printf("%+v\n", entry.Json())
	buffer.Reset()     //刷新缓冲区
	l.pool.Put(buffer) //放入池
}

// Sync 用于同步写入文件日志
func (l *Log) Sync() {
	// 开启线程 读取日志条目

}

// Head 添加一个log头字段
func (l *Log) Head(head map[string]interface{}) {
	if l.head == nil {
		l.head = make([]interface{}, 0)
	}
	l.head = append(l.head, head)
}

// HeadList 批量添加log头字段
func (l *Log) HeadList(head ...map[string]interface{}) {
	if l.head == nil {
		l.head = make([]interface{}, 0)
	}
	for i := 0; i < len(head); i++ {
		l.head = append(l.head, head[i])
	}
}

// LevelColor 修改日志等级输出颜色
func (l *Log) LevelColor(level int, color int) {
	l.formats[level][3] = color
}

// MessageColor 修改日志消息输出颜色
func (l *Log) MessageColor(level int, color int) {
	l.logFormats[level][3] = color
}

// LevelBackground 修改日志等级背景色,不推荐修改背景色，在不同主题下对控制台其他提示会有影响
func (l *Log) LevelBackground(level int, color int) {
	l.formats[level][1] = color + 10
}

// fileInfo 栈信息获取
func (l *Log) fileInfo() string {
	_, file, line, ok := runtime.Caller(4)
	pc := make([]uintptr, 1)
	runtime.Callers(4, pc) //初始化 pc 不能删除
	f := runtime.FuncForPC(pc[0])
	if ok {
		itoa := strconv.Itoa(line)
		return ",code line:" + file[l.length+1:] + ":" + itoa + ",func:" + "(" + f.Name() + ")"
	}
	return ""
}

// colorFormat 彩色日志打印处理方法
// level:日志等级
// format:彩色等级模板参数
// args:用户最终日志消息
func (l *Log) colorFormat(level int, format string, args ...interface{}) (Entry, string) {
	color := l.formats[level]
	lcolor := l.logFormats[level]
	file := l.fileInfo() //获取文件位置信息

	logType := fmt.Sprintf(format, color[0], color[1], color[2], color[3], color[4], l.level[level], color[5]) //解析日志类型及其颜色参数,需要拼接上log头字段

	message := fmt.Sprint(args...) //解析需要打印的日志为字符串

	entry := Entry{
		Time:     time.Now().Format("2006/01/02 15:04:05"),
		Level:    l.level[level],
		Head:     l.head,
		Message:  message,
		FileInfo: file,
	}
	message += file //拼接文件位置信息
	//h := fmt.Sprint(entry.Head...)
	marshal, err := json.Marshal(entry.Head)
	if err != nil {
		return Entry{}, ""
	}
	sprintf := fmt.Sprintf("%s %s"+logType+messageFormat, entry.Time, marshal, lcolor[0], lcolor[1], lcolor[2], lcolor[3], lcolor[4], message, lcolor[5])
	return entry, sprintf
}

// defaultFormats 普通日志处理方法
// format 用户定义的消息模板，不采用彩色模板
func (l *Log) defaultFormats(format string, level int, args ...interface{}) (Entry, string) {
	logType := fmt.Sprintf(defaultFormat, l.level[level])
	message := fmt.Sprint(args...)
	file := l.fileInfo() //获取文件位置信息
	entry := Entry{
		Time:     time.Now().Format("2006/01/02 15:04:05"),
		Level:    l.level[level],
		Head:     l.head,
		Message:  message,
		FileInfo: file,
	}
	marshal, err := json.Marshal(entry.Head)
	if err != nil {
		return Entry{}, ""
	}
	message += file //拼接文件位置信息
	sprintf := fmt.Sprintf("%s %s"+logType+format, entry.Time, marshal, message)
	return entry, sprintf
}

// Json 格式化输出日志明细
func (e *Entry) Json() string {
	if e.Head != nil {
		marshal, err := json.Marshal(e.Head)
		if err != nil {
			return ""
		}
		//		return "{" + " \"TIME\" :" + "\"" + e.Time + "\"" + "," + " \"HEAD\" :" + "\"" + fmt.Sprintf("%s", e.Head) + "\"" + "," + " \"LEVEL\" :" + "\"" + e.Level + "\"" + "," + " \"MESSAGE\" :" + "\"" + e.Message + "\"" + "," + " \"FILE\" :" + "\"" + e.FileInfo + "\"" + "}"
		return "{" + "\"TIME\" :" + "\"" + e.Time + "\"" + "," + "\"HEAD\" :" + string(marshal) + "," + "\"LEVEL\" :" + "\"" + e.Level + "\"" + "," + "\"MESSAGE\" :" + "\"" + e.Message + "\"" + "," + "\"FILE\" :" + "\"" + e.FileInfo + "\"" + "}"
	}
	return "{" + "\"TIME\" :" + "\"" + e.Time + "\"" + "," + "\"LEVEL\" :" + "\"" + e.Level + "\"" + "," + " \"MESSAGE\" :" + "\"" + e.Message + "\"" + "," + "\"FILE\" :" + "\"" + e.FileInfo + "\"" + "}"
}
