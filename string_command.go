package gohome

import "fmt"

type StringCommand struct {
	Value    string
	Friendly string
	Device   *Device
	Type     CommandType
	Args     []interface{}
}

func (c *StringCommand) Execute() error {
	str := fmt.Sprintf(c.Value, c.Args...)
	fmt.Println("Sending command:", str)

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

func (c *StringCommand) String() string {
	return c.Value
}

func (c *StringCommand) FriendlyString() string {
	return c.Friendly
}

func (c *StringCommand) CMDType() CommandType {
	return c.Type
}
