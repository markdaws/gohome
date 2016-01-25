package gohome

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/connectedbytcp"
	"github.com/markdaws/gohome/zone"
)

type Importer interface {
	ImportFromFile(path, importerID string, cmdProcessor CommandProcessor) (*System, error)
}

type importer struct {
}

func (i importer) ImportFromFile(path, importerID string, cp CommandProcessor) (*System, error) {
	switch importerID {
	case "L-BDGPRO2-WH":
		return importL_BDGPRO2_WH(path, "1", cp)
	default:
		return nil, errors.New("unknown import type: " + importerID)
	}
}

func NewImporter() Importer {
	return importer{}
}

// Used for integration reports from Lutron Smart Bridge Pro
func importL_BDGPRO2_WH(integrationReportPath, smartBridgeProID string, cmdProcessor CommandProcessor) (*System, error) {

	//TODO: Handle non runtime panic
	bytes, err := ioutil.ReadFile(integrationReportPath)
	if err != nil {
		return nil, err
	}

	var configJson map[string]interface{}
	if err = json.Unmarshal(bytes, &configJson); err != nil {
		return nil, err
	}

	system := NewSystem("Lutron Smart Bridge Pro", "Lutron Smart Bridge Pro", cmdProcessor)

	root, ok := configJson["LIPIdList"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Missing LIPIdList key, or value not a map")
	}
	devices, ok := root["Devices"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Devices key, or value not a map")
	}

	fmt.Println("\nDEVICES")

	var makeDevice = func(modelNumber, address string, deviceMap map[string]interface{}, sys *System, stream bool, auth *comm.Auth) Device {
		var deviceName string = deviceMap["Name"].(string)

		device := NewDevice(
			modelNumber,
			address,
			sys.NextGlobalID(),
			deviceName,
			"",
			stream,
			auth)

		for _, buttonMap := range deviceMap["Buttons"].([]interface{}) {
			button := buttonMap.(map[string]interface{})
			btnNumber := strconv.FormatFloat(button["Number"].(float64), 'f', 0, 64)

			var btnName string
			if name, ok := button["Name"]; ok {
				btnName = name.(string)
			} else {
				btnName = "Button " + btnNumber
			}

			b := &Button{
				Address:     btnNumber,
				ID:          sys.NextGlobalID(),
				Name:        btnName,
				Description: btnName,
				Device:      device,
			}
			device.Buttons()[btnNumber] = b
			system.AddButton(b)
		}

		return device
	}

	var makeScenes = func(sceneContainer map[string]*Scene, deviceMap map[string]interface{}, sbp Device) error {
		buttons, ok := deviceMap["Buttons"].([]interface{})
		if !ok {
			return errors.New("Missing Buttons key, or value not array")
		}

		//var deviceID string = strconv.FormatFloat(deviceMap["ID"].(float64), 'f', 0, 64)
		for _, buttonMap := range buttons {
			button, ok := buttonMap.(map[string]interface{})
			if !ok {
				return errors.New("Expected Button elements to be objects")
			}
			if name, ok := button["Name"]; ok && !strings.HasPrefix(name.(string), "Button ") {
				fmt.Printf("  Scene %d: %s\n", int(button["Number"].(float64)), name)

				var buttonID string = strconv.FormatFloat(button["Number"].(float64), 'f', 0, 64)
				var buttonName = button["Name"].(string)

				var globalID = system.NextGlobalID()
				var btn = sbp.Buttons()[buttonID]
				sceneContainer[globalID] = &Scene{
					Address:     buttonID,
					ID:          globalID,
					Name:        buttonName,
					Description: buttonName,
					Commands: []cmd.Command{
						&cmd.ButtonPress{
							ButtonAddress: btn.Address,
							ButtonID:      btn.ID,
							DeviceName:    sbp.Name(),
							DeviceAddress: sbp.Address(),
							DeviceID:      sbp.ID(),
						},
						&cmd.ButtonRelease{
							ButtonAddress: btn.Address,
							ButtonID:      btn.ID,
							DeviceName:    sbp.Name(),
							DeviceAddress: sbp.Address(),
							DeviceID:      sbp.ID(),
						},
					},
				}
			}
		}

		return nil
	}

	// First need to find the Smart Bridge Pro since it is needed to make scenes and zones
	var sbp Device
	for _, deviceMap := range devices {
		device, ok := deviceMap.(map[string]interface{})
		if !ok {
			return nil, errors.New("Expected Devices elements to be objects")
		}

		var deviceID = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			sbp = makeDevice("L-BDGPRO2-WH", "192.168.0.10:23", device, system, true, &comm.Auth{
				Login:    "lutron",
				Password: "integration",
			})
			sbp.Auth().Authenticator = sbp
			makeScenes(system.Scenes, device, sbp)
			break
		}
	}

	if sbp == nil {
		return nil, errors.New("Did not find Smart Bridge Pro with ID:" + smartBridgeProID)
	}
	system.AddDevice(sbp)

	for _, deviceMap := range devices {
		device, ok := deviceMap.(map[string]interface{})
		if !ok {
			return nil, errors.New("Expected Devices elements to be objects")
		}

		fmt.Printf("%d: %s\n", int(device["ID"].(float64)), device["Name"])

		// Don't want to re-add the SBP
		var deviceID = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			continue
		}
		gohomeDevice := makeDevice("", deviceID, device, system, false, nil)
		system.AddDevice(gohomeDevice)
		sbp.Devices()[gohomeDevice.Address()] = gohomeDevice
	}

	zones, ok := root["Zones"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Zones key")
	}

	fmt.Println("\nZONES")
	for _, zoneMap := range zones {
		z := zoneMap.(map[string]interface{})
		fmt.Printf("%d: %s\n", int(z["ID"].(float64)), z["Name"])

		var zoneID = strconv.FormatFloat(z["ID"].(float64), 'f', 0, 64)
		var zoneName = z["Name"].(string)
		var zoneTypeFinal = zone.ZTLight
		if zoneType, ok := z["Type"].(string); ok {
			switch zoneType {
			case "light":
				zoneTypeFinal = zone.ZTLight
			case "shade":
				zoneTypeFinal = zone.ZTShade
			}
		}
		var outputTypeFinal = zone.OTContinuous
		if outputType, ok := z["Output"].(string); ok {
			switch outputType {
			case "binary":
				outputTypeFinal = zone.OTBinary
			case "continuous":
				outputTypeFinal = zone.OTContinuous
			}
		}
		newZone := &zone.Zone{
			Address:     zoneID,
			ID:          system.NextGlobalID(),
			Name:        zoneName,
			Description: zoneName,
			DeviceID:    sbp.ID(),
			Type:        zoneTypeFinal,
			Output:      outputTypeFinal,
		}
		system.AddZone(newZone)
		sbp.Zones()[newZone.Address] = newZone
	}

	//TODO: Move
	importConnectedByTCP(system)
	importGoHomeHub(system)
	return system, nil
}

