package logx

import (
	"fmt"
	"io"
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
