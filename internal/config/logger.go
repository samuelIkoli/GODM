package config

import (
	"fmt"
	"os"
	"runtime"

	log "github.com/jeanphorn/log4go"
)

type Log struct {
	logger *log.Filter
}

var logDir = "./logs"

type Logger struct {
	Level   string
	Date    string
	Source  string
	Message string
}

func GetConfigLogger(level, date, source, message string) *Logger {
	return &Logger{Level: level,
		Date:    date,
		Source:  source,
		Message: message,
	}
}

func NewAppLogger() *Log {
	logSettingsPath := "./log.json"
	_, err := log.ReadFile(logSettingsPath)
	if err != nil {
		fmt.Println(err, "-----")
		return &Log{}
	}
	// check if folder exists
	if _, err = os.Stat(logDir); os.IsNotExist(err) {
		err = os.Mkdir(logDir, os.ModePerm)
		fmt.Println(err)
	}
	log.LoadConfiguration(logSettingsPath)

	return &Log{
		logger: log.LOGGER("filelogs"),
	}
}

func (l *Log) Info(args interface{}, data ...interface{}) {
	l.logger.Log(log.INFO, getSource(), fmt.Sprintf("%s %v", args.(string), data))
}

func (l *Log) Debug(args interface{}, data ...interface{}) {
	l.logger.Log(log.DEBUG, getSource(), fmt.Sprintf("%s %v", args.(string), data))
}

func (l *Log) Warning(args interface{}, data ...interface{}) {
	l.logger.Log(log.WARNING, getSource(), fmt.Sprintf("%s %v", args.(string), data))
}

func (l *Log) Error(args interface{}, data ...interface{}) {
	l.logger.Log(log.ERROR, getSource(), fmt.Sprintf("%s %v", args.(string), data))
}

func (l *Log) Fatal(args interface{}, data ...interface{}) {
	l.logger.Log(log.CRITICAL, getSource(), fmt.Sprintf("%s %v", args.(string), data))
	os.Exit(1)
}

func getSource() string {
	if pc, _, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf("%s:%d", runtime.FuncForPC(pc), line)
	}
	return ""
}
