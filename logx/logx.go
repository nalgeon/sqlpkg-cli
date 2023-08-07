// Package logx provides a logging utility.
package logx

import (
	"io"
	"os"
)

var logger = NewLogger(os.Stdout)

// IsVerbose returns true if the logger is verbose.
func IsVerbose() bool {
	return logger.IsVerbose
}

// SetVerbose changes the verboseness of the logger.
func SetVerbose(val bool) {
	logger.IsVerbose = val
}

// Output returns the logger destination.
func Output() io.Writer {
	return logger.out
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

// Mock creates a new Memory and installs it as the logger output
// instead of the default one. Should be used for testing purposes only.
func Mock(path ...string) *Memory {
	memory := NewMemory("log")
	SetOutput(memory)
	SetVerbose(true)
	return memory
}
