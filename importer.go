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

	var makeDevice = func(modelNumber string, deviceMap map[string]interface{}, sys *System, stream bool, ci comm.ConnectionInfo) Device {
		var deviceID string = strconv.FormatFloat(deviceMap["ID"].(float64), 'f', 0, 64)
		var deviceName string = deviceMap["Name"].(string)

		device := NewDevice(
			modelNumber,
			deviceID,
			sys.NextGlobalID(),
			deviceName,
			deviceName,
			stream,
			sys,
			cmdProcessor,
			ci)

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
					LocalID:     buttonID,
					GlobalID:    globalID,
					Name:        buttonName,
					Description: buttonName,
					Commands: []cmd.Command{
						&cmd.ButtonPress{
							ButtonAddress:  btn.Address,
							ButtonID:       btn.ID,
							DeviceName:     sbp.Name(),
							DeviceLocalID:  sbp.LocalID(),
							DeviceGlobalID: sbp.GlobalID(),
						},
						&cmd.ButtonRelease{
							ButtonAddress:  btn.Address,
							ButtonID:       btn.ID,
							DeviceName:     sbp.Name(),
							DeviceLocalID:  sbp.LocalID(),
							DeviceGlobalID: sbp.GlobalID(),
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

		var deviceID string = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			//ModelNumber: L-BDGPRO2-WH
			sbp = makeDevice("L-BDGPRO2-WH", device, system, true, &comm.TelnetConnectionInfo{
				Network:       "tcp",
				Address:       "192.168.0.10:23",
				Login:         "lutron",
				Password:      "integration",
				PoolSize:      2,
				Authenticator: sbp,
			})
			sbp.ConnectionInfo().(*comm.TelnetConnectionInfo).Authenticator = sbp
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
		var deviceID string = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			continue
		}
		gohomeDevice := makeDevice("", device, system, false, nil)
		system.AddDevice(gohomeDevice)
		sbp.Devices()[gohomeDevice.LocalID()] = gohomeDevice
	}

	zones, ok := root["Zones"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Zones key")
	}

	fmt.Println("\nZONES")
	for _, zoneMap := range zones {
		zone := zoneMap.(map[string]interface{})
		fmt.Printf("%d: %s\n", int(zone["ID"].(float64)), zone["Name"])

		var zoneID string = strconv.FormatFloat(zone["ID"].(float64), 'f', 0, 64)
		var zoneName string = zone["Name"].(string)
		var zoneTypeFinal ZoneType = ZTLight
		if zoneType, ok := zone["Type"].(string); ok {
			switch zoneType {
			case "light":
				zoneTypeFinal = ZTLight
			case "shade":
				zoneTypeFinal = ZTShade
			}
		}
		var outputTypeFinal OutputType = OTContinuous
		if outputType, ok := zone["Output"].(string); ok {
			switch outputType {
			case "binary":
				outputTypeFinal = OTBinary
			case "continuous":
				outputTypeFinal = OTContinuous
			}
		}
		z := &Zone{
			Address:     zoneID,
			ID:          system.NextGlobalID(),
			Name:        zoneName,
			Description: zoneName,
			Device:      sbp,
			Type:        zoneTypeFinal,
			Output:      outputTypeFinal,
		}
		system.AddZone(z)
		sbp.Zones()[z.Address] = z
	}

	//TODO: Move
	importConnectedByTCP(system)
	importGoHomeHub(system)
	return system, nil
}

//TODO: Temp function - import from UI
func importConnectedByTCP(system *System) {
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
	tcp := NewDevice(
		"TCP600GWB",
		"tcphub",
		system.NextGlobalID(),
		"ConnectedByTcp Hub",
		"Description",
		false,
		system,
		//TODO: Remove from NewDevice
		system.CmdProcessor,
		nil)

	/*
		//TODO: Fix
		tcp2 := tcp.(*Tcp600gwbDevice)
		tcp2.Token = "79tz3vbbop9pu5fcen60p97ix3mbvd3sblhjmz21"
		tcp2.Host = "https://192.168.0.23"
	*/

	zoneID := "216438039298518643"
	z := &Zone{
		Address:     zoneID,
		ID:          system.NextGlobalID(),
		Name:        "bulb1",
		Description: "tcp - bulb1",
		Device:      tcp,
		Type:        ZTLight,
		Output:      OTContinuous,
		Controller:  "TCP - LED A19 11W",
	}
	fmt.Println("BULB ID: " + z.ID)
	tcp.Zones()[z.Address] = z
	system.AddZone(z)
	system.AddDevice(tcp)

	z1 := system.Zones["142"]
	z2 := system.Zones["153"]
	s := &Scene{
		LocalID:     "xxx",
		GlobalID:    system.NextGlobalID(),
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
	ti := &comm.TelnetConnectionInfo{}
	ghh := NewDevice(
		"GoHomeHub",
		"gohomehub",
		system.NextGlobalID(),
		"GoHome Hub",
		"GoHome Hub Description",
		false,
		system,
		//TODO: Remove from NewDevice
		system.CmdProcessor,
		ti)

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

	z := &Zone{
		Address:     "192.168.0.24:5577",
		ID:          system.NextGlobalID(),
		Name:        "FluxBulb",
		Description: "Flux wifi bulb",
		Device:      ghh,
		Type:        ZTLight,
		Output:      OTRGB,
		Controller:  ZCFluxWIFI,
	}
	ghh.Zones()[z.Address] = z
	z2 := &Zone{
		Address:     "192.168.0.24:55777",
		ID:          system.NextGlobalID(),
		Name:        "FluxBulb2",
		Description: "Flux wifi bulb",
		Device:      ghh,
		Type:        ZTLight,
		Output:      OTRGB,
		Controller:  ZCFluxWIFI,
	}
	ghh.Zones()[z2.Address] = z2

	system.AddDevice(ghh)

	//TODO: Remove
	c := comm.NewTelnetConnection(comm.TelnetConnectionInfo{
		PoolSize: 2,
		Network:  "tcp",
		Address:  "192.168.0.24:5577",
	})
	err := c.Open()
	if err != nil {
		fmt.Printf("ERROR CONNECTING: %s", err)
	} else {
		/*
			system.CmdProcessor.Enqueue(&cmd.ZoneSetLevel{
				ZoneLocalID:  z.LocalID,
				ZoneGlobalID: z.GlobalID,
				ZoneName:     z.Name,
				Level:        cmd.Level{R: 255, G: 0, B: 0},
			})

			b := []byte{0x31, 0x00, 0x00, 0xff, 0x00, 0xf0, 0x0f}
			var t int = 0
			for _, v := range b {
				t += int(v)
			}
			cs := t & 0xff
			b = append(b, byte(cs))
			n, err := c.Write(b)
			if err != nil {
				fmt.Printf("ERROR SENDING %s\n", err)
			} else {
				fmt.Printf("Send data: %d\n", n)
			}

			_ = n
			//		c.Close()
		*/
	}
}
