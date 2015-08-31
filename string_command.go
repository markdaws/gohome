package gohome

import "fmt"

type StringCommand struct {
	Value    string
	Friendly string
	Device   *Device
	Type     CommandType
}

func (c *StringCommand) Execute(args ...interface{}) error {
	str := fmt.Sprintf(c.Value, args...)
	fmt.Println("Setting command:", str)

	conn, err := c.Device.Connect()
	if err != nil {
		return err
	}

	defer func() {
		conn.Close()
	}()

	//TODO: If n < data, keep going
	_, err = conn.Write([]byte(str))
	return err
}

func (c *StringCommand) String() string {
	return c.Value
}

func (c *StringCommand) FriendlyString() string {
	return c.Friendly
}

func (c *StringCommand) GetType() CommandType {
	return c.Type
}
