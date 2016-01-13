package log

import (
	"fmt"
	"time"
)

func V(m string, args ...interface{}) {
	fmt.Printf("%s::%s::%s", "I", time.Now().UTC(), fmt.Sprintf(m, args...))
}

func W(m string, args ...interface{}) {
	fmt.Printf("%s::%s::%s", "W", time.Now().UTC(), fmt.Sprintf(m, args...))
}

func E(m string, args ...interface{}) {
	fmt.Printf("%s::%s::%s", "E", time.Now().UTC(), fmt.Sprintf(m, args...))
}
