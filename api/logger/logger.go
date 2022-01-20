package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	logg "github.com/sirupsen/logrus"
)

type mgformatter struct{}

func (m *mgformatter) Format(entry *logg.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//  Timestamp(time.RFC3339)
	timestamp := entry.Time.Format(time.RFC3339)
	// var newlog string
	// newlog = fmt.Sprintf("[%s][%s]%s\n", timestamp, entry.Level, entry.Message)
	b.WriteString(fmt.Sprintf("[%s][%s]%s,%s\n", timestamp, entry.Level, entry.Message, entry.Data))
	return b.Bytes(), nil
}
func init() {
	// 设置日志格式为json格式
	// logg.SetFormatter(&logg.JSONFormatter{
	// 	PrettyPrint: true,
	// })
	logg.SetFormatter(&mgformatter{})
	// 设置将日志输出到标准输出（默认的输出为stderr,标准错误）
	// 日志消息输出可以是任意的io.writer类型
	logg.SetOutput(os.Stdout)
	// 设置日志级别为info以上
	logg.SetLevel(logg.InfoLevel)

}

func InfoFields(l interface{}, f logg.Fields) {
	// logg.SetFormatter(&mgformatter{})
	entry := logg.WithFields(f)
	//logg.WithFields(logg.Fields(f))
	entry.Data["file"] = fileInfo(2)
	entry.Info(l)
}
func InfoJsonFields(l interface{}, f logg.Fields) {

	logg.SetFormatter(&logg.JSONFormatter{
		PrettyPrint: true,
	})
	entry := logg.WithFields(logg.Fields(f))
	entry.Data["file"] = fileInfo(2)
	entry.Info(l)
	logg.SetFormatter(&mgformatter{})
}
func InfoJsonField(l interface{}, k string, v interface{}) {
	logg.SetFormatter(&logg.JSONFormatter{
		PrettyPrint: true,
	})
	entry := logg.WithField(k, v)
	entry.Data["file"] = fileInfo(2)
	entry.Info(l)
	logg.SetFormatter(&mgformatter{})
}
func InfoField(l interface{}, k string, v interface{}) {
	// logg.SetFormatter(&mgformatter{})
	entry := logg.WithField(k, v)
	entry.Data["file"] = fileInfo(2)
	entry.Info(l)
}
func ErrorField(l interface{}, k string, v interface{}) {
	entry := logg.WithField(k, v)
	entry.Data["file"] = fileInfo(2)
	entry.Error(l)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