/*
type tcpListener struct {
}

func (tcpListener) Response(m gossdp.ResponseMessage) {
	fmt.Printf("%+v\n", m)
}*/

//TODO: Temp function - import from UI
func importConnectedByTCP(system *System) {

	/*
		b := tcpListener{}
		c, err := gossdp.NewSsdpClient(b)
		if err != nil {
			log.Println("Failed to start client: ", err)
			return
		}
		defer c.Stop()
		go c.Start()

		err = c.ListenFor("urn:greenwavereality-com:service:gop:1")
		time.Sleep(30 * time.Second)
	*/
	/*
		{MaxAge:7200 SearchType:urn:greenwavereality-com:service:gop:1 DeviceId:71403833960916 Usn:uuid:71403833960916::urn:greenwavereality-com:service:gop:1 Location:https://192.168.0.23 Server:linux UPnP/1.1 Apollo3/3.0.74 RawResponse:0xc2080305a0 Urn:urn:greenwavereality-com:service:gop:1}
	*/

	/*
			//1. Press sync button on hub
			//2. Execute following url
			//https://192.168.0.23/gwr/gop.php?cmd=GWRLogin&data=%3Cgip%3E%3Cversion%3E1%3C/version%3E%3Cemail%3Etest%3C/email%3E%3Cpassword%3Etest%3C/password%3E%3C/gip%3E
			//3. Get response: <gip><version>1</version><rc>200</rc><token>ar6thtpqg6yinh219pn0c4t814dqkye1f0j3sfye</token></gip>
			//4. Use token in commands

		data := "cmd=GWRBatch&data=<gwrcmds><gwrcmd><gcmd>RoomGetCarousel</gcmd><gdata><gip><version>1</version><token>79tz3vbbop9pu5fcen60p97ix3mbvd3sblhjmz21</token><fields>name,control,power,product,class,realtype,status</fields></gip></gdata></gwrcmd></gwrcmds>&fmt=xml"
		_ = data
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		slc := "cmd=GWRBatch&data=<gwrcmds><gwrcmd><gcmd>DeviceSendCommand</gcmd><gdata><gip><version>1</version><token>79tz3vbbop9pu5fcen60p97ix3mbvd3sblhjmz21</token><did>216438039298518643</did><value>100</value><type>level</type></gip></gdata></gwrcmd></gwrcmds>&fmt=xml"
		_ = slc
		resp, err := client.Post("https://192.168.0.23/gwr/gpo.php", "text/xml; charset=\"utf-8\"", bytes.NewReader([]byte(data)))
		xx, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(xx))
		fmt.Println(err)
	*/

	data, err := connectedbytcp.RoomGetCarousel("https://192.168.0.23", "79tz3vbbop9pu5fcen60p97ix3mbvd3sblhjmz21")
	fmt.Printf("%+v\n", data)
	fmt.Println(err)

	token, err := connectedbytcp.GetToken("https://192.168.0.23")
	fmt.Printf("TOKEN: %s", token)
	fmt.Println(err)

	tcp := NewDevice(
		"TCP600GWB",
		"https://192.168.0.23",
		system.NextGlobalID(),
		"ConnectedByTcp Hub",
		"Description",
		false,
		//TODO: Remove from NewDevice
		&comm.Auth{
			Token: "79tz3vbbop9pu5fcen60p97ix3mbvd3sblhjmz21",
		})

	/*
		//TODO: Fix
		tcp2 := tcp.(*Tcp600gwbDevice)
		tcp2.Token = "79tz3vbbop9pu5fcen60p97ix3mbvd3sblhjmz21"
		tcp2.Host = "https://192.168.0.23"
	*/

	zoneID := "216438039298518643"
	z := &zone.Zone{
		Address:     zoneID,
		ID:          system.NextGlobalID(),
		Name:        "bulb1",
		Description: "tcp - bulb1",
		DeviceID:    tcp.ID(),
		Type:        zone.ZTLight,
		Output:      zone.OTContinuous,
		Controller:  "TCP - LED A19 11W",
	}
	fmt.Println("BULB ID: " + z.ID)
	tcp.Zones()[z.Address] = z
	system.AddZone(z)
	system.AddDevice(tcp)

	z1 := system.Zones["142"]
	z2 := system.Zones["153"]
	s := &Scene{
		Address:     "xxx",
		ID:          system.NextGlobalID(),
		Name:        "Synthetic Scene",
		Description: "Scene to control lutron + tcp lights",
		Commands: []cmd.Command{
			&cmd.ZoneSetLevel{
				ZoneAddress: z1.Address,
				ZoneID:      z1.ID,
				ZoneName:    z1.Name,
				Level:       cmd.Level{Value: 30},
			},
			&cmd.ZoneSetLevel{
				ZoneAddress: z2.Address,
				ZoneID:      z2.ID,
				ZoneName:    z2.Name,
				Level:       cmd.Level{Value: 30},
			},
		},
	}
	system.AddScene(s)
}

