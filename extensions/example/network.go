package example

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/log"
)

type network struct {
	Device *gohome.Device
}

func (d *network) NewConnection(sys *gohome.System, dev *gohome.Device) (func(pool.Config) (net.Conn, error), error) {
	return func(cfg pool.Config) (net.Conn, error) {

		// Here you need to create an open connection to a device e.g. call net.DialTimeout, and perform
		// any authentication you need before returning the connection back to the caller.

		// For example. when we connect to our fake example hardware it shows:
		// login:
		// We need to wait for that text, then write the login, if then shows:
		// password:
		// We need to write the password, only at this point can we pass the connection back.  Your
		// device might have different protocols that you have to implement here.

		// In your hardware maybe you dont need to do any authentication, in which case you can
		// just do a net.DialTimeout and return.

		// NOTE: Make sure you check all error values and also use Timeouts to make sure this
		// can't block forever.

		conn, err := net.DialTimeout("tcp", dev.Address, time.Second*10)
		if err != nil {
			log.E("Failed to connect to Device[%s] %s, %s", dev.Name, dev.Address, err)
			return nil, err
		}

		// Wait for login: to appear on the stream
		r := bufio.NewReader(conn)
		_, err = r.ReadString(':')
		if err != nil {
			return nil, fmt.Errorf("authenticate login failed: %s", err)
		}

		// Write the login information
		_, err = conn.Write([]byte(dev.Auth.Login + "\r\n"))
		if err != nil {
			return nil, fmt.Errorf("authenticate write login failed: %s", err)
		}

		// Wait for password: to appear on the stream
		_, err = r.ReadString(':')
		if err != nil {
			return nil, fmt.Errorf("authenticate password failed: %s", err)
		}

		// Write the password
		_, err = conn.Write([]byte(dev.Auth.Password + "\r\n"))
		if err != nil {
			return nil, fmt.Errorf("authenticate password failed: %s", err)
		}
		return conn, nil
	}, nil
}
