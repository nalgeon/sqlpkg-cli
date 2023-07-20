// Logging.
package cmd

import "fmt"

var IsVerbose bool

// Log prints a message to the screen.
func Log(message string, args ...any) {
	if len(args) == 0 {
		fmt.Println(message)
	} else {
		fmt.Printf(message+"\n", args...)
	}
}

// Debug prints a message to the screen if the verbose mode is on.
func Debug(message string, args ...any) {
	if !IsVerbose {
		return
	}
	Log(".."+message, args...)
}
