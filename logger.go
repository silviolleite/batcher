package batcher

import (
	"log"
	"os"
	"sync"
)

var mu sync.Mutex

// A Logger is a minimalistic interface for the batcher to log messages to. Should
// be used to provide custom logging writers for the batcher to use.
type Logger interface {
	Log(...interface{})
}

// A LoggerFunc is a convenience type to convert a function taking a variadic
// list of arguments and wrap it so the Logger interface can be used.
//
// Example:
//
//	batcher.New(&batcher.Options{Logger: batcher.LoggerFunc(func(args ...interface{}) {
//	    fmt.Fprintln(os.Stdout, args...)
//	})})
type LoggerFunc func(...interface{})

// Log calls the wrapped function with the arguments provided
func (f LoggerFunc) Log(args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	f(args...)
}

// newDefaultLogger returns a Logger which will write log messages to stdout, and
// use same formatting runes as the stdlib log.Logger
func newDefaultLogger() Logger {
	return &defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// A defaultLogger provides a minimalistic logger satisfying the Logger interface.
type defaultLogger struct {
	logger *log.Logger
}

// Log logs the parameters to the stdlib logger. See log.Println.
func (l defaultLogger) Log(args ...interface{}) {
	l.logger.Println(args...)
}
