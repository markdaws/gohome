package honeywell

import (
	"context"
	"strconv"
	"time"

	"github.com/go-home-iot/event-bus"
	honeywellExt "github.com/go-home-iot/honeywell"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/feature"
	"github.com/markdaws/gohome/log"
)

type consumer struct {
	System *gohome.System
	Device *gohome.Device
}

func (c *consumer) ConsumerName() string {
	return "honeywell"
}

func (c *consumer) StartConsuming(ch chan evtbus.Event) {
	go func() {
		var thermostat honeywellExt.Thermostat

		for e := range ch {
			switch evt := e.(type) {
			case *gohome.FeaturesReportEvt:
				for _, f := range c.Device.OwnedFeatures(evt.FeatureIDs) {
					if thermostat == nil {
						devID, err := strconv.Atoi(c.Device.Address)
						if err != nil {
							log.V("honeywell device does not have valid device ID in the address field %s, feature ID: %s",
								c.Device.Address, f.ID)
							continue
						}

						thermostat = honeywellExt.NewThermostat(devID)
						ctx := context.TODO()
						err = thermostat.Connect(ctx, c.Device.Auth.Login, c.Device.Auth.Password)
						if err != nil {
							log.V("failed to connect to honeywell thermostat: %s", err)
							thermostat = nil
							continue
						}
					}

					ctx := context.TODO()
					status, err := thermostat.FetchStatus(ctx)
					if err != nil {
						log.V("failed to fetch honeywell status: %s", err)

						// Set this to nil so that next time we try to reconnect again
						thermostat = nil
						continue
					}

					if !status.DeviceLive {
						continue
					}
					current := status.LatestData.UIData.DispTemperature
					target := status.LatestData.UIData.HeatSetpoint

					currentTemp, targetTemp := feature.HeatZoneCloneAttrs(f)
					currentTemp.Value = current
					targetTemp.Value = target

					c.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(currentTemp, targetTemp),
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
	System    *gohome.System
	Device    *gohome.Device
	Producing bool
}

func (p *producer) ProducerName() string {
	return "honeywell"
}

func (p *producer) StartProducing(b *evtbus.Bus) {
	log.V("producer [%s] start producing", p.ProducerName())

	go func() {
		p.Producing = true

		var thermostat honeywellExt.Thermostat
		for p.Producing {
			// Since we can't get push notification, just poll every 30 seconds for the current value
			time.Sleep(time.Second * 30)

			var f *feature.Feature
			for _, f = range p.Device.Features {
				if f.Type == feature.FTHeatZone {
					break
				}
			}
			if f == nil {
				log.V("unable to find honeywell heat zone")
				continue
			}

			if thermostat == nil {
				devID, err := strconv.Atoi(p.Device.Address)
				if err != nil {
					log.V("honeywell device does not have valid device ID in the address field %s, feature ID: %s",
						p.Device.Address, f.ID)
					continue
				}

				thermostat = honeywellExt.NewThermostat(devID)
				ctx := context.TODO()
				err = thermostat.Connect(ctx, p.Device.Auth.Login, p.Device.Auth.Password)
				if err != nil {
					log.V("failed to connect to honeywell thermostat: %s", err)
					thermostat = nil
					continue
				}
			}

			ctx := context.TODO()
			status, err := thermostat.FetchStatus(ctx)
			if err != nil {
				log.V("failed to fetch honeywell status: %s", err)

				// Set this to nil so that next time we try to reconnect again
				thermostat = nil
				continue
			}

			if !status.DeviceLive {
				continue
			}

			current := status.LatestData.UIData.DispTemperature
			target := status.LatestData.UIData.HeatSetpoint

			currentTemp, targetTemp := feature.HeatZoneCloneAttrs(f)
			currentTemp.Value = current
			targetTemp.Value = target

			p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
				FeatureID: f.ID,
				Attrs:     feature.NewAttrs(currentTemp, targetTemp),
			})
		}
		log.V("%s - stopped producing events", p.ProducerName())
	}()
}

func (p *producer) StopProducing() {
	p.Producing = false
}
