package aurora

import (
	"fmt"
	"runtime"
	"sync"
)

/*
	web 链路链路调用
*/

type Monitor struct {
	FunName  string   //函数名
	FileName string   //函数所在文件名
	CodeLine int      //所在行数
	Err      error    //是否头错误
	Next     *Monitor //下一个调用
}

func (m *Monitor) ToString() string {
	if m.Err != nil {
		s := fmt.Sprintf("method: %s at %s:%d by error: %s ", m.FunName, m.FileName, m.CodeLine, m.Err.Error())
		return s
	}
	s := fmt.Sprintf("method: %s at %s:%d by error: nil ", m.FunName, m.FileName, m.CodeLine)
	return s
}

type LocalMonitor struct {
	mx   *sync.Mutex
	Head *Monitor
	End  *Monitor
}

func (l *LocalMonitor) En(monitor *Monitor) {
	l.mx.Lock()
	defer l.mx.Unlock()
	if l.Head == nil {
		l.Head = monitor
		l.End = monitor
		return
	}
	l.Head.Next = monitor
	l.End = monitor
}

func (l *LocalMonitor) Message() string {
	t := l
	s := "Monitor Error List: start "
	for t.Head != nil {
		s = s + "\n" + t.Head.ToString()
		t.Head = t.Head.Next
	}
	s = s + " end"
	return s
}

func ExecuteInfo(err error) *Monitor {
	caller, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("ExecuteInfo Caller filed! ")
	}
	FunName := runtime.FuncForPC(caller).Name()
	m := &Monitor{
		FunName:  FunName,
		FileName: file,
		CodeLine: line,
		Err:      err,
	}
	return m
}
