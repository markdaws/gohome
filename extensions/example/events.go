package example

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
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

			// For other events you can react and produce, see gohome/events.go

			// Pretend this was a zone level change event, report it
			for _, zone := range p.Device.Zones {
				lvl := float32(rand.Intn(101))

				p.System.Services.EvtBus.Enqueue(&gohome.ZoneLevelReportingEvt{
					ZoneName: zone.Name,
					ZoneID:   zone.ID,
					Level: cmd.Level{
						// Just send back some random number between 0-100
						Value: lvl,
						R:     0.0,
						G:     0.0,
						B:     0.0,
					},
				})
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

			var evt *gohome.ZonesReportEvt
			var ok bool
			if evt, ok = e.(*gohome.ZonesReportEvt); !ok {
				// Not an event we care about, ignore
				continue
			}

			// Look at all the zones which the system is requesting values for, if we own
			// any of them we need to get their latest values
			for _, zone := range c.Device.OwnedZones(evt.ZoneIDs) {
				log.V("%s - %s", c.ConsumerName(), evt)

				// Get the latest value for this device, we will just return a random
				// value for the purpose of this example
				c.System.Services.EvtBus.Enqueue(&gohome.ZoneLevelReportingEvt{
					ZoneName: zone.Name,
					ZoneID:   zone.ID,
					Level: cmd.Level{
						Value: float32(rand.Intn(101)),
						R:     0.0,
						G:     0.0,
						B:     0.0,
					},
				})
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
