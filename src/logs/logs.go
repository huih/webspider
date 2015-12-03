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
	"strings"
	"runtime"
	"sync"
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
	LOG_EXTEND_NULL
)

const(
	LOG_OUTPUT_STDOUT = iota
	LOG_OUTPUT_FILE
)

type LogFile struct{
	baseName string
	extendName string
	suffixName string
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
	mu sync.Mutex
	useprefix bool
}

var localLog = &Log{level:LOG_INFO, logger:nil, 
file:&LogFile{"", "", "", 1024 * 1024 * 5, LOG_EXTEND_DATETIME, "", 0}, 
outputType:LOG_OUTPUT_STDOUT, outputHandle:nil, useprefix:true} 

var fileSeqIndex = 1 

//set level
func LogSetLevel(level int){
	localLog.level = level
}

func LogSetUsePrefix(useprefix bool) {
	localLog.useprefix = useprefix
}

//set filename, for example c:/logs/log.txt
func LogSetFilePath(logPath string){
	
	//set directory and file name
	d, f := path.Split(logPath)
	localLog.file.baseName = f
	localLog.file.filePath = d
	
	//handle file prefix name
	i := strings.LastIndex(f, ".")
	localLog.file.baseName = f[:i]
	localLog.file.suffixName = f[i:]
	
	//set output type
	if len(localLog.file.baseName) <= 0 {
		localLog.outputType = LOG_OUTPUT_STDOUT
	} else {
		localLog.outputType = LOG_OUTPUT_FILE
	}
	
	resetExtendName()
}

func LogSetFileMaxSize(maxSize int) {
	localLog.file.maxSize = maxSize
}

func LogSetOutputType(outType int) {
	localLog.outputType = outType
}

func LogSetFileExtendType(extendType int) {
	localLog.file.extendType = extendType
	resetExtendName()
}

func resetExtendName(){
	if localLog.file.extendType == LOG_EXTEND_SEQUENCE {
		localLog.file.extendName = strconv.Itoa(fileSeqIndex)
	} else if localLog.file.extendType == LOG_EXTEND_DATETIME {
		localLog.file.extendName = time.Now().Format("20060102235959")
	} else if localLog.file.extendType == LOG_EXTEND_NULL {
		localLog.file.extendName = string("")
	}
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

func logserror() bool {
	return localLog.level <= LOG_ERROR
}

func fatal() bool {
	return localLog.level <= LOG_FATAL
}

func start(){
	if localLog.outputType == LOG_OUTPUT_STDOUT {
		localLog.outputHandle = os.Stdout
	} else {
		var fileNameBuffer bytes.Buffer
		fileNameBuffer.WriteString(localLog.file.filePath)
		fileNameBuffer.WriteString(localLog.file.baseName)
		if localLog.file.extendName != string("") {
			fileNameBuffer.WriteString("_")
			fileNameBuffer.WriteString(localLog.file.extendName)
		}
		fileNameBuffer.WriteString(localLog.file.suffixName)
		
		//close old file
		if localLog.outputHandle != nil {
			localLog.outputHandle.Close()
		}
		
		fileHandle, err := os.OpenFile(fileNameBuffer.String(), os.O_WRONLY | os.O_CREATE | os.O_SYNC, 0755)
		if err != nil {
			localLog.outputHandle = os.Stdout
		} else {
			localLog.outputHandle = fileHandle
		}
	}
	
	if (localLog.useprefix){
		localLog.logger = log.New(localLog.outputHandle, "", log.Ldate | log.Ltime)	
	}
}

func logSetCurrentFileSize(size int) {
	localLog.file.currentSize += size
	
	if localLog.outputType != LOG_OUTPUT_FILE {
		return
	}
	
	if localLog.file.currentSize < localLog.file.maxSize {
		return
	}
	
	if localLog.file.extendType == LOG_EXTEND_NULL {
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
	start()
	
}

func output (prefix string, format string, v ...interface{}) {
	localLog.mu.Lock()
	defer localLog.mu.Unlock()
	
	//if destination file is not open
	if localLog.outputHandle == nil {
		start()
	}
	
	if (localLog.useprefix) {
		localLog.logger.SetPrefix(prefix)
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "????"
			line = 0
		}
		_, file = path.Split(file)
		format = fmt.Sprintf("%s:%d %s", file, line, format)
		localLog.logger.Printf(format, v...)
	} else {
		outlog := fmt.Sprintf(format, v...)
		localLog.outputHandle.WriteString(outlog);
	}
	
	logSetCurrentFileSize(len(fmt.Sprintf(format, v...)))
}

func Info(format string, v ...interface{}) {
	if info() {
		output("[info]>> ", format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	if debug(){
		output("[debug]>> ", format, v...)
	}
}

func Warning(format string, v ...interface{}){
	if warning() {
		output("[warning]>> ", format, v...)
	}
}

func Error(format string, v ...interface{}){
	if logserror() {
		output("[error]>> ", format, v...)
	}
}

func Fatal(format string, v ...interface{}){
	if fatal() {
		output("[fatal]>> ", format, v...)
		os.Exit(1)
	}
}
