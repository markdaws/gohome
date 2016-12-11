package lutron

import (
	"time"

	"github.com/go-home-iot/event-bus"
	lutronExt "github.com/go-home-iot/lutron"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/feature"
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
				p.Device.Connections.Release(conn, nil)
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
					/*
						sourceDevice, ok := p.Device.Devices[e.DeviceAddress]
						if !ok {
							return
						}
						_ = sourceDevice
					*/

					//TODO: Fix
					/*
						if btn := sourceDevice.ButtonByAddress(e.Address); btn != nil {
							p.System.Services.EvtBus.Enqueue(&gohome.ButtonPressEvt{
								BtnAddress:    btn.Address,
								BtnID:         btn.ID,
								BtnName:       btn.Name,
								DeviceName:    sourceDevice.Name,
								DeviceAddress: sourceDevice.Address,
								DeviceID:      sourceDevice.ID,
							})
						}*/

				case *lutronExt.BtnReleaseEvt:
					/*
						sourceDevice, ok := p.Device.Devices[e.DeviceAddress]
						if !ok {
							return
						}
						_ = sourceDevice
					*/

					//TODO: Fix
					/*
						if btn := sourceDevice.ButtonByAddress(e.Address); btn != nil {
							p.System.Services.EvtBus.Enqueue(&gohome.ButtonReleaseEvt{
								BtnAddress:    btn.Address,
								BtnID:         btn.ID,
								BtnName:       btn.Name,
								DeviceName:    sourceDevice.Name,
								DeviceAddress: sourceDevice.Address,
								DeviceID:      sourceDevice.ID,
							})
						}*/

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
