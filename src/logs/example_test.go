package logs_test

import(
	"logs"
)

func ExampleLog(){
	logs.LogSetFilePath("D:\\work\\log.txt")
	logs.LogStart()
	logs.LogInfo("hello web spider log file")
	logs.LogInfo("hello web spider log file")
}