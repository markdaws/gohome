package lutron

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

type eventProducer struct {
	Name   string
	System *gohome.System
	Device *gohome.Device
}

func (p *eventProducer) ProducerName() string {
	return "LutronEventProducer"
}

func (p *eventProducer) StartProducing(b *evtbus.Bus) {

	//TODO: These producers shouldn't block the bus, make bus more tolerant
	go func() {
		for {
			log.V("%s attemping to stream events", p.Device)
			conn, err := p.Device.Connections.Get(time.Second * 20)
			if err != nil {
				log.V("%s unable to connect to stream events: %s", p.Device, err)
				continue
			}

			defer func() {
				p.Device.Connections.Release(conn)
			}()

			log.V("%s streaming events", p.Device)
			scanner := bufio.NewScanner(conn)
			split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

				//Match first instance of ~OUTPUT|~DEVICE.*\r\n
				str := string(data[0:])
				indices := regexp.MustCompile("[~|#][OUTPUT|DEVICE].+\r\n").FindStringIndex(str)

				//TODO: Don't let input grow forever - remove beginning chars after reaching max length

				if indices != nil {
					token = []byte(string([]rune(str)[indices[0]:indices[1]]))
					advance = indices[1]
					err = nil
				} else {
					advance = 0
					token = nil
					err = nil
				}
				return
			}

			scanner.Split(split)
			for scanner.Scan() {
				orig := scanner.Text()
				if evt := p.parseCommandString(orig); evt != nil {
					p.System.EvtBus.Enqueue(evt)
				}
			}

			log.V("%s stopped streaming events", p.Device)
			if err := scanner.Err(); err != nil {
				log.V("%s error streaming events, streaming stopped: %s", p.Device, err)
			}
		}
	}()
	//p.System.EvtBus.Enqueue
}

func (p *eventProducer) StopProducing() {
	//TODO:
}

func (p *eventProducer) parseCommandString(cmd string) evtbus.Event {
	switch {
	case strings.HasPrefix(cmd, "~OUTPUT"),
		strings.HasPrefix(cmd, "#OUTPUT"):
		return p.parseZoneCommand(cmd)

	case strings.HasPrefix(cmd, "~DEVICE"),
		strings.HasPrefix(cmd, "#DEVICE"):
		//TODO:
		//return p.parseDeviceCommand(cmd)
		return nil
	default:
		// Ignore commands we don't care about
		return nil
	}
}

func (p *eventProducer) parseDeviceCommand(command string) evtbus.Event {
	//TODO:
	/*
		matches := regexp.MustCompile("[~|#]DEVICE,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(command)
		if matches == nil || len(matches) != 4 {
			return nil
		}

		deviceID := matches[1]
		componentID := matches[2]
		cmdID := matches[3]
		sourceDevice := p.Device.Devices[deviceID]
		if sourceDevice == nil {
			return nil
		}

		var finalCmd cmd.Command
		switch cmdID {
		case "3":
			if btn := sourceDevice.Buttons()[componentID]; btn != nil {
				finalCmd = &cmd.ButtonPress{
					ButtonAddress: btn.Address,
					ButtonID:      btn.ID,
					DeviceName:    d.Name(),
					DeviceAddress: d.Address(),
					DeviceID:      d.ID(),
				}
			}
		case "4":
			if btn := sourceDevice.Buttons()[componentID]; btn != nil {
				finalCmd = &cmd.ButtonRelease{
					ButtonAddress: btn.Address,
					ButtonID:      btn.ID,
					DeviceName:    d.Name(),
					DeviceAddress: d.Address(),
					DeviceID:      d.ID(),
				}
			}
		default:
			return nil
		}

		return finalCmd*/
	return nil
}

func (p *eventProducer) parseZoneCommand(command string) evtbus.Event {
	matches := regexp.MustCompile("[~|?]OUTPUT,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(command)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	zoneID := matches[1]
	cmdID := matches[2]
	level, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return nil
	}

	z := p.Device.Zones[zoneID]
	if z == nil {
		return nil
	}

	var finalCmd cmd.Command
	switch cmdID {
	case "1":
		return &gohome.ZoneLevelChanged{
			ZoneName: z.Name,
			ZoneID:   z.ID,
			Level:    cmd.Level{Value: float32(level)},
		}
	default:
		return nil
	}

	return finalCmd
}
