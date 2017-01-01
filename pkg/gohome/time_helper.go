package gohome

import (
	"time"

	"github.com/cpucycle/astrotime"
	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/pkg/clock"
	"github.com/markdaws/gohome/pkg/log"
)

// TimeHelper helps with some time related functionality, such as firing sunrise/sunset events
type TimeHelper struct {
	Time      clock.Time
	System    *System
	Latitude  float64
	Longitude float64
	Produce   bool
}

func (th *TimeHelper) ProducerName() string {
	return "TimeHelper"
}

func (th *TimeHelper) StartProducing(b *evtbus.Bus) {

	th.Produce = true

	if th.Latitude == 0 && th.Longitude == 0 {
		log.V("Sunrise/Sunset events will not be fired, location not set.  Update config.json with the correct lat/long values then restart the server")
		return
	}

	log.V("TimeHelper - initializing")

	go func() {
		for {
			now := th.Time.Now()
			t := astrotime.NextSunrise(now, th.Latitude, th.Longitude)
			tzname, _ := t.Zone()
			log.V("The next sunrise (lat:%f, long:%f) is %d:%02d %s on %d/%d/%d.",
				th.Latitude, th.Longitude, t.Hour(), t.Minute(), tzname, t.Month(), t.Day(), t.Year())

			<-th.Time.After(t.Sub(now))

			if th.Produce {
				th.System.Services.EvtBus.Enqueue(&SunriseEvt{})
			}

			// Small delay so we don't fire multiple times in the loop for the same event
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			now := th.Time.Now()
			t := astrotime.NextSunset(now, th.Latitude, th.Longitude)
			tzname, _ := t.Zone()
			log.V("The next sunset (lat:%f, long:%f) is %d:%02d %s on %d/%d/%d.",
				th.Latitude, th.Longitude, t.Hour(), t.Minute(), tzname, t.Month(), t.Day(), t.Year())

			<-th.Time.After(t.Sub(now))

			if th.Produce {
				th.System.Services.EvtBus.Enqueue(&SunsetEvt{})
			}
			time.Sleep(time.Second)
		}
	}()
}

func (th *TimeHelper) StopProducing() {
	th.Produce = false
}
