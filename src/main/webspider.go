package main

import (
	"fmt"
	"tasklist"
	"logs"
)

func main(){
	
	tasklist.AddPageTask("http://www.baidu.com", "http://www.baidu.com/image")
	tasklist.AddPageTask("http://www.baidu.com1", "http://www.baidu.com/image1")
	pageobject := tasklist.PopPageTask()
	pageobject = tasklist.PopPageTask()
	fmt.Println(pageobject.PageUrl, pageobject.HostUrl)
	
//	//输出到终端
//	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
//	logger.Print("Hello, log file!")

//	//输出到文件
	//logFile,_ := os.OpenFile("log.txt", os.O_WRONLY | os.O_CREATE | os.O_SYNC, 0755) 
//	defer logFile.Close()
//	logger = log.New(logFile, "logger: ", log.Lshortfile)
//	logger.Print("hello write log file")

	//logs.LogSetFilePath("D:\\work\\log.txt")
	logs.LogStart()
	logs.LogInfo("hello web spider log file: %s", "info")
	logs.LogDebug("hello web spider log file: %s", "debug")
	logs.LogWarning("hello web spider log file: %s", "warning")
	logs.LogError("hello web spider log file: %s", "error")
	logs.LogFatal("hello web spider log file: %s", "fatal")
}
