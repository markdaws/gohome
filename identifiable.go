package gohome

import "fmt"

type Identifiable struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (i Identifiable) String() string {
	return fmt.Sprintf("Id:\"%s\",Name:%s,Description:\"%s\"", i.Id, i.Name, i.Description)
}
