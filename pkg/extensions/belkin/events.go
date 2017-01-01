package belkin

import (
	"html"
	"regexp"
	"sync"
	"time"

	belkinExt "github.com/go-home-iot/belkin"
	"github.com/go-home-iot/event-bus"
	"github.com/go-home-iot/upnp"
	"github.com/markdaws/gohome/pkg/attr"
	"github.com/markdaws/gohome/pkg/feature"
	"github.com/markdaws/gohome/pkg/gohome"
	"github.com/markdaws/gohome/pkg/log"
)

type consumer struct {
	System     *gohome.System
	Device     *gohome.Device
	Name       string
	DeviceType belkinExt.DeviceType
}

func (c *consumer) ConsumerName() string {
	return c.Name
}

func (c *consumer) StartConsuming(ch chan evtbus.Event) {
	go func() {
		for e := range ch {
			switch evt := e.(type) {
			case *gohome.FeaturesReportEvt:
				switch c.DeviceType {
				case belkinExt.DTInsight:
					c.reportInsight(evt)
				case belkinExt.DTMaker:
					c.reportMaker(evt)
				}
			}
		}
	}()
}

func (c *consumer) reportInsight(evt *gohome.FeaturesReportEvt) {

	dev := &belkinExt.Device{
		Scan: belkinExt.ScanResponse{
			SearchType: string(c.DeviceType),
			Location:   c.Device.Address,
		},
	}

	// Get all of the features we own, then fetch latest values
	for _, f := range c.Device.OwnedFeatures(evt.FeatureIDs) {
		switch f.Type {
		case feature.FTOutlet:
			state, err := dev.FetchBinaryState(time.Second * 5)
			if err != nil {
				log.V("Belkin - failed to fetch binary state: %s", err)
				continue
			}

			var onOffVal int32
			switch state {
			case 0:
				onOffVal = attr.OnOffOff
			case 1, 8:
				onOffVal = attr.OnOffOn
			}

			onoff := feature.OutletCloneAttrs(f)
			onoff.Value = onOffVal
			c.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
				FeatureID: f.ID,
				Attrs:     feature.NewAttrs(onoff),
			})
		}
	}
}

func (c *consumer) reportMaker(evt *gohome.FeaturesReportEvt) {
	dev := &belkinExt.Device{
		Scan: belkinExt.ScanResponse{
			SearchType: string(c.DeviceType),
			Location:   c.Device.Address,
		},
	}
	var once sync.Once
	var fetchErr error
	var attrs *belkinExt.DeviceAttributes

	// Get all of the features we own, then fetch latest values
	for _, f := range c.Device.OwnedFeatures(evt.FeatureIDs) {

		// We only need to do one call to the device to get all of the feature info
		// so usong once will only call the request once
		once.Do(func() {
			attrs, fetchErr = dev.FetchAttributes(time.Second * 5)
		})

		if fetchErr != nil {
			log.V("Belkin - failed to fetch attrs: %s", fetchErr)
			continue
		}

		switch f.Type {
		case feature.FTSwitch:
			if attrs.Switch == nil {
				continue
			}

			var onOffVal int32
			switch *attrs.Switch {
			case 0:
				onOffVal = attr.OnOffOff
			case 1, 8:
				onOffVal = attr.OnOffOn
			}

			onoff := feature.SwitchCloneAttrs(f)
			onoff.Value = onOffVal
			c.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
				FeatureID: f.ID,
				Attrs:     feature.NewAttrs(onoff),
			})

		case feature.FTSensor:

			if attrs.Sensor == nil {
				continue
			}
			var openCloseVal int32
			switch *attrs.Sensor {
			case 0:
				openCloseVal = attr.OpenCloseClosed
			case 1:
				openCloseVal = attr.OpenCloseOpen
			}

			// Sensors only have one attribute, it can be any kind of attribute with any local
			// ID so we just grab the only attribute and update it
			attribute := attr.Only(f.Attrs).Clone()
			attribute.Value = openCloseVal
			c.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
				FeatureID: f.ID,
				Attrs:     feature.NewAttrs(attribute),
			})
		}
	}
}

func (c *consumer) StopConsuming() {
	//TODO:
}

type producer struct {
	System     *gohome.System
	Device     *gohome.Device
	Name       string
	SID        string
	Producing  bool
	DeviceType belkinExt.DeviceType
}

var attrRegexp = regexp.MustCompile(`(<attributeList>.*</attributeList>)`)
var binaryRegexp = regexp.MustCompile(`(<BinaryState>.*</BinaryState>)`)

//==================== upnp.Subscriber interface ========================

