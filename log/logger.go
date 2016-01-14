package log

import (
	"fmt"
	"time"
)

func V(m string, args ...interface{}) {
	fmt.Printf("%s::%s::%s\n", "V", time.Now().UTC(), fmt.Sprintf(m, args...))
}

func W(m string, args ...interface{}) {
	fmt.Printf("%s::%s::%s\n", "W", time.Now().UTC(), fmt.Sprintf(m, args...))
}

func E(m string, args ...interface{}) {
	fmt.Printf("%s::%s::%s\n", "E", time.Now().UTC(), fmt.Sprintf(m, args...))
}
