package testing

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/feature"
)

type producer struct {
	producing bool
	System    *gohome.System
	Device    *gohome.Device
}

//============evtbus.Producer interface =============================================

func (p *producer) ProducerName() string {
	return fmt.Sprintf("example-producer-%s", p.Device.Name)
}

func (p *producer) StartProducing(b *evtbus.Bus) {
	p.producing = true
	go func() {
		for p.producing {
			time.Sleep(time.Second * 10)

			// For each feature we own, just send back some random value. We can switch
			// on the feature address to know what kind of feature it is because we assigned
			// those in discover.go
			for _, f := range p.Device.Features {
				switch f.Address {
				case "1":
					onoff, _, _ := feature.LightZoneCloneAttrs(f)
					onoff.Value = int32(1 + rand.Intn(2))
					p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(onoff),
					})

				case "2":
					_, brightness, _ := feature.LightZoneCloneAttrs(f)
					brightness.Value = float32(rand.Intn(101))
					p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(brightness),
					})

				case "3":
					openclose := attr.Only(f.Attrs).Clone()
					openclose.Value = int32(1 + rand.Intn(2))
					p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(openclose),
					})

				case "4":
					onoff := feature.SwitchCloneAttrs(f)
					onoff.Value = int32(1 + rand.Intn(2))
					p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(onoff),
					})

				case "5":
					currentTemp, targetTemp := feature.HeatZoneCloneAttrs(f)
					currentTemp.Value = int32(40 + rand.Intn(41))
					targetTemp.Value = int32(40 + rand.Intn(41))
					p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(currentTemp, targetTemp),
					})

				case "6":
					_, offset := feature.WindowTreatmentCloneAttrs(f)
					offset.Value = float32(rand.Intn(101))
					p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(offset),
					})
				}
			}
		}
	}()
}

func (p *producer) StopProducing() {
	// This is called by the system at some point, make sure we don't produce any
	// more events after this is called
	p.producing = false
}

//===================================================================================

//============evtbus.Consumer interface =============================================

type consumer struct {
	consuming bool
	Device    *gohome.Device
	System    *gohome.System
}

func (c *consumer) ConsumerName() string {
	return fmt.Sprintf("example-consumer-%s", c.Device.Name)
}

func (c *consumer) StartConsuming(ch chan evtbus.Event) {
	c.consuming = true
}

func (c *consumer) StopConsuming() {
	c.consuming = false
}

//===================================================================================
