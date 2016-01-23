package log

import (
	"fmt"
	"log"
	"os"
)

var (
	verbose *log.Logger
	error   *log.Logger
)

func init() {
	verbose = log.New(os.Stdout, "V: ", log.Ldate|log.Ltime)
	error = log.New(os.Stderr, "E: ", log.Ldate|log.Ltime)
}

// V logs a verbose message to the app log
func V(m string, args ...interface{}) {
	verbose.Printf("%s\n", fmt.Sprintf(m, args...))
}

// E logs an error message to the app log
func E(m string, args ...interface{}) {
	error.Printf("%s\n", fmt.Sprintf(m, args...))
}
