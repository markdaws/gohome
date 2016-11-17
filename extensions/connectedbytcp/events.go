package connectedbytcp

import (
	"context"
	"fmt"
	"time"

	connectedExt "github.com/go-home-iot/connectedbytcp"
	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
	errExt "github.com/pkg/errors"
)

type consumer struct {
	Name   string
	System *gohome.System
	Device *gohome.Device
}

func (c *consumer) ConsumerName() string {
	return fmt.Sprintf("ConnectedByTCPEventConsumer - %s", c.Name)
}

func (c *consumer) StartConsuming(ch chan evtbus.Event) {
	go func() {
		for e := range ch {

			var evt *gohome.ZonesReportEvt
			var ok bool
			if evt, ok = e.(*gohome.ZonesReportEvt); !ok {
				continue
			}

			// Find all of the zones in the report that we own, if non we can skip
			zones := c.Device.OwnedZones(evt.ZoneIDs)
			if len(zones) == 0 {
				continue
			}

			// For each zone the device owns, get the current address
			zoneValueByAddress, err := getZoneValuesByAddress(c.Device)
			if err != nil {
				log.V(err.Error())
				continue
			}

			for _, zone := range zones {
				log.V("%s - %s", c.ConsumerName(), evt)

				value, ok := zoneValueByAddress[zone.Address]
				if !ok {
					continue
				}

				c.System.Services.EvtBus.Enqueue(&gohome.ZoneLevelReportingEvt{
					ZoneName: zone.Name,
					ZoneID:   zone.ID,
					Level: cmd.Level{
						Value: value,
					},
				})
			}
		}
	}()
}

func (c *consumer) StopConsuming() {
	//TODO:
}

type producer struct {
	producing bool
	Name      string
	Device    *gohome.Device
	System    *gohome.System
}

func (p *producer) ProducerName() string {
	return fmt.Sprintf("ConnectedByTCPEventProducer - %s", p.Name)
}
func (p *producer) StartProducing(b *evtbus.Bus) {
	p.producing = true

	go func() {
		// Since we don't have any mechanism to automatically get updated from the
		// bulb, we just poll every 10 seconds to get the latest state
		for p.producing {
			time.Sleep(time.Second * 10)

			zoneValueByAddress, err := getZoneValuesByAddress(p.Device)
			if err != nil {
				log.V(err.Error())
				continue
			}

			for _, zone := range p.Device.Zones {
				value, ok := zoneValueByAddress[zone.Address]
				if !ok {
					continue
				}

				p.System.Services.EvtBus.Enqueue(&gohome.ZoneLevelReportingEvt{
					ZoneName: zone.Name,
					ZoneID:   zone.ID,
					Level: cmd.Level{
						Value: value,
					},
				})
			}
		}
	}()
}
func (p *producer) StopProducing() {
	p.producing = false
}

func getZoneValuesByAddress(d *gohome.Device) (map[string]float32, error) {
	// Have to get the room report, this gets the state of all zones we own
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	resp, err := connectedExt.RoomGetCarousel(ctx, d.Address, d.Auth.Token)
	if err != nil {
		return nil, errExt.Wrap(err, "failed to get state")
	}

	var zoneValueByAddress = make(map[string]float32)
	for _, room := range resp.Rooms {
		for _, device := range room.Devices {
			//fmt.Println(device.DID, " state:", device.State, ", lvl:", device.Level)
			if device.State == 0 {
				zoneValueByAddress[device.DID] = 0
				continue
			}

			zoneValueByAddress[device.DID] = float32(device.Level)
		}
	}
	return zoneValueByAddress, nil
}
