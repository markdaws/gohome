package belkin

import (
	"strings"
	"time"

	belkinExt "github.com/go-home-iot/belkin"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/zone"
)

var infos = []gohome.DiscovererInfo{gohome.DiscovererInfo{
	ID:          "belkin.wemo.insight",
	Name:        "Belkin WeMo Insight",
	Description: "Discover Belkin WeMo Insight devices",
}, gohome.DiscovererInfo{
	ID:          "belkin.wemo.maker",
	Name:        "Belkin WeMo Maker",
	Description: "Discover Belkin WeMo Maker devices",
}}

type discovery struct {
	//TODO: Needed?
	System *gohome.System
}

func (d *discovery) Discoverers() []gohome.DiscovererInfo {
	return infos
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	switch ID {
	case "belkin.wemo.insight":
		return &discoverer{scanType: belkinExt.DTInsight, info: infos[0]}
	case "belkin.wemo.maker":
		return &discoverer{scanType: belkinExt.DTMaker, info: infos[1]}
	default:
		return nil
	}
}

type discoverer struct {
	scanType belkinExt.DeviceType
	info     gohome.DiscovererInfo
}

func (d *discoverer) Info() gohome.DiscovererInfo {
	return d.info
}

func (d *discoverer) ScanDevices(sys *gohome.System, uiFields map[string]string) (*gohome.DiscoveryResults, error) {

	log.V("scanning belkin")

	responses, err := belkinExt.Scan(d.scanType, 5)
	if err != nil {
		log.V("scan err: %s", err)
		return nil, err
	}

	devices := make([]*gohome.Device, len(responses))
	for i, devInfo := range responses {
		err := devInfo.Load(time.Second * 5)
		if err != nil {
			// Keep going, try to get as many as we can
			log.V("failed to load device information: %s", err)
			continue
		}

		dev := gohome.NewDevice(
			sys.NextGlobalID(),
			devInfo.FriendlyName,
			devInfo.ModelDescription,
			devInfo.ModelNumber,
			devInfo.ModelName,
			devInfo.FirmwareVersion,
			strings.Replace(devInfo.Scan.Location, "/setup.xml", "", -1),
			nil,
			nil,
			nil,
			nil,
		)

		z := &zone.Zone{
			ID:          sys.NextGlobalID(),
			Address:     "1",
			Name:        devInfo.FriendlyName,
			Description: devInfo.ModelDescription,
			DeviceID:    dev.ID,
			Type:        zone.ZTSwitch,
			Output:      zone.OTBinary,
		}
		dev.AddZone(z)

		if d.scanType == belkinExt.DTMaker {
			sensor := &gohome.Sensor{
				ID:          sys.NextGlobalID(),
				Address:     "1",
				Name:        devInfo.FriendlyName + " - sensor",
				Description: "",
				DeviceID:    dev.ID,
				Attr: gohome.SensorAttr{
					Name:     "sensor",
					Value:    "-1",
					DataType: gohome.SDTInt,
					States: map[string]string{
						"0": "Closed",
						"1": "Open",
					},
				},
			}
			dev.AddSensor(sensor)
		}
		devices[i] = dev
	}

	return &gohome.DiscoveryResults{
		Devices: devices,
	}, nil

}
