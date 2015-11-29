package logs

import (
	"testing"
)

func TestLogSetLevel(t *testing.T){
	LogSetLevel(LOG_INFO)
}

func TestLogWriteToFile(t *testing.T) {
	LogSetFilePath("D:\\work\\log.txt")
	Start()
	Info("hello web spider log file")
	Info("hello web spider log file")
}

func TestLogWriteToStdOut(t *testing.T) {
	Start()
	Info("write log to stdout")
}

func BenchmarkLogInfo(b *testing.B) {
	Start()
	Info("write log to stdout")
}