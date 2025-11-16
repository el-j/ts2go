package console

import (
	"fmt"
	"log"
)

// Log prints to stdout (equivalent to console.log)
func Log(args ...interface{}) {
	fmt.Println(args...)
}

// Error prints to stderr (equivalent to console.error)
func Error(args ...interface{}) {
	log.Println(args...)
}

// Warn prints a warning (equivalent to console.warn)
func Warn(args ...interface{}) {
	log.Println(args...)
}
