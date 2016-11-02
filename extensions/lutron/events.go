package lutron

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-home-iot/event-bus"
	lutronExt "github.com/go-home-iot/lutron"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

type eventConsumer struct {
	Name   string
	System *gohome.System
	Device *gohome.Device
}

func (p *eventConsumer) ConsumerName() string {
	return "LutronEventConsumer"
}
func (p *eventConsumer) StartConsuming(ch chan evtbus.Event) {
	go func() {
		for e := range ch {

			// If we have a backlog, merge all of the requests in to one
			zoneRpt := &gohome.ZonesReportEvt{ZoneIDs: make(map[string]bool)}
			for {
				switch evt := e.(type) {
				case *gohome.ZonesReportEvt:
					zoneRpt.Merge(evt)
				}

				if len(ch) > 0 {
					e = <-ch
				} else {
					break
				}
			}

			if len(zoneRpt.ZoneIDs) == 0 {
				continue
			}

			// The system wants zones to report their current status, check if
			// we own any of these zones, if so report them
			dev, err := lutronExt.DeviceFromModelNumber(p.Device.ModelNumber)
			if err != nil {
				log.V("%s - error, unsupported device %s inside consumer", p.ConsumerName(), p.Device.ModelNumber)
				continue
			}

			log.V("%s - %s", p.ConsumerName(), zoneRpt)

			for _, zone := range p.Device.Zones {
				if _, ok := zoneRpt.ZoneIDs[zone.ID]; ok {
					conn, err := p.Device.Connections.Get(time.Second * 10)
					if err != nil {
						log.V("%s - unable to get connection to device: %s, timeout", p.ConsumerName(), p.Device)
						continue
					}
					err = dev.RequestLevel(zone.Address, conn)
					if err != nil {
						log.V("%s - Failed to request level for lutron, zoneID:%s, %s", p.ConsumerName(), zone.ID, err)
						conn.IsBad = true
					}
					p.Device.Connections.Release(conn)
				}
			}
		}
	}()
}
func (p *eventConsumer) StopConsuming() {
	//TODO:
}

type eventProducer struct {
	Name      string
	System    *gohome.System
	Device    *gohome.Device
	producing bool
}

func (p *eventProducer) ProducerName() string {
	return "LutronEventProducer: " + p.Name
}

func (p *eventProducer) StartProducing(b *evtbus.Bus) {
	p.producing = true

	go func() {
		for p.producing {
			log.V("%s attempting to stream events", p.Device)
			conn, err := p.Device.Connections.Get(time.Second * 20)
			if err != nil {
				log.V("%s unable to connect to stream events: %s", p.Device, err)
				continue
			}

			log.V("%s streaming events", p.Device)
			time.Sleep(time.Second * 30)
			scanner := bufio.NewScanner(conn)
			split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

				//Match first instance of ~OUTPUT|~DEVICE.*\r\n
				str := string(data[0:])
				//log.V("From lutron: " + str)
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

			// Let the system know we are ready to process events
			b.Enqueue(&gohome.DeviceProducingEvt{
				Device: p.Device,
			})

			for scanner.Scan() {
				orig := scanner.Text()
				if evt := p.parseCommandString(orig); evt != nil {
					p.System.Services.EvtBus.Enqueue(evt)
				}
			}

			//TODO: There should be a ping mechanism here so we know this connection
			//is still alive...
			log.V("%s stopped streaming events", p.Device)
			conn.IsBad = true
			p.Device.Connections.Release(conn)
			if err := scanner.Err(); err != nil {
				log.V("%s error streaming events, streaming stopped: %s", p.Device, err)
			}
		}
	}()
}

func (p *eventProducer) StopProducing() {
	p.producing = false
	//TODO: get out of the event loop, stop the scanner
}

//TODO: Move all this parsing to go-home-iot/lutron
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

//TODO: Move this to go-home-iot/lutron
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
		return &gohome.ZoneLevelChangedEvt{
			ZoneName: z.Name,
			ZoneID:   z.ID,
			Level:    cmd.Level{Value: float32(level)},
		}
	default:
		return nil
	}

	return finalCmd
}
