package logger

import "github.com/kelchy/go-lib/log"

// Logger is a simple interface for logging.
type Logger interface {
	Debug(key string, message string)
	Out(key string, message string)
	Error(key string, err error)
}

var DefaultLogger, _ = log.New("standard")
