package gohome

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type TelnetConnection struct {
	Login    string
	Password string
	Network  string
	Address  string
	conn     net.Conn
}

func stream(c net.Conn) {
	buf := make([]byte, 4096)
	for {
		n, err := c.Read(buf)
		if err != nil || n == 0 {
			fmt.Println("connection closed")
			c.Close()
			return
		}
		str := string(buf[0:n])

		events := strings.Split(str, "\r\n")
		for _, event := range events {
			//#,~,?
			fmt.Printf("%s\n", event)
		}
	}
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

	go func() {
		stream(conn)
	}()
	return nil
}

func (c *TelnetConnection) Send(data []byte) {
	c.conn.Write(data)
}
