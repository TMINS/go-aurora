package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	INFO = iota					//程序正常运行日志
	WARNING						//警告日志
	ERRORS						//错误输出日志

	INIT
	INIT_ERRORS

	SERVER_REQUEST_INFO
	SERVER_REQUEST_WARNING
	SERVER_REQUEST_ERRORS

	SERVER_RESPONSE_INFO
	SERVER_RESPONSE_WARNING
	SERVER_RESPONSE_ERRORS

)

// 全局默认日志
var logger *Log
var weblog *WebLogs
// Out 面向控制台输出
var Out=os.Stdout

var pc = make([]uintptr,1)
// 日志格式

// Logs 顶级接口，
type Logs interface {
	info (info...interface{})
	warning(warning...interface{})
	errors (errs ...interface{})
	WebLog
}

// Log 日志 顶级
type Log struct {
	Head string
	Logger *log.Logger
}

func (l *Log) Info(info...interface{})  {
	l.Logger.Println(info...)
}

func (l *Log) Error(err...interface{})  {

}

func (l *Log) Warning(warning...interface{})  {

}


func (l Log) info(info...interface{})  {
		l.outMessage(Out,INFO,info...)
}

func (l Log) warning(warning...interface{})  {
	l.outMessage(Out,WARNING,warning...)
}

func (l Log) errors(errs ...interface{})  {
	l.outMessage(Out,ERRORS,errs...)
}
func (l Log) start(errs ...interface{})  {

}

// WebLog Web系统日志接口
type WebLog interface {
	start(message...interface{})
	requestInfo(message...interface{})
}
// WebLogs 继承了Log 从写Logs中的函数来配置需要的模板
type WebLogs struct {
	Log
}

func (l WebLogs) start(errs ...interface{})  {
	l.WebMessage(Out,INIT,errs...)
}

func (l WebLogs) requestInfo(errs ...interface{})  {
	l.WebMessage(Out,SERVER_REQUEST_INFO,errs...)
}

func (l WebLogs) requestError(errs ...interface{})  {
	l.WebMessage(Out,SERVER_REQUEST_ERRORS,errs...)
}


func (l WebLogs) errors(errs ...interface{})  {
	l.WebMessage(Out,INIT_ERRORS,errs...)
}



// OutMessage 重新实现 switch 可以自定义样式
func (l Log) outMessage(write io.Writer,MessageType int,args...interface{})  {
	switch MessageType {
		case INFO:
			fmt.Fprintf(Out,"%c[%dm%s[%s] | TIME:%v | [INFO] -> %s %c[0m\n",0x1B,33,"",funInfo(),LogNowTime(),toMassage(args...),0x1B)
		case WARNING:
			fmt.Fprintf(Out,"%c[%dm%s[%s] | TIME:%v | [WARN] -> %s %c[0m\n",0x1B,33,"",funInfo(),LogNowTime(),toMassage(args...),0x1B)
		case ERRORS:
			fmt.Fprintf(Out,"%c[%dm%s[%s] | TIME:%v | [ERRO] -> %s %c[0m\n",0x1B,31,"",funInfo(),LogNowTime(),toMassage(args...),0x1B)
		default:
	}
}

// WebMessage Web应用log日志模板
func (l WebLogs) WebMessage(write io.Writer,MessageType int,args...interface{})  {
	switch MessageType {
	case INIT:
		fmt.Fprintf(Out,"%c[%dm%s[%s]|TIME:%v|[START] -> %s %c[0m\n",0x1B,32,"",l.Head,LogNowTime(),toMassage(args...),0x1B)
	case SERVER_REQUEST_INFO:
		fmt.Fprintf(Out,"%c[%dm%s[%s][%s]|TIME:%v|[INFO ] -> %s %c[0m\n",0x1B,32,"",l.Head,funInfo(),LogNowTime(),toMassage(args...),0x1B)
	case SERVER_REQUEST_ERRORS:
		fmt.Fprintf(Out,"%c[%dm%s[%s][%s]|TIME:%v|[ERROR] -> %s %c[0m\n",0x1B,31,"",l.Head,funInfo(),LogNowTime(),toMassage(args...),0x1B)
	case INIT_ERRORS:
		fmt.Fprintf(Out,"%c[%dm%s[%s]|TIME:%v|[ERROR] -> %s ,%s %c[0m\n",0x1B,31,"",l.Head,LogNowTime(),toMassage(args...),funInfo(),0x1B)
	default:
	}
}

// BeginLog 显得有些多余 暂时弃用
func BeginLog(log *Log) {
	if log==nil {
		logger=&Log{}
	}else {
		logger=log
	}
}

// LoadWebLog 初始化Web 日志同时初始化 顶级日志
func LoadWebLog(log *WebLogs) {
	if log==nil {
		logger=&Log{}
		weblog=&WebLogs{*logger}
	}else {
		logger=&log.Log
		weblog=log
	}
}



func Info(info...interface{})  {
	logger.info(info...)
}

func Warning(info...interface{})  {
	logger.warning(info...)
}

func Errors(info...interface{})  {
	logger.errors(info...)
}



func WebStart(info...interface{})  {
	weblog.start(info...)
}

func WebError(info...interface{})  {
	weblog.errors(info...)
}

// WebLogger 项目启动 日志输出模板
func WebLogger(args...interface{})  {
	 weblog.start(args...)
}

// WebRequestInfoLogger 项目启动 日志输出模板
func WebRequestInfoLogger(args...interface{})  {
	weblog.start(args...)
}

// WebErrorLogger   项目启动 错误输出
func WebErrorLogger(args...interface{})  {
	weblog.errors(args...)
}

func WebRequestError(args...interface{})  {
	weblog.requestError(args...)
}



func funInfo() string  {
	runtime.Callers(5,pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
func LogNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func toMassage(args...interface{}) string {
	return fmt.Sprint(args...)
}




