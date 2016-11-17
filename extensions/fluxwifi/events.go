package fluxwifi

import (
	"fmt"
	"time"

	"github.com/go-home-iot/event-bus"
	fluxwifiExt "github.com/go-home-iot/fluxwifi"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

type consumer struct {
	Name   string
	System *gohome.System
	Device *gohome.Device
}

func (c *consumer) ConsumerName() string {
	return fmt.Sprintf("FluxwifiEventConsumer - %s", c.Name)
}

func (c *consumer) StartConsuming(ch chan evtbus.Event) {
	go func() {
		for e := range ch {

			var evt *gohome.ZonesReportEvt
			var ok bool
			if evt, ok = e.(*gohome.ZonesReportEvt); !ok {
				continue
			}

			for _, zone := range c.Device.OwnedZones(evt.ZoneIDs) {
				log.V("%s - %s", c.ConsumerName(), evt)

				conn, err := c.Device.Connections.Get(time.Second*5, true)
				if err != nil {
					log.V("%s - failed to get connection: %s", c.ConsumerName(), err)
					continue
				}

				state, err := fluxwifiExt.FetchState(conn)
				c.Device.Connections.Release(conn, err)
				if err != nil {
					log.V("%s - failed to get state: %s", c.ConsumerName(), err)
				} else {
					if state.Power < 2 {
						// 2 -> unknown, so only process if it is 0 or 1
						c.System.Services.EvtBus.Enqueue(&gohome.ZoneLevelReportingEvt{
							ZoneName: zone.Name,
							ZoneID:   zone.ID,
							Level: cmd.Level{
								Value: float32(state.Power),
								R:     state.R,
								G:     state.G,
								B:     state.B,
							},
						})
					}
				}

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
	return fmt.Sprintf("FluxwifiEventProducer - %s", p.Name)
}
func (p *producer) StartProducing(b *evtbus.Bus) {
	p.producing = true

	go func() {
		// Since we don't have any mechanism to automatically get updated from the
		// bulb, we just poll every 10 seconds to get the latest state
		for p.producing {
			time.Sleep(time.Second * 10)

			for _, zone := range p.Device.Zones {
				conn, err := p.Device.Connections.Get(time.Second*10, false)
				if err != nil {
					log.V("%s - failed to get connection to check status: %s", p.ProducerName(), err)
					continue
				}

				state, err := fluxwifiExt.FetchState(conn)
				p.Device.Connections.Release(conn, err)
				if err != nil {
					log.V("%s - failed to get bulb state: %s", p.ProducerName(), err)
				} else {
					if state.Power < 2 {
						// 2 -> unknown, so only process if it is 0 or 1
						p.System.Services.EvtBus.Enqueue(&gohome.ZoneLevelReportingEvt{
							ZoneName: zone.Name,
							ZoneID:   zone.ID,
							Level: cmd.Level{
								Value: float32(state.Power),
								R:     state.R,
								G:     state.G,
								B:     state.B,
							},
						})
					}
				}

			}
		}
	}()
}
func (p *producer) StopProducing() {
	p.producing = false
}
