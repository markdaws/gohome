package gohome

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Device struct {
	Identifiable
	System     *System
	Connection Connection
	evpDone    chan bool
	evpFire    chan Event
}

func (d *Device) Connect() error {
	conn, err := d.Connection.Connect()
	if err != nil {
		return err
	}

	//TODO: This should be a connection pool ...
	//TODO: Should be an option to persist connection
	go func() {
		stream(d, conn)
	}()

	return nil
}

func (d *Device) GetEventProducerChans() (<-chan Event, <-chan bool) {
	d.evpDone = make(chan bool)
	d.evpFire = make(chan Event)
	return d.evpFire, d.evpDone
}

func stream(d *Device, r io.Reader) {
	scanner := bufio.NewScanner(r)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		//Remove leading GET>
		//Remove leading spaces
		//find first instance of \r\n return if one
		//TODO: Replace with regex
		//TODO: Different devices will need different split functions ...
		str := string(data[0:])
		origLen := len(str)
		str = strings.TrimLeft(str, "GNET>")
		str = strings.TrimLeft(str, " ")
		offsetTotal := origLen - len(str)
		index := strings.Index(str, "\r\n")

		//TODO: Don't let input grow forever - remove beginning chars after reaching max length
		//TODO: Save raw stream
		if index != -1 {
			// Ignore lines with just \r\n
			if index == 0 {
				token = nil
			} else {
				token = []byte(string([]rune(str)[0 : index+2]))
			}
			advance = index + 2 + offsetTotal
			err = nil
		} else {
			advance = 0
			token = nil
			err = nil
		}
		return
	}

	scanner.Split(split)
	for scanner.Scan() {
		//		fmt.Printf("scanner: %s\n", scanner.Text())

		if d.evpFire != nil {
			orig := scanner.Text()
			cmd := ParseCommandString(d, orig)
			d.evpFire <- NewEvent(d, cmd, orig)
		}
	}

	if d.evpDone != nil {
		close(d.evpDone)
	}

	//TODO: What about connect/disconnect event
	fmt.Println("Done scanning")
	if err := scanner.Err(); err != nil {
		fmt.Printf("something happened", err)
	}
}
