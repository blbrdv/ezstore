package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var trace = false
var logFileName *string = nil

func setTraceLevel() {
	ex, _ := os.Executable()
	dir := filepath.Dir(ex)
	data := strings.Split(dir, "go-build")

	s := ""
	if len(data) > 1 {
		trace = true
		s = fmt.Sprintf("%s.log", strings.Replace(data[1], "\\", "", -1))
	}
	logFileName = &s
}

func IsTraceLevel() bool {
	if logFileName == nil {
		setTraceLevel()
	}

	return trace
}

func GetLogFileName() *string {
	if logFileName == nil {
		setTraceLevel()
	}

	return logFileName
}
