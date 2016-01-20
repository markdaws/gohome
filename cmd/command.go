package cmd

import "fmt"

type Command interface {
	FriendlyString() string
	fmt.Stringer
}
