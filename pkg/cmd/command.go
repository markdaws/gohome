package cmd

import "fmt"

type Command interface {
	GetID() string
	FriendlyString() string
	fmt.Stringer
}
