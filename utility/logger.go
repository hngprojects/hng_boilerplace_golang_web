package utility

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/davecgh/go-spew/spew"
	log "github.com/jeanphorn/log4go"
)

var (
	AppDir    string
	AppDirErr error
)

// Logger application logger
type Logger struct {
	logger *log.Filter
}

type LogRecord struct {
	Level    string // The log level
	Date     string // The time at which the log message was created (nanoseconds)
	Source   string // The message source
	Message  string // The log message
	Category string // The log group
}

// AuditLog Audit log
type AuditLog struct {
	Date           time.Time   `json:"date"`
	Username       string      `json:"Username"`
	RequestHeader  interface{} `json:"request_header"`
	Request        interface{} `json:"request"`
	StatusCode     int         `json:"status_code"`
	ResponseHeader interface{} `json:"response_header"`
	Response       interface{} `json:"response"`
	ClientID       string      `json:"client_id"`
	Route          string      `json:"route"`
	Duration       float64     `json:"duration (seconds)"`
}

// NewLogger constructs a logger object
func NewLogger() *Logger {
	folder := "./logs"
	logSettingsPath := "./log.json"
	// appDir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Printf("Could not load log location >> ", err)
	// }
	_, err := log.ReadFile(logSettingsPath)
	if err != nil {
		logSettingsPath = "../log.json"
		_, err := log.ReadFile(logSettingsPath)
		if err != nil {
			return &Logger{}
		} else {
			folder = "../logs"
		}
	}

	_ = os.Mkdir(folder, os.ModePerm)

	// log.LoadConfiguration(appDir + string(os.PathSeparator) + logSettingsPath)
	log.LoadConfiguration(logSettingsPath)

	return &Logger{
		logger: log.LOGGER("fileLogs"),
	}
}

// Info log information
func (l *Logger) Info(arg0 interface{}, args ...interface{}) {
	//record := LogRecord{
	//	Level:   "INFO",
	//	Date: time.Now().Local().String(),
	//	Source:  getSource(),
	//	Message: fmt.Sprintf(arg0.(string), args...),
	//}
	//go record.Save(l)
	l.logger.Log(log.INFO, getSource(), fmt.Sprintf(arg0.(string), args...))
}

// Debug log debug
func (l *Logger) Debug(arg0 interface{}, args ...interface{}) {
	//record := LogRecord{
	//	Level:   "DEBUG",
	//	Date: time.Now().Local().String(),
	//	Source:  getSource(),
	//	Message: fmt.Sprintf(arg0.(string), args...),
	//}
	//go record.Save(l)
	l.logger.Log(log.DEBUG, getSource(), fmt.Sprintf(arg0.(string), args...))
}

// Warning log warnings
func (l *Logger) Warning(arg0 interface{}, args ...interface{}) {
	//record := LogRecord{
	//	Level:   "WARNING",
	//	Date: time.Now().Local().String(),
	//	Source:  getSource(),
	//	Message: fmt.Sprintf(arg0.(string), args...),
	//}
	//go record.Save(l)
	l.logger.Log(log.WARNING, getSource(), fmt.Sprintf(arg0.(string), args...))
}

// Error log errors
func (l *Logger) Error(arg0 interface{}, args ...interface{}) {
	//record := LogRecord{
	//	Level:   "ERROR",
	//	Date: time.Now().Local().String(),
	//	Source:  getSource(),
	//	Message: fmt.Sprintf(arg0.(string), args...),
	//}
	//go record.Save(l)
	l.logger.Log(log.ERROR, getSource(), fmt.Sprintf(arg0.(string), args...))
}

// Error log errors
func (l *Logger) mongoError(arg0 interface{}, args ...interface{}) {
	l.logger.Log(log.ERROR, getSource(), fmt.Sprintf(arg0.(string), args...))
}

// Info log errors
func (l *Logger) mongoInfo(arg0 interface{}, args ...interface{}) {
	l.logger.Log(log.INFO, getSource(), fmt.Sprintf(arg0.(string), args...))
}

// Fatal log fatal errors
func (l *Logger) Fatal(arg0 interface{}, args ...interface{}) {
	//record := LogRecord{
	//	Level:   "FATAL",
	//	Date: time.Now().Local().String(),
	//	Source:  getSource(),
	//	Message: fmt.Sprintf(arg0.(string), args...),
	//}
	//go record.Save(l)
	l.logger.Log(log.CRITICAL, getSource(), fmt.Sprintf(arg0.(string), args...))
	l.logger.Close()
	os.Exit(1)
}

// Audit : log information on api request and response
func (l *Logger) Audit(record *AuditLog) {
	js, _ := json.Marshal(record)
	l.logger.Log(log.INFO, getSource(), string(js))
}

func Header2Map(header http.Header) map[string]interface{} {
	head := make(map[string]interface{})
	for k, v := range header {
		head[k] = v
	}
	return head
}

func getSource() (source string) {
	if pc, _, line, ok := runtime.Caller(2); ok {
		source = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), line)
	}
	return
}

func SpewResultForDebugging(description string, v interface{}) {
	fmt.Println()
	fmt.Println("**** Start Result ******")
	fmt.Println(description)
	spew.Dump(v)
	fmt.Println("**** End Result ******")
}
