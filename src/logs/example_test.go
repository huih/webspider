package logs_test

import(
	"logs"
)

func ExampleLog(){
	logs.LogSetFilePath("D:\\work\\log.txt")
	logs.LogStart()
	logs.LogInfo("hello web spider log file: %s", info)
	logs.LogDebug("hello web spider log file: %s", "debug")
	logs.LogWarning("hello web spider log file: %s", "warning")
	logs.LogError("hello web spider log file: %s", "error")
	logs.LogFatal("hello web spider log file: %s", "fatal")
}