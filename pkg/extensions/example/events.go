package example

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/pkg/attr"
	"github.com/markdaws/gohome/pkg/feature"
	"github.com/markdaws/gohome/pkg/gohome"
	"github.com/markdaws/gohome/pkg/log"
)

type producer struct {
	producing bool
	System    *gohome.System
	Device    *gohome.Device
}

//============evtbus.Producer interface =============================================

func (p *producer) ProducerName() string {
	// A friendly string that will appear in debug information
	return fmt.Sprintf("example-producer-%s", p.Device.Name)
}

func (p *producer) StartProducing(b *evtbus.Bus) {
	// At some point when the extension is initialized this function will be called
	// you can then start listening for events from your device and report them
	// back to the system

	// If your device supports UPNP notifications, you can look at the
	// extensions/belkin/events.go file to see how you can easily subscribe
	// to those events

	// In this example we will pretend that we are getting updates from some hardware
	// just simulating with a timeout, we then put an event on the event bus so other
	// parts of the system can react to it.

	// The system can stop the producer by calling StopProducing(), so this bool keeps
	// track of that for us
	p.producing = true

	// IMORTANT: This function cannot be blocking, you must wrap all your code in
	// a go routine so the function can return immediately, if you don't do this you
	// may cause the system to eventually block

	go func() {
		for p.producing {
			// Pretend we are getting events from hardware, just sleep for a small time
			time.Sleep(time.Second * 10)

			// For other events you can react and produce, see gohome/events.go.

			// Loop through all of the features that the device exports, such as button,
			// light zone, sensor etc. enqueuing some random value for the value to
			// simulate things are changing
			for _, f := range p.Device.Features {

				// In discovery.go we gave the light feature address 1 and the
				// sensor feature address 2, that lets us distinguish them here,
				// we could also use feature.Type to choose between them
				switch f.Address {
				case "1":
					// A feature, such as a light zone, can have one or more attributes. An attribute is
					// just a value, so a light zone has brightness, onoff and hue attributes. The Features
					// field is a map of the attributes keyed by the attributes LocalID field. LocalID is just
					// a name that identifies the attribute in the feature.

					// When we want to update an attribute value, we access it by localID first, in the feature.go
					// file there are a bunch of consts defined that map to the names. Then you need to Clone() the
					// attribute before you update it

					onoff, brightness, _ := feature.LightZoneCloneAttrs(f)

					// Send back some random values
					brightness.Value = float32(rand.Intn(101))
					onoff.Value = int32(1 + rand.Intn(2))

					// Finally we send an event to the event bus so the system knows that this feature has been
					// updated.  We include only the attributes that have changed.
					p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(onoff, brightness),
					})

				case "2":
					// A sensor only has one attribute, so we can use attr.Only() to pick it out from the
					// map returned in f.Attrs. Only just picks the first item in the map.n
					openclose := attr.Only(f.Attrs).Clone()
					openclose.Value = int32(1 + rand.Intn(2))

					p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(openclose),
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
	// Return a friendly string that will be shown in debug information
	return fmt.Sprintf("example-consumer-%s", c.Device.Name)
}

func (c *consumer) StartConsuming(ch chan evtbus.Event) {
	// This function will be called at some point by the system, the channel that is
	// passed in will contain all the events that are published in the system, you can
	// check to make sure it is an event you care about and handle it

	c.consuming = true

	go func() {
		// Sit here waiting for events to inspect
		for e := range ch {
			if !c.consuming {
				return
			}

			// In this example, we check the events, if it is a ZoneReportEvt, which is an
			// event where the system is asking for specific zone to report their values, we
			// will check to see if the zones are ones we own, and if so, get their latest
			// values and report them back to the system.

			// For other events you can handle, look at gohome/events.go

			var evt *gohome.FeaturesReportEvt
			var ok bool
			if evt, ok = e.(*gohome.FeaturesReportEvt); !ok {
				continue
			}

			log.V("%s - %s", c.ConsumerName(), evt)

			for _, f := range c.Device.OwnedFeatures(evt.FeatureIDs) {
				// We gave the features different addresses in the discovery.go file, so we can
				// use it to distinguish them here, between the light and the sensor
				switch f.Address {
				case "1":
				case "2":
					// Get the latest value for this device, we will just return a random
					// value for the purpose of this example. We create the attribute with
					// a new value and set that in the event
					openclose := attr.Only(f.Attrs)
					openclose.Value = int32(1 + rand.Intn(2))

					c.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
						FeatureID: f.ID,
						Attrs:     feature.NewAttrs(openclose),
					})
				}
			}
		}
	}()
}

func (c *consumer) StopConsuming() {
	// Called by the system at some point, after this you should not consume
	// any more events
	c.consuming = false
}

//===================================================================================
