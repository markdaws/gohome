package gohome

import "fmt"

type Identifiable struct {
	ID          string
	Name        string
	Description string
}

func (i Identifiable) String() string {
	return fmt.Sprintf("ID:\"%s\",Name:%s,Description:\"%s\"", i.ID, i.Name, i.Description)
}
