package gohome

import "fmt"

type Command interface {
	Execute() error
	FriendlyString() string
	CMDType() CommandType
	fmt.Stringer
}
