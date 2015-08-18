package main

import (
	"fmt"
	"time"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/www"
)

func main() {
	fmt.Println("hi")

	//	serverDone := make(chan bool)
	go func() {
		s := www.NewServer("./www")
		err := s.ListenAndServe(":8000")
		if err != nil {
			fmt.Println("error with server")
		}
		//		close(serverDone)
	}()

	r := &gohome.Recipe{
		Id:          "123",
		Name:        "Test",
		Description: "Test desc",
		Trigger: &gohome.TimeTrigger{
			Iterations: 5,
			Forever:    true,
			Interval:   time.Second * 10,
			At:         time.Now(),
		},
		Action: &gohome.PrintAction{},
	}
	doneChan := r.Start()

	go func() {
		time.Sleep(time.Second * 10)
		fmt.Println("stopping")
		r.Stop()
	}()

	//What is the lifetime of a recipe? How to know when done?
	<-doneChan

	//	<-serverDone
}
