package logs

import (
	"testing"
)

func TestLogSetLevel(t *testing.T){
	LogSetLevel(LOG_INFO)
}

func TestLogWriteToFile(t *testing.T) {
	LogSetFilePath("D:\\work\\log.txt")
	LogStart()
	LogInfo("hello web spider log file")
	LogInfo("hello web spider log file")
}

func TestLogWriteToStdOut(t *testing.T) {
	LogStart()
	LogInfo("write log to stdout")
}

func BenchmarkLogInfo(b *testing.B) {
	LogStart()
	LogInfo("write log to stdout")
}