func importGoHomeHub(system *System) {

	ghh := NewDevice(
		"GoHomeHub",
		"gohomehub",
		system.NextGlobalID(),
		"GoHome Hub",
		"GoHome Hub Description",
		false,
		nil)

	/*
		//TODO: Fix
		tcp2 := tcp.(*Tcp600gwbDevice)
		tcp2.Token = "79tz3vbbop9pu5fcen60p97ix3mbvd3sblhjmz21"
		tcp2.Host = "https://192.168.0.23"
	*/

	// 192.168.0.24 / fluxbulb
	/*
				zoneID := "216438039298518643"
				z := &Zone{
					LocalID:     zoneID,
					GlobalID:    system.NextGlobalID(),
					Name:        "bulb1",
					Description: "tcp - bulb1",
					Device:      tcp,
					Type:        ZTLight,
					Output:      OTContinuous,
				}
				fmt.Println("BULB ID: " + z.GlobalID)
				tcp.Zones()[z.LocalID] = z
				system.AddZone(z)

				s := &Scene{
					LocalID:     "xxx",
					GlobalID:    system.NextGlobalID(),
					Name:        "Synthetic Scene",
					Description: "Scene to control lutron + tcp lights",
		            Managed: true,
					Commands: []Command{
						&ZoneSetLevelCommand{Zone: system.Zones["142"], Level: 30},
						&ZoneSetLevelCommand{Zone: system.Zones["153"], Level: 75},
					},
				}
				system.AddScene(s)*/

	//TODO:
	//1. Discover bulbs
	//2. Configure bulb
	//3. Control bulb

	//Aim: Be able to configure and control bulb completely from gohome app

	z := &zone.Zone{
		Address:     "192.168.0.24:5577",
		ID:          system.NextGlobalID(),
		Name:        "FluxBulb",
		Description: "Flux wifi bulb",
		DeviceID:    ghh.ID(),
		Type:        zone.ZTLight,
		Output:      zone.OTRGB,
		Controller:  zone.ZCFluxWIFI,
	}
	ghh.Zones()[z.Address] = z
	z2 := &zone.Zone{
		Address:     "192.168.0.24:55777",
		ID:          system.NextGlobalID(),
		Name:        "FluxBulb2",
		Description: "Flux wifi bulb",
		DeviceID:    ghh.ID(),
		Type:        zone.ZTLight,
		Output:      zone.OTRGB,
		Controller:  zone.ZCFluxWIFI,
	}
	ghh.Zones()[z2.Address] = z2

	system.AddDevice(ghh)
}
