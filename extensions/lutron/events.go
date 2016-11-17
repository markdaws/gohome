package lutron

import (
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

func (c *eventConsumer) ConsumerName() string {
	return "LutronEventConsumer"
}
func (c *eventConsumer) StartConsuming(ch chan evtbus.Event) {
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
			dev, err := lutronExt.DeviceFromModelNumber(c.Device.ModelNumber)
			if err != nil {
				log.V("%s - error, unsupported device %s inside consumer", c.ConsumerName(), c.Device.ModelNumber)
				continue
			}

			log.V("%s - %s", c.ConsumerName(), zoneRpt)

			for _, zone := range c.Device.OwnedZones(zoneRpt.ZoneIDs) {
				conn, err := c.Device.Connections.Get(time.Second*10, true)
				if err != nil {
					log.V("%s - unable to get connection to device: %s, timeout", c.ConsumerName(), c.Device)
					continue
				}

				err = dev.RequestLevel(zone.Address, conn)
				c.Device.Connections.Release(conn, err)
				if err != nil {
					log.V("%s - Failed to request level for lutron, zoneID:%s, %s", c.ConsumerName(), zone.ID, err)
				}
			}
		}
	}()
}
func (c *eventConsumer) StopConsuming() {
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
			// If we have broken out the loop, have a small wait period so we don't end up
			// in a tight loop if the error keep occuring
			time.Sleep(time.Second * 5)

			log.V("%s attempting to stream events", p.Device)

			conn, err := p.Device.Connections.Get(time.Second*20, true)
			if err != nil {
				log.V("%s unable to connect to stream events: %s", p.Device, err)
				continue
			}

			dev, err := lutronExt.DeviceFromModelNumber(p.Device.ModelNumber)
			if err != nil {
				log.V("unable to get lutron device for model number %s", p.Device.ModelNumber)
				continue
			}

			log.V("%s streaming events", p.Device)

			// Let the system know we are ready to process events
			b.Enqueue(&gohome.DeviceProducingEvt{
				Device: p.Device,
			})

			err = dev.Stream(conn, func(evt lutronExt.Event) {
				if !p.producing {
					return
				}

				switch e := evt.(type) {
				case *lutronExt.ZoneLevelEvt:
					z, ok := p.Device.Zones[e.Address]
					if !ok {
						return
					}
					p.System.Services.EvtBus.Enqueue(&gohome.ZoneLevelReportingEvt{
						ZoneName: z.Name,
						ZoneID:   z.ID,
						Level:    cmd.Level{Value: e.Level},
					})

				case *lutronExt.BtnPressEvt:
					sourceDevice, ok := p.Device.Devices[e.DeviceAddress]
					if !ok {
						return
					}

					if btn := sourceDevice.Buttons[e.Address]; btn != nil {
						p.System.Services.EvtBus.Enqueue(&gohome.ButtonPressEvt{
							BtnAddress:    btn.Address,
							BtnID:         btn.ID,
							BtnName:       btn.Name,
							DeviceName:    sourceDevice.Name,
							DeviceAddress: sourceDevice.Address,
							DeviceID:      sourceDevice.ID,
						})
					}

				case *lutronExt.BtnReleaseEvt:
					sourceDevice, ok := p.Device.Devices[e.DeviceAddress]
					if !ok {
						return
					}

					if btn := sourceDevice.Buttons[e.Address]; btn != nil {
						p.System.Services.EvtBus.Enqueue(&gohome.ButtonReleaseEvt{
							BtnAddress:    btn.Address,
							BtnID:         btn.ID,
							BtnName:       btn.Name,
							DeviceName:    sourceDevice.Name,
							DeviceAddress: sourceDevice.Address,
							DeviceID:      sourceDevice.ID,
						})
					}

				case *lutronExt.UnknownEvt:
					// Don't care, ignore
				}
			})

			log.V("%s stopped streaming events", p.Device)
			p.Device.Connections.Release(conn, err)

			if err != nil {
				log.V("%s error streaming events, streaming stopped: %s", p.Device, err)
			}
		}
	}()
}

func (p *eventProducer) StopProducing() {
	p.producing = false
	//TODO: Stop the scanner
}
