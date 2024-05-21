package entity

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Console logger

type ConsoleLogger struct {
	mu     sync.Mutex
	out    *io.Writer
	prefix string
}

func NewConsoleLogger(out io.Writer) Logger {
	return &ConsoleLogger{
		out: &out,
	}
}

func (l *ConsoleLogger) Prefix() string {
	return l.prefix
}

func (l *ConsoleLogger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *ConsoleLogger) Info(msg string) {
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
func (l *ConsoleLogger) Warning(msg string) {
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
func (l *ConsoleLogger) Error(msg string) {
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
func (l *ConsoleLogger) Debug(msg string) {
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
