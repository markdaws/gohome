package fluxwifi

import (
	"fmt"
	"time"

	"github.com/go-home-iot/event-bus"
	fluxwifiExt "github.com/go-home-iot/fluxwifi"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/feature"
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

			var evt *gohome.FeaturesReportEvt
			var ok bool
			if evt, ok = e.(*gohome.FeaturesReportEvt); !ok {
				continue
			}

			for _, f := range c.Device.OwnedFeatures(evt.FeatureIDs) {
				// All features are LightZone type, no need to filter between different types

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
					continue
				}

				if state.Power < 2 {
					// Get the cloned attributes
					onoff, _, hsl := feature.LightZoneCloneAttrs(f)

					var onoffVal int32
					if state.Power > 0 {
						onoffVal = attr.OnOffOn
					} else {
						onoffVal = attr.OnOffOff
					}
					onoff.Value = onoffVal
					hsl.Value = attr.RGBToHSLString(int(state.R), int(state.G), int(state.B))

					c.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(onoff, hsl),
					})
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

			for _, f := range p.Device.Features {
				// We only support FTLightZone features, so no need to type check

				conn, err := p.Device.Connections.Get(time.Second*10, false)
				if err != nil {
					log.V("%s - failed to get connection to check status: %s", p.ProducerName(), err)
					continue
				}

				state, err := fluxwifiExt.FetchState(conn)
				p.Device.Connections.Release(conn, err)
				if err != nil {
					log.V("%s - failed to get bulb state: %s", p.ProducerName(), err)
					continue
				}

				// 2 is unknown so ignore
				if state.Power < 2 {
					onoff, _, hsl := feature.LightZoneCloneAttrs(f)

					if state.Power > 0 {
						onoff.Value = attr.OnOffOn
					} else {
						onoff.Value = attr.OnOffOff
					}

					hsl.Value = attr.RGBToHSLString(int(state.R), int(state.G), int(state.B))

					p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(onoff, hsl),
					})
				}
			}
		}
	}()
}
func (p *producer) StopProducing() {
	p.producing = false
}
