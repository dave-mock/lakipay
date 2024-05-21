package entity

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warning(msg string)
	Error(msg string)
}

// Console logger

// File logger

// Web logger
