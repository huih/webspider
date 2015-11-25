//Copyright 2015 huan ghui. All rights reserved.
package logs

import(
	"log"
	"path"
	"os"
	"time"
	"strconv"
	"bytes"
	"fmt"
)

const (
	LOG_INFO = (iota + 1) * 10
	LOG_DEBUG
	LOG_WARNING
	LOG_ERROR
	LOG_FATAL
)

const (
	LOG_EXTEND_SEQUENCE = iota
	LOG_EXTEND_DATETIME
)

const(
	LOG_OUTPUT_STDOUT = iota
	LOG_OUTPUT_FILE
)

type LogFile struct{
	baseName string
	extendName string
	maxSize int // the max size of log file
	extendType int //the value must is LOG_EXTEND_SEQUENCE,LOG_EXTEND_DATETIME
	filePath string
	currentSize int //the current size of file
}

type Log struct{
	level int
	logger *log.Logger
	file *LogFile
	outputType int
	outputHandle *os.File
}

var localLog = &Log{LOG_INFO, nil, 
&LogFile{"","", 1024 * 1024 * 5, LOG_EXTEND_DATETIME, "", 0}, LOG_OUTPUT_STDOUT, nil} 

var fileSeqIndex = 1 

//set level
func LogSetLevel(level int){
	localLog.level = level
}

//set filename, for example c:/logs/log.txt
func LogSetFilePath(logPath string){
	
	//set directory and file name
	d, f := path.Split(logPath)
	localLog.file.baseName = f
	localLog.file.filePath = d
	
	//set output type
	if len(localLog.file.baseName) <= 0 {
		localLog.outputType = LOG_OUTPUT_STDOUT
	} else {
		localLog.outputType = LOG_OUTPUT_FILE
	}
	
	if localLog.file.extendType == LOG_EXTEND_SEQUENCE {
		localLog.file.extendName = strconv.Itoa(fileSeqIndex)
	} else if localLog.file.extendType == LOG_EXTEND_DATETIME {
		localLog.file.extendName = time.Now().Format("20060102235959")
	}
	
}

func LogSetFileMaxSize(maxSize int) {
	localLog.file.maxSize = maxSize
}

func LogSetOutputType(outType int) {
	localLog.outputType = outType
}

func LogSetFileExtendType(extendType int) {
	localLog.file.extendType = extendType	
}

func info()bool {
	return localLog.level <= LOG_INFO
}

func debug()bool {
	return localLog.level <= LOG_DEBUG
}

func warning() bool {
	return localLog.level <= LOG_WARNING
}

func error() bool {
	return localLog.level <= LOG_ERROR
}

func fatal() bool {
	return localLog.level <= LOG_FATAL
}

func LogStart(){
	if localLog.outputType == LOG_OUTPUT_STDOUT {
		fmt.Print("xxxxxxxxxxxxxx")
		localLog.outputHandle = os.Stdout
	} else {
		var fileNameBuffer bytes.Buffer
		fileNameBuffer.WriteString(localLog.file.filePath)
		fileNameBuffer.WriteString(localLog.file.baseName)
		fileNameBuffer.WriteString("_")
		fileNameBuffer.WriteString(localLog.file.extendName)
		
		fileHandle, err := os.OpenFile(fileNameBuffer.String(), os.O_WRONLY | os.O_CREATE | os.O_SYNC, 0755)
		if err != nil {
			localLog.outputHandle = os.Stdout
		} else {
			localLog.outputHandle = fileHandle
		}
	}
	
	localLog.logger = log.New(localLog.outputHandle, "", log.Lshortfile)
}

func logSetCurrentFileSize(size int) {
	localLog.file.currentSize += size
	
	if localLog.outputType != LOG_OUTPUT_FILE {
		return
	}
	
	if localLog.file.currentSize < localLog.file.maxSize {
		return
	}
	
	//set new file
	localLog.outputHandle.Close()
	
	//new file index
	if localLog.file.extendType == LOG_EXTEND_SEQUENCE {
		fileSeqIndex += 1
		localLog.file.extendName = strconv.Itoa(fileSeqIndex)
	} else if localLog.file.extendType == LOG_EXTEND_DATETIME {
		localLog.file.extendName = time.Now().Format("20060102235959")
	}
	
	//start new log file
	LogStart()
	
}
func LogInfo(log string) {
	if info() == false {
		return
	}
	localLog.logger.Print(log)
	
	logSetCurrentFileSize(len(log))
}

