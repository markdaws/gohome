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

func (c *TelnetConnection) Connect() (net.Conn, error) {
	fmt.Println("trying to connect")
	conn, err := net.Dial(c.Network, c.Address)
	if err != nil {
		fmt.Printf("Dial failed\n")
		return nil, err
	}

	r := bufio.NewReader(conn)
	_, err = r.ReadString(':')
	if err != nil {
		fmt.Println("Failed to read login", err)
		return nil, err
	}
	fmt.Println("Got past login")
	c.conn = conn
	_, err = conn.Write([]byte(c.Login + "\r\n"))
	if err != nil {
		fmt.Println("Failed to write password", err)
		return nil, err
	}
	fmt.Println("Wrote login")
	_, err = r.ReadString(':')
	if err != nil {
		fmt.Println("error waiting for password", err)
		return nil, err
	}
	_, err = conn.Write([]byte(c.Password + "\r\n"))
	if err != nil {
		fmt.Println("Error writing password")
		return nil, err
	}
	fmt.Println("wrote password")

	return c.conn, nil
}

func (c *TelnetConnection) Send(data []byte) {
	c.conn.Write(data)
}
