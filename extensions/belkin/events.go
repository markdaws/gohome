package belkin

import (
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
				for _, zone := range c.Device.Zones {
					if _, ok := evt.ZoneIDs[zone.ID]; !ok {
						continue
					}

					dev := &belkinExt.Device{
						Scan: belkinExt.ScanResponse{
							SearchType: belkinExt.DTMaker,
							Location:   c.Device.Address,
						},
					}
					attrs, err := dev.FetchAttributes()
					if err != nil {
						log.V("Belkin - failed to fetch attrs: %s", err)
						continue
					}

					if attrs.Switch != nil {
						c.System.Services.EvtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
							ZoneName: c.Zone.Name,
							ZoneID:   c.Zone.ID,
							Level:    cmd.Level{Value: float32(*attrs.Switch)},
						})
					}
				}

			case *gohome.SensorsReportEvt:
				for _, sensor := range c.Device.Sensors {
					if _, ok := evt.SensorIDs[sensor.ID]; !ok {
						continue
					}

					dev := &belkinExt.Device{
						Scan: belkinExt.ScanResponse{
							SearchType: belkinExt.DTMaker,
							Location:   c.Device.Address,
						},
					}
					attrs, err := dev.FetchAttributes()
					if err != nil {
						log.V("Belkin - failed to fetch attrs: %s", err)
						continue
					}

					if attrs.Sensor != nil {
						attr := sensor.Attr
						attr.Value = strconv.Itoa(*attrs.Sensor)
						c.System.Services.EvtBus.Enqueue(&gohome.SensorAttrChangedEvt{
							SensorName: sensor.Name,
							SensorID:   sensor.ID,
							Attr:       attr,
						})
					}

				}
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

	if attrs.Sensor != nil {
		p.System.Services.EvtBus.Enqueue(&gohome.SensorAttrChangedEvt{
			SensorID:   p.Sensor.ID,
			SensorName: p.Sensor.Name,
			Attr: gohome.SensorAttr{
				Name:     "sensor",
				Value:    strconv.Itoa(*attrs.Sensor),
				DataType: gohome.SDTInt,
				States:   p.Sensor.Attr.States,
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
