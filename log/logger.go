package log

import (
	"fmt"
	"log"
	"os"
)

var (
	verbose *log.Logger
	warning *log.Logger
	error   *log.Logger
)

func init() {
	verbose = log.New(os.Stdout, "V: ", log.Ldate|log.Ltime)
	warning = log.New(os.Stdout, "W: ", log.Ldate|log.Ltime)
	error = log.New(os.Stderr, "E: ", log.Ldate|log.Ltime)
}

func V(m string, args ...interface{}) {
	verbose.Printf("%s\n", fmt.Sprintf(m, args...))
}

func W(m string, args ...interface{}) {
	warning.Printf("%s\n", fmt.Sprintf(m, args...))
}

func E(m string, args ...interface{}) {
	error.Printf("%s\n", fmt.Sprintf(m, args...))
}
