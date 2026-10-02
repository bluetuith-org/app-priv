package ui

import (
	"fmt"
	"os"
	"time"
)

const timeFmt = "2006-01-02 15:04:05"

var lg = logFile{}

func init() {
	f, err := os.OpenFile("app.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err == nil {
		lg.file = f
	}
}

type logFile struct {
	file *os.File
}

// Log logs a message.
// WARN: REMOVE.
func Log(v ...any) {
	if lg.file == nil {
		return
	}

	fmt.Fprintln(lg.file, append([]any{time.Now().Format(timeFmt)}, v...)...)
}

// Logf logs a message
// WARN: REMOVE.
func Logf(f string, v ...any) {
	if lg.file == nil {
		return
	}

	f = "%s " + f

	fmt.Fprintf(lg.file, f, append([]any{time.Now().Format(timeFmt)}, v...)...)
}
