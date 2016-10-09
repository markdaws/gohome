package gohome

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/event"
	"github.com/markdaws/gohome/validation"
)

type Tcp600gwbDevice struct {
	device
}

func (d *Tcp600gwbDevice) ModelNumber() string {
	return "TCP600GWB"
}

func (d *Tcp600gwbDevice) InitConnections() {
}

func (d *Tcp600gwbDevice) StartProducingEvents() (<-chan event.Event, <-chan bool) {
	return nil, nil
}

func (d *Tcp600gwbDevice) Authenticate(c comm.Connection) error {
	return nil
}

func (d *Tcp600gwbDevice) Connect() (comm.Connection, error) {
	return nil, fmt.Errorf("unsupported function connect")
}

func (d *Tcp600gwbDevice) ReleaseConnection(c comm.Connection) {
}

func (d *Tcp600gwbDevice) BuildCommand(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneSetLevel:
		return d.buildZoneSetLevelCommand(command)
	default:
		return nil, fmt.Errorf("unsupported command tcp600gwbdevice")
	}
}

func (d *Tcp600gwbDevice) buildZoneSetLevelCommand(c *cmd.ZoneSetLevel) (*cmd.Func, error) {

	//TODO: Use connected.go library
	sendLevel := func(level int32) error {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		output := int32(c.Level.Value)

		var data string
		if output == 0 || output == 1 {
			data = "<gip><version>1</version><token>%s</token><did>%s</did><value>%d</value></gip>"
		} else {
			data = "<gip><version>1</version><token>%s</token><did>%s</did><value>%d</value><type>level</type></gip>"
		}
		data = fmt.Sprintf(data, d.Auth().Token, c.ZoneAddress, output)

		client := &http.Client{Transport: tr}
		slc := fmt.Sprintf("cmd=GWRBatch&data=<gwrcmds><gwrcmd><gcmd>DeviceSendCommand</gcmd><gdata>%s</gdata></gwrcmd></gwrcmds>&fmt=xml", data)
		fmt.Println(slc)
		_, err := client.Post(d.Address()+"/gwr/gpo.php", "text/xml; charset=\"utf-8\"", bytes.NewReader([]byte(slc)))
		return err
	}

	return &cmd.Func{
		Func: func() error {
			level := c.Level.Value
			if level != 0 {
				// 0 -> off, 1 -> on, if the light was set to 0 then you have to set a 1 first
				// before trying to set any other level
				if err := sendLevel(1); err != nil {
					return err
				}
				if err := sendLevel(int32(level)); err != nil {
					return err
				}
				return nil
			} else {
				return sendLevel(0)
			}
		},
	}, nil
}

func (d *Tcp600gwbDevice) Validate() *validation.Errors {
	errors := d.device.Validate()
	if errors == nil {
		errors = &validation.Errors{}
	}

	if d.Address() == "" {
		errors.Add("required field", "Address")
	}

	if d.Auth() == nil {
		errors.Add("required field", "Token")
	} else {
		if d.Auth().Token == "" {
			errors.Add("required field", "Token")
		}
	}

	if errors.Has() {
		return errors
	}
	return nil
}
