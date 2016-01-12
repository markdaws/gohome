package gohome

import "fmt"

type Command interface {
	Execute(args ...interface{}) error
	FriendlyString() string
	CMDType() CommandType
	fmt.Stringer
}
