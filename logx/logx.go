// Package logx provides a logging utility.
package logx

import (
	"fmt"
	"io"
	"os"
)

// Logger logs messages to a destination.
type Logger struct {
	IsVerbose bool
	out       io.Writer
}

// NewLogger creates a new logger.
func NewLogger(out io.Writer) *Logger {
	return &Logger{out: out}
}

// SetOutput changes the logger destination.
func (l *Logger) SetOutput(out io.Writer) {
	l.out = out
}

// Log prints a message.
func (l *Logger) Log(message string, args ...any) {
	if len(args) == 0 {
		fmt.Fprintln(l.out, message)
	} else {
		fmt.Fprintf(l.out, message+"\n", args...)
	}
}

// Debug prints a message if the verbose mode is on.
func (l *Logger) Debug(message string, args ...any) {
	if !l.IsVerbose {
		return
	}
	l.Log(".."+message, args...)
}

var logger = NewLogger(os.Stdout)

// IsVerbose returns true if the logger is verbose.
func IsVerbose(out io.Writer) bool {
	return logger.IsVerbose
}

// SetVerbose changes the verboseness of the logger.
func SetVerbose(val bool) {
	logger.IsVerbose = val
}

// SetOutput changes the logger destination.
func SetOutput(out io.Writer) {
	logger.out = out
}

// Log prints a message to the console.
func Log(message string, args ...any) {
	logger.Log(message, args...)
}

// Debug prints a message to the console if the verbose mode is on.
func Debug(message string, args ...any) {
	logger.Debug(message, args...)
}
