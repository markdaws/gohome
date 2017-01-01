package lutron

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-home-iot/connection-pool"
	"github.com/go-home-iot/event-bus"
	lutronExt "github.com/go-home-iot/lutron"
	"github.com/markdaws/gohome/pkg/attr"
	"github.com/markdaws/gohome/pkg/feature"
	"github.com/markdaws/gohome/pkg/gohome"
	"github.com/markdaws/gohome/pkg/log"
)

type event struct {
	Name               string
	System             *gohome.System
	Device             *gohome.Device
	producing          bool
	waitingForResponse bool
	scannerConn        *pool.Connection
	gotResponse        bool
	scannerMutex       sync.RWMutex
}

func (c *event) ConsumerName() string {
	return "LutronEvent"
}
func (c *event) StartConsuming(ch chan evtbus.Event) {
	go func() {
		dev, err := lutronExt.DeviceFromModelNumber(c.Device.ModelNumber)
		if err != nil {
			log.V("%s - error, unsupported device %s inside consumer", c.ConsumerName(), c.Device.ModelNumber)
			return
		}

		for e := range ch {
			switch evt := e.(type) {
			case *gohome.FeaturesReportEvt:
				log.V("%s - %s", c.ConsumerName(), evt)

				for _, f := range c.Device.OwnedFeatures(evt.FeatureIDs) {
					switch f.Type {
					case feature.FTLightZone:
					case feature.FTWindowTreatment:
					default:
						// Only support lights + window treatments, ignore all other features,
						// which are buttons
						continue
					}

					conn, err := c.Device.Connections.Get(time.Second*10, true)
					if err != nil {
						log.V("%s - unable to get connection to device: %s, %s", c.ConsumerName(), c.Device, err)
						continue
					}

					// Set a timeout to make sure we get a valid response
					c.responseTimeout()

					// Here we simply ask the for levels, the Lutron gateway will then stream back
					// the results at some point in time, which we will process in the producer loop
					// below, since this is all async
					err = dev.RequestLevel(f.Address, conn)
					c.Device.Connections.Release(conn, err)
					if err != nil {
						log.V("%s - Failed to request level for lutron, featureID:%s, %s", c.ConsumerName(), f.ID, err)
					}
				}
			}
		}
	}()
}
func (c *event) StopConsuming() {
	//TODO:
}

func (p *event) responseTimeout() {
	p.scannerMutex.Lock()
	if p.waitingForResponse {
		p.scannerMutex.Unlock()
		return
	}
	p.waitingForResponse = true
	p.scannerMutex.Unlock()

	p.gotResponse = false

	// When we get a request for data, we will wait and if we haven't received any response
	// by 30 seconds we will assume the lutron smart bridge connection for the scanner is lost
	// and will forcefully close the connection
	go func() {
		time.Sleep(time.Second * 30)

		p.scannerMutex.RLock()
		if !p.gotResponse {
			log.V("no response from lutron smart bridge after 30 seconds, closing connection")

			// Force close the connection
			p.scannerConn.Conn.Close()
		}
		p.scannerMutex.RUnlock()

		p.scannerMutex.Lock()
		p.waitingForResponse = false
		p.scannerMutex.Unlock()
	}()
}

func (p *event) ProducerName() string {
	return "LutronEvent: " + p.Name
}

func (p *event) StartProducing(b *evtbus.Bus) {
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
				p.Device.Connections.Release(conn, nil)
				continue
			}

			log.V("%s streaming events", p.Device)

			// Let the system know we are ready to process events
			b.Enqueue(&gohome.DeviceProducingEvt{
				Device: p.Device,
			})

			p.scannerConn = conn
			err = dev.Stream(conn, func(evt lutronExt.Event) {
				p.gotResponse = true

				if !p.producing {
					return
				}

				switch e := evt.(type) {
				case *lutronExt.ZoneLevelEvt:
					// Zones can be either lights or window treatment, can't distinguish from the
					// lutron data, so have to try both
					f := p.Device.FeatureTypeByAddress(feature.FTLightZone, e.Address)
					if f == nil {

						f = p.Device.FeatureTypeByAddress(feature.FTWindowTreatment, e.Address)
						if f == nil {
							return
						}
					}

					switch f.Type {
					case feature.FTLightZone:
						onoff, brightness, _ := feature.LightZoneCloneAttrs(f)

						// Check if the light zone is dimmable, if so it will have a brightness attribute
						if brightness != nil {
							brightness.Value = e.Level
						}

						if e.Level > 0 {
							onoff.Value = attr.OnOffOn
						} else {
							onoff.Value = attr.OnOffOff
						}

						p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
							FeatureID: f.ID,
							Attrs:     feature.NewAttrs(onoff, brightness),
						})

					case feature.FTWindowTreatment:
						openClose, offset := feature.WindowTreatmentCloneAttrs(f)
						offset.Value = e.Level
						if e.Level > 0 {
							openClose.Value = attr.OpenCloseOpen
						} else {
							openClose.Value = attr.OpenCloseClosed
						}

						p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
							FeatureID: f.ID,
							Attrs:     feature.NewAttrs(openClose, offset),
						})
					}

				case *lutronExt.BtnPressEvt:
					sourceDev := p.System.DeviceByAddress(e.DeviceAddress)
					if sourceDev == nil {
						return
					}

					btn := sourceDev.FeatureTypeByAddress(feature.FTButton, e.Address)
					if btn == nil {
						return
					}

					state := feature.ButtonCloneAttrs(btn)
					state.Value = attr.ButtonStatePressed

					// TODO: What happens here if we are firing this event, and the monitor
					// system eventually handles buttons, need to make sure we don't fire
					// two events (fame for button release)
					p.System.Services.EvtBus.Enqueue(&gohome.FeatureAttrsChangedEvt{
						FeatureID: btn.ID,
						Attrs:     feature.NewAttrs(state),
					})

				case *lutronExt.BtnReleaseEvt:
					sourceDev := p.System.DeviceByAddress(e.DeviceAddress)
					if sourceDev == nil {
						return
					}

					btn := sourceDev.FeatureTypeByAddress(feature.FTButton, e.Address)
					if btn == nil {
						return
					}

					state := feature.ButtonCloneAttrs(btn)
					state.Value = attr.ButtonStateReleased

					p.System.Services.EvtBus.Enqueue(&gohome.FeatureAttrsChangedEvt{
						FeatureID: btn.ID,
						Attrs:     feature.NewAttrs(state),
					})

				case *lutronExt.UnknownEvt:
					// Don't care, ignore
				}
			})
			p.scannerConn = nil

			log.V("%s stopped streaming events", p.Device)

			// stream doesn't always return an error, if the underlyinc connection closes
			if err == nil {
				err = fmt.Errorf("unknown error, stream ended unexpectedly")
			}
			p.Device.Connections.Release(conn, err)

			if err != nil {
				log.V("%s error streaming events, streaming stopped: %s", p.Device, err)
			}
		}
	}()
}

func (p *event) StopProducing() {
	p.producing = false
}
