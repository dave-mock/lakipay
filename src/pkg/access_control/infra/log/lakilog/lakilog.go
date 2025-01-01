package lakilog

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu     sync.Mutex
	prefix string
}

func New(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
	}
}

func (l *Logger) Prefix() string {
	return l.prefix
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Info(msg string) {
	// Append color
	// Red 31
	// Green 32
	// Yellow 33
	l.mu.Lock()
	defer l.mu.Unlock()

	l.mu.Unlock()
	var ok bool
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	l.mu.Lock()
	fmt.Printf("\x1b[30m:: %s :: %s:%d :: %s :: %s\x1b\n", time.Now().Format("2006-01-02 15:04:05.000"), strings.Split(file, "/")[len(strings.Split(file, "/"))-1], line, l.prefix, msg)
}
func (l *Logger) Warning(msg string) {
	// Append color
	// Red 31
	// Green 32
	// Yellow 33
	l.mu.Lock()
	defer l.mu.Unlock()

	l.mu.Unlock()
	var ok bool
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	l.mu.Lock()
	fmt.Printf("\x1b[33m:: %s :: %s:%d :: %s :: %s\x1b\n", time.Now().Format("2006-01-02 15:04:05.000"), strings.Split(file, "/")[len(strings.Split(file, "/"))-1], line, l.prefix, msg)
}
func (l *Logger) Error(msg string) {
	// Append color
	// Red 31
	// Green 32
	// Yellow 33
	l.mu.Lock()
	defer l.mu.Unlock()

	l.mu.Unlock()
	var ok bool
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	l.mu.Lock()
	fmt.Printf("\x1b[31m:: %s :: %s:%d :: %s :: %s\x1b\n", time.Now().Format("2006-01-02 15:04:05.000"), strings.Split(file, "/")[len(strings.Split(file, "/"))-1], line, l.prefix, msg)
}
func (l *Logger) Success(msg string) {
	// Append color
	// Red 31
	// Green 32
	// Yellow 33
	l.mu.Lock()
	defer l.mu.Unlock()

	l.mu.Unlock()
	var ok bool
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	l.mu.Lock()
	fmt.Printf("\x1b[32m:: %s :: %s:%d :: %s :: %s\x1b\n", time.Now().Format("2006-01-02 15:04:05.000"), strings.Split(file, "/")[len(strings.Split(file, "/"))-1], line, l.prefix, msg)
}
