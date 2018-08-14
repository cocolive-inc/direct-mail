package log

import (
	"github.com/ngonghi/VyLog"
	"github.com/ngonghi/VyLog/common"
	"github.com/ngonghi/VyLog/handler"
)

type Logger struct {
	vLog *VyLog.Vylog
}

var l *Logger

func Init() {
	l = NewLogger()
}

func NewLogger() *Logger {
	return &Logger{
		vLog : &VyLog.Vylog{},
	}
}

func GetLogger() *VyLog.Vylog {
	return l.GetLogger()
}

func (l *Logger) GetLogger() *VyLog.Vylog {
	return l.vLog
}

func AddFileHandler(lvName string, logPath string) {
	l.AddFileHandler(lvName, logPath)
}

func (l *Logger) AddFileHandler(lvName string, logPath string) {

	if lvName == "" || logPath == "" {
		return
	}

	lv := VyLog.GetLevelByName(lvName)
	if lv < common.TRACE || lv > common.FATAL {
		return
	}

	handle, err := handler.GetFileHandler(logPath, handler.APPEND)
	if err != nil {
		return
	}

	handle.SetLevel(lv)
	l.vLog.AddHandler(handle)
}

func AddSlackHandler(lvName string, channel string, hookURL string) {
	l.AddSlackHandler(lvName, channel, hookURL)
}

func (l *Logger) AddSlackHandler(lvName string, channel string, hookURL string) {

	if lvName == "" || channel == "" || hookURL == "" {
		return
	}

	lv := VyLog.GetLevelByName(lvName)
	if lv < common.TRACE || lv > common.FATAL {
		return
	}

	handle := handler.GetSlackHandler(channel, hookURL, "メール配信システム", "")
	handle.SetLevel(lv)
	l.vLog.AddHandler(handle)
}