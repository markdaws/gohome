package gohome

import (
	"fmt"
	"strings"

	"github.com/markdaws/gohome/log"
)

type StringCommand struct {
	Value  string
	Device Device
	Args   []interface{}
}

func (c *StringCommand) Execute() error {
	str := fmt.Sprintf(c.Value, c.Args...)
	log.V("sending command \"%s\" to Device \"%s\"", strings.Replace(strings.Replace(str, "\r", "\\r", -1), "\n", "\\n", -1), c.Device.Name())

	conn, err := c.Device.Connect()
	if err != nil {
		return fmt.Errorf("StringCommand - error connecting %s", err)
	}

	defer func() {
		c.Device.ReleaseConnection(conn)
	}()
	_, err = conn.Write([]byte(str))
	if err != nil {
		fmt.Printf("Failed to string string_command %s\n", err)
	}
	return err
}