func (p *producer) UPNPNotify(e upnp.NotifyEvent) {
	if !p.Producing {
		return
	}

	var sensor *feature.Feature
	var swtch *feature.Feature
	var outlet *feature.Feature
	for _, f := range p.Device.Features {
		switch f.Type {
		case feature.FTSensor:
			sensor = f
		case feature.FTSwitch:
			swtch = f
		case feature.FTOutlet:
			outlet = f
		}
	}

	// Contents are double HTML encoded when returned from the device
	body := html.UnescapeString(html.UnescapeString(e.Body))

	// This could be a response with an attribute list, or it could be a binary state property
	attrList := attrRegexp.FindStringSubmatch(body)
	if attrList != nil && len(attrList) != 0 {
		attrs := belkinExt.ParseAttributeList(attrList[0])
		if attrs == nil {
			return
		}

		if attrs.Sensor != nil && sensor != nil {
			var openCloseVal int32
			switch *attrs.Sensor {
			case 0:
				openCloseVal = attr.OpenCloseClosed
			case 1:
				openCloseVal = attr.OpenCloseOpen
			}

			attribute := attr.Only(sensor.Attrs).Clone()
			attribute.Value = openCloseVal
			p.System.Services.EvtBus.Enqueue(&gohome.FeatureAttrsChangedEvt{
				FeatureID: sensor.ID,
				Attrs:     feature.NewAttrs(attribute),
			})
		}

		if attrs.Switch != nil {
			var onOffVal int32
			switch *attrs.Switch {
			case 0:
				onOffVal = attr.OnOffOff
			case 1:
				onOffVal = attr.OnOffOn
			}

			if swtch != nil {
				onoff := feature.SwitchCloneAttrs(swtch)
				onoff.Value = onOffVal
				p.System.Services.EvtBus.Enqueue(&gohome.FeatureAttrsChangedEvt{
					FeatureID: swtch.ID,
					Attrs:     feature.NewAttrs(onoff),
				})
			}
			if outlet != nil {
				onoff := feature.OutletCloneAttrs(outlet)
				onoff.Value = onOffVal
				p.System.Services.EvtBus.Enqueue(&gohome.FeatureAttrsChangedEvt{
					FeatureID: outlet.ID,
					Attrs:     feature.NewAttrs(onoff),
				})
			}
		}
	} else if swtch != nil || outlet != nil {
		binary := binaryRegexp.FindStringSubmatch(body)
		if binary == nil || len(binary) == 0 {
			return
		}

		states := belkinExt.ParseBinaryState(binary[0])

		// Note for onoff 1 and 8 mean on, normalize to 1
		level := states.OnOff
		if level == 8 {
			level = 1
		}

		var onOffVal int32
		switch level {
		case 0:
			onOffVal = attr.OnOffOff
		case 1:
			onOffVal = attr.OnOffOn
		}

		if swtch != nil {
			onoff := feature.SwitchCloneAttrs(swtch)
			onoff.Value = onOffVal
			p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
				FeatureID: swtch.ID,
				Attrs:     feature.NewAttrs(onoff),
			})
		}
		if outlet != nil {
			onoff := feature.OutletCloneAttrs(outlet)
			onoff.Value = onOffVal
			p.System.Services.EvtBus.Enqueue(&gohome.FeatureReportingEvt{
				FeatureID: outlet.ID,
				Attrs:     feature.NewAttrs(onoff),
			})
		}
	}
}

//=======================================================================

func (p *producer) ProducerName() string {
	return p.Name
}

func (p *producer) StartProducing(b *evtbus.Bus) {
	log.V("producer [%s] start producing", p.ProducerName())

	go func() {
		p.Producing = true

		for p.Producing {
			log.V("%s - subscribing to UPNP, SID:%s", p.ProducerName(), p.SID)

			// The make has a sensor and a switch state, need to notify these changes
			// to the event bus
			sid, err := p.System.Services.UPNP.Subscribe(
				p.Device.Address+"/upnp/event/basicevent1",
				p.SID,
				120,
				true,
				p)

			if err != nil {
				// log failure, keep trying to subscribe to the target device
				// there may be network issues, if this is a renew, the old SID
				// might have expired, so reset so we get a new one
				log.V("[%s] failed to subscribe to upnp: %s", p.ProducerName(), err)
				p.SID = ""
				time.Sleep(time.Second * 10)
			} else {
				// We got a sid, now sleep then renew the subscription
				p.SID = sid
				log.V("%s - subscribed to UPNP, SID:%s", p.ProducerName(), sid)
				time.Sleep(time.Second * 100)
			}
		}

		log.V("%s - stopped producing events", p.ProducerName())
	}()
}

func (p *producer) StopProducing() {
	p.Producing = false

	err := p.System.Services.UPNP.Unsubscribe(p.SID)
	if err != nil {
		log.V("error during unsusbscribe [%s]: %s", p.ProducerName(), err)
	}
}
