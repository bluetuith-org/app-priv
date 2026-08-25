package views

import (
	"fmt"
	"iter"
	"os"
	"strings"
	"sync"
	"time"
)

type logLevel uint8

const (
	logNone logLevel = iota
	logLevelInfo
	logLevelDebug
	logLevelError
)

func (l logLevel) String() string {
	switch l {
	case logLevelInfo:
		return "INFO"

	case logLevelDebug:
		return "DEBUG"

	case logLevelError:
		return "ERROR"
	}

	return ""
}

type logger struct {
	*logBuffer[logData]

	msgCh chan logData

	app AppBinder
}

func newLogger(app AppBinder, filePath string, maxCapacity int) (*logger, error) {
	var file *os.File
	var err error

	if filePath != "" {
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	log := &logger{
		app:       app,
		logBuffer: newLogBuffer[logData](maxCapacity, 10),
		msgCh:     make(chan logData, 1),
	}

	go func() {
		defer close(log.msgCh)

		t := time.NewTimer(250 * time.Millisecond)

		updateCh := log.UpdateChan()
		update := false

		for {
			select {
			case val, ok := <-updateCh:
				if !ok {
					return
				}

				t.Reset(250 * time.Millisecond)
				update = true

				select {
				case log.msgCh <- val:
				default:
				}

			case <-t.C:
				if update {
					log.app.SendMsg(msgLogUpdate())
					update = false
				}
			}
		}
	}()

	if file == nil {
		return log, nil
	}

	go func() {
		defer file.Close()

		var sb strings.Builder

		for val := range log.msgCh {
			sb.WriteString(time.Now().Format("2006-01-02 15:04:05"))
			sb.WriteString(" ")

			sb.WriteString(val.prefix.String())
			sb.WriteString(": ")
			sb.WriteString(val.text)

			if _, err := file.WriteString(sb.String()); err != nil {
				logError(err)
				return
			}

			sb.Reset()
		}
	}()

	return log, nil
}

func (l *logger) Stop() {
	l.Discard()
}

type logData struct {
	prefix logLevel
	text   string
}

func newLogData(prefix logLevel, text string) logData {
	return logData{prefix, text}
}

type logBuffer[T any] struct {
	buffer []T
	ch     chan T

	count, index int
	closed       bool

	mu sync.RWMutex
}

func newLogBuffer[T any](bufCapacity, chanCapacity int) *logBuffer[T] {
	lb := &logBuffer[T]{
		buffer: make([]T, max(1, bufCapacity)),
		ch:     make(chan T, max(2, chanCapacity)),
	}

	return lb
}

func (l *logBuffer[T]) UpdateChan() chan T {
	return l.ch
}

func (l *logBuffer[T]) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.count
}

func (l *logBuffer[T]) Add(val T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return
	}

	l.sendUpdate(val)

	buflen := len(l.buffer)

	l.buffer[l.index] = val
	l.index = (l.index + 1) % buflen

	if l.count < buflen {
		l.count++
	}
}

func (l *logBuffer[T]) Peek() (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.count == 0 {
		var val T
		return val, false
	}

	buflen := len(l.buffer)
	index := ((l.index - 1) + buflen) % buflen

	return l.buffer[index], true
}

func (l *logBuffer[T]) IterAsc() iter.Seq[T] {
	return func(yield func(T) bool) {
		l.mu.RLock()
		defer l.mu.RUnlock()

		startIndex := l.startIndex()
		buflen := len(l.buffer)

		for i := 0; i < l.count; i++ {
			index := (startIndex + i) % buflen
			if !yield(l.buffer[index]) {
				return
			}
		}
	}
}

func (l *logBuffer[T]) IterDesc() iter.Seq[T] {
	return func(yield func(T) bool) {
		l.mu.RLock()
		defer l.mu.RUnlock()

		startIndex := l.startIndex()
		buflen := len(l.buffer)

		for i := l.count - 1; i >= 0; i-- {
			index := (startIndex + i) % buflen
			if !yield(l.buffer[index]) {
				return
			}
		}
	}
}

func (l *logBuffer[T]) Discard() {
	l.mu.Lock()
	defer l.mu.Unlock()

	close(l.ch)
	l.closed = true
}

func (l *logBuffer[T]) startIndex() int {
	if l.count == len(l.buffer) {
		return l.index
	}

	return 0
}

func (l *logBuffer[T]) sendUpdate(val T) {
	if l.ch == nil {
		return
	}

	select {
	case l.ch <- val:
	default:
	}
}

var _log *logger

func initLogger(app AppBinder, filePath string, maxCapacity int) error {
	log, err := newLogger(app, filePath, maxCapacity)
	if err != nil {
		return err
	}

	_log = log

	return nil
}

func stopLogger() {
	_log.Stop()
}

func logFmt(format string, text ...any) string {
	return useStringBuffer(len(text)+len(format), func(b *strings.Builder) {
		fmt.Fprintf(b, format, text...)
	})
}

func logInfo(text string) {
	_log.Add(newLogData(logLevelInfo, text))
}

func logInfof(format string, text ...any) {
	_log.Add(newLogData(logLevelInfo, logFmt(format, text...)))
}

func logError(err error) {
	_log.Add(newLogData(logLevelError, err.Error()))
}

func logErrorf(format string, text ...any) {
	_log.Add(newLogData(logLevelError, logFmt(format, text...)))
}

func logDebug(text string) {
	_log.Add(newLogData(logLevelDebug, text))
}

func logDebugf(format string, text ...any) {
	_log.Add(newLogData(logLevelDebug, logFmt(format, text...)))
}
