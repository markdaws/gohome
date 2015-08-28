package gohome

import "fmt"

type Identifiable struct {
	//TODO: Rename ID
	Id          string
	Name        string
	Description string
}

func (i Identifiable) String() string {
	return fmt.Sprintf("Id:\"%s\",Name:%s,Description:\"%s\"", i.Id, i.Name, i.Description)
}
