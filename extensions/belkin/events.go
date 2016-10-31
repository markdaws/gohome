package belkin

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"time"

	belkinExt "github.com/go-home-iot/belkin"
	"github.com/go-home-iot/event-bus"
	"github.com/go-home-iot/upnp"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/zone"
)

type makerConsumer struct {
	System *gohome.System
	Device *gohome.Device
	Sensor *gohome.Sensor
	Zone   *zone.Zone
	Name   string
}

func (c *makerConsumer) ConsumerName() string {
	return c.Name
}
func (c *makerConsumer) StartConsuming(ch chan evtbus.Event) {
	go func() {
		for e := range ch {
			switch evt := e.(type) {
			case *gohome.ZonesReportEvt:
				_ = evt
				//TODO:
			case *gohome.SensorsReportEvt:
				//TODO:
			}
		}
	}()
}
func (c *makerConsumer) StopConsuming() {
	//TODO:
}

type makerProducer struct {
	System    *gohome.System
	Device    *gohome.Device
	Sensor    *gohome.Sensor
	Zone      *zone.Zone
	Name      string
	SID       string
	Producing bool
}

var attrRegexp = regexp.MustCompile(`(<attributeList>.*</attributeList>)`)

//==================== upnp.Subscriber interface ========================

func (p *makerProducer) UPNPNotify(e upnp.NotifyEvent) {
	if !p.Producing {
		return
	}

	// Contents are double HTML encoded when returned from the device
	body := html.UnescapeString(html.UnescapeString(e.Body))

	attrList := attrRegexp.FindStringSubmatch(body)
	if attrList == nil || len(attrList) == 0 {
		return
	}

	attrs := belkinExt.ParseAttributeList(attrList[0])
	if attrs == nil {
		return
	}

	fmt.Printf("%+v\n", attrs)
	//TODO: If this is a switch state change then we need to log
	//      that as a seperate event type

	if attrs.Sensor != nil {
		p.System.Services.EvtBus.Enqueue(&gohome.SensorAttrChangedEvt{
			SensorID:   p.Sensor.ID,
			SensorName: p.Sensor.Name,
			Attr: gohome.SensorAttr{
				Name:     "sensor",
				Value:    strconv.Itoa(*attrs.Sensor),
				DataType: gohome.SDTInt,
			},
		})
	} else if attrs.Switch != nil {
		p.System.Services.EvtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
			ZoneName: p.Zone.Name,
			ZoneID:   p.Zone.ID,
			Level:    cmd.Level{Value: float32(*attrs.Switch)},
		})
	}
}

//=======================================================================

func (p *makerProducer) ProducerName() string {
	return p.Name
}

func (p *makerProducer) StartProducing(b *evtbus.Bus) {
	log.V("producer [%s] start producing", p.ProducerName())

	go func() {
		p.Producing = true
		for p.Producing {
			//TODO: What about if we lose connection to this device or need a new SID?

			// The make has a sensor and a switch state, need to notify these changes
			// to the event bus
			sid, err := p.System.Services.UPNP.Subscribe(
				p.Device.Address+"/upnp/event/basicevent1",
				"",
				300,
				true,
				p)

			if err != nil {
				// log failure, keep trying to subscribe to the target device
				// there may be network issues
				log.V("[%s] failed to subscribe to upnp: %s", p.ProducerName(), err)
				time.Sleep(time.Second * 10)
			} else {
				p.SID = sid
				break
			}
		}
	}()
}

func (p *makerProducer) StopProducing() {
	p.Producing = false

	//TODO: upnp should have want to ping subscribers and see
	// if they still want events then evict if not
	err := p.System.Services.UPNP.Unsubscribe(p.SID)
	if err != nil {
		log.V("error during unsusbscribe [%s]: %s", p.ProducerName(), err)
	}
}
