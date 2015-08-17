package main

import (
	"fmt"

	"github.com/markdaws/gohome/www"
)

func main() {
	fmt.Println("hi")

	s := www.NewServer("./www")
	s.ListenAndServe(":8000")
}
