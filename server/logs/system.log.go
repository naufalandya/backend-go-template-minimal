package logs

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var processID = os.Getpid()

func Info(message string, metadata map[string]interface{}) {
	log("INFO", message, metadata)
}

func Error(message string, metadata map[string]interface{}) {
	log("ERROR", message, metadata)
}

func log(level, message string, metadata map[string]interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	file, line, function := getCallerInfo()

	metaStr := ""
	if len(metadata) > 0 {
		metaItems := []string{}
		for k, v := range metadata {
			metaItems = append(metaItems, fmt.Sprintf("%s=%v", k, v))
		}
		metaStr = " | " + strings.Join(metaItems, ", ")
	}

	logMsg := fmt.Sprintf("[%s] [%s] [PID:%d] [%s:%d %s] %s%s",
		now, level, processID, file, line, function, message, metaStr,
	)

	if level == "ERROR" {
		fmt.Fprintln(os.Stderr, logMsg)
	} else {
		fmt.Println(logMsg)
	}
}

func getCallerInfo() (string, int, string) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???", 0, "???"
	}
	fn := runtime.FuncForPC(pc)
	functionName := "unknown"
	if fn != nil {
		functionName = fn.Name()
	}

	shortFile := file
	if idx := strings.LastIndex(file, "/"); idx != -1 {
		shortFile = file[idx+1:]
	}

	return shortFile, line, functionName
}
