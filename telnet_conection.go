package gohome

import (
	"bufio"
	"fmt"
	"net"
)

type TelnetConnection struct {
	Login    string
	Password string
	Network  string
	Address  string
	conn     net.Conn
}

func (c *TelnetConnection) Connect() error {
	fmt.Println("trying to connect")
	conn, err := net.Dial(c.Network, c.Address)
	if err != nil {
		fmt.Printf("Dial failed\n")
		return err
	}

	r := bufio.NewReader(conn)
	_, err = r.ReadString(':')
	if err != nil {
		fmt.Println("Failed to read login", err)
		return err
	}
	fmt.Println("Got past login")
	c.conn = conn
	_, err = conn.Write([]byte(c.Login + "\r\n"))
	if err != nil {
		fmt.Println("Failed to write password", err)
		return err
	}
	fmt.Println("Wrote login")
	_, err = r.ReadString(':')
	if err != nil {
		fmt.Println("error waiting for password", err)
		return err
	}
	_, err = conn.Write([]byte(c.Password + "\r\n"))
	if err != nil {
		fmt.Println("Error writing password")
		return err
	}
	fmt.Println("wrote password")

	//	conn.Write([]byte("#DEVICE,1,4,3\r\n"))

	//	time.Sleep(time.Second * 3)
	return nil
}

func (c *TelnetConnection) Send(data []byte) {
	c.conn.Write(data)
}
