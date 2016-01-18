package gohome

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/markdaws/gohome/comm"
)

type Tcp600gwbDevice struct {
	device

	/* //TODO: Need to export to config file
	Host  string
	Token string*/
}

func (d *Tcp600gwbDevice) ModelNumber() string {
	return "TCP600GWB"
}

func (d *Tcp600gwbDevice) InitConnections() {
}

func (d *Tcp600gwbDevice) StartProducingEvents() (<-chan Event, <-chan bool) {
	return nil, nil
}

func (d *Tcp600gwbDevice) Authenticate(c comm.Connection) error {
	return nil
}

//TODO:
func (d *Tcp600gwbDevice) ZoneSetLevel(z *Zone, level float32) error {

	sendLevel := func(level int32) error {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		output := int32(level)

		// TODO: Move into connection info, user configurable
		token := "79tz3vbbop9pu5fcen60p97ix3mbvd3sblhjmz21"
		host := "https://192.168.0.23"

		var data string
		if output == 0 || output == 1 {
			data = "<gip><version>1</version><token>%s</token><did>%s</did><value>%d</value></gip>"
		} else {
			data = "<gip><version>1</version><token>%s</token><did>%s</did><value>%d</value><type>level</type></gip>"
		}
		data = fmt.Sprintf(data, token, z.LocalID, output)

		client := &http.Client{Transport: tr}
		slc := fmt.Sprintf("cmd=GWRBatch&data=<gwrcmds><gwrcmd><gcmd>DeviceSendCommand</gcmd><gdata>%s</gdata></gwrcmd></gwrcmds>&fmt=xml", data)
		fmt.Println(slc)
		_, err := client.Post(host+"/gwr/gpo.php", "text/xml; charset=\"utf-8\"", bytes.NewReader([]byte(slc)))
		return err
	}

	exec := func() error {
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
	}

	d.cmdProcessor.Enqueue(&FuncCommand{
		Func: exec,
		//TODO:
		Friendly:    "Some friendly string",
		CommandType: CTZoneSetLevel,
	})
	return nil
}

func (d *Tcp600gwbDevice) Enqueue(c Command) error {
	return fmt.Errorf("//TODO: unsupported tcp600gwbdevice")
}

func (d *Tcp600gwbDevice) BuildCommand(c Command) (*FCommand, error) {
	return nil, fmt.Errorf("//TODO: unsupported tcp600gwbdevice")
}
