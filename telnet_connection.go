package gohome

import (
	"bufio"
	"fmt"
	"net"
)

type TelnetConnection struct {
	conn net.Conn
	info ConnectionInfo
}

func NewTelnetConnection(i ConnectionInfo) *TelnetConnection {
	c := TelnetConnection{}
	c.SetConnectionInfo(i)
	return &c
}

func (c *TelnetConnection) SetConnectionInfo(i ConnectionInfo) {
	c.info = i
}

func (c *TelnetConnection) Open() error {
	//TODO: log properly
	fmt.Println("trying to connect")
	conn, err := net.Dial(c.info.Network, c.info.Address)
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
	_, err = conn.Write([]byte(c.info.Login + "\r\n"))
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
	_, err = conn.Write([]byte(c.info.Password + "\r\n"))
	if err != nil {
		fmt.Println("Error writing password")
		return err
	}
	fmt.Println("wrote password")
	return nil
}

func (c *TelnetConnection) Close() error {
	//TODO: return to pool
	return nil
}

func (c *TelnetConnection) Read(p []byte) (n int, err error) {
	n, err = c.conn.Read(p)
	return
}

func (c *TelnetConnection) Write(p []byte) (n int, err error) {
	n, err = c.conn.Write(p)
	return
}
