// Package belkin provides support for Belkin devices, such as the WeMo Switch
package belkin

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fromkeith/gossdp"
)

// CREDIT: All the knowledge of how to control this product came from:
// https://github.com/timonreinhard/wemo-client

// DeviceType represents an identifier for the type of Belkin device you want to
// scan the network for
type DeviceType string

const (
	DTBridge      DeviceType = "urn:Belkin:device:bridge:1"
	DTSwitch                 = "urn:Belkin:device:controllee:1"
	DTMotion                 = "urn:Belkin:device:sensor:1"
	DTMaker                  = "urn:Belkin:device:Maker:1"
	DTInsight                = "urn:Belkin:device:insight:1"
	DTLightSwitch            = "urn:Belkin:device:lightswitch:1"
)

// Service contains information about a service exposed by the Belkin device
type Service struct {
	ServiceType string `xml:"serviceType"`
	ServiceID   string `xml:"serivceId"`
	ControlURL  string `xml:"controlURL"`
	EventSubURL string `xml:"eventSubURL"`
	SCPDURL     string `xml:"SCPDURL"`
}

// Device contains information about a device that has been found on the network
type Device struct {
	DeviceType       string    `xml:"deviceType"`
	FriendlyName     string    `xml:"friendlyName"`
	Manufacturer     string    `xml:"manufacturer"`
	ManufacturerURL  string    `xml:"manufacturerURL"`
	ModelDescription string    `xml:"modelDescription"`
	ModelName        string    `xml:"modelName"`
	ModelNumber      string    `xml:"modelNumber"`
	ModelURL         string    `xml:"modelURL"`
	SerialNumber     string    `xml:"serialNumber"`
	UDN              string    `xml:"UDN"`
	UPC              string    `xml:"UPC"`
	MACAddress       string    `xml:"macAddress"`
	FirmwareVersion  string    `xml:"firmwareVersion"`
	IconVersion      string    `xml:"iconVersion"`
	BinaryState      int       `xml:"binaryState"`
	ServiceList      []Service `xml:"serviceList>service"`
}

type root struct {
	Device Device `xml:"device"`
}

// ScanResponse contains information from a device that responded to a scan response
type ScanResponse struct {
	MaxAge     int
	SearchType string
	DeviceID   string
	USN        string
	Location   string
	Server     string
	URN        string
}

// Scan detects Belkin devices on the network
func Scan(dt DeviceType, waitTimeSeconds int) ([]ScanResponse, error) {
	responses := make([]ScanResponse, 0)
	l := belkinListener{
		URN:       string(dt),
		Responses: &responses,
	}

	c, err := gossdp.NewSsdpClientWithLogger(l, l)
	if err != nil {
		return nil, fmt.Errorf("failed to start ssdp discovery client: %s", err)
	}

	defer c.Stop()
	go c.Start()
	err = c.ListenFor(string(dt))
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %s", err)
	}

	time.Sleep(time.Duration(waitTimeSeconds) * time.Second)
	return responses, nil
}

func LoadDevice(scanResponse ScanResponse) (*Device, error) {
	client := http.Client{}
	resp, err := client.Get(scanResponse.Location)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching device info: %s", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response from device: %s", err)
	}

	var root root
	err = xml.Unmarshal(b, &root)
	if err != nil {
		return nil, err
	}

	return &root.Device, nil
}

func TurnOn(location string) error {
	return SendSOAP(
		location,
		"urn:Belkin:service:basicevent:1",
		"/upnp/control/basicevent1",
		"SetBinaryState",
		"<BinaryState>1</BinaryState>",
	)
}

func TurnOff(location string) error {
	return SendSOAP(
		location,
		"urn:Belkin:service:basicevent:1",
		"/upnp/control/basicevent1",
		"SetBinaryState",
		"<BinaryState>0</BinaryState>",
	)
}

func SendSOAP(location, serviceType, controlURL, action, body string) error {
	url := location + controlURL
	resp, err := postData(url, action, serviceType, body)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("error sending command: %s", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %s", err)
	}
	fmt.Println(string(b))
	if resp.StatusCode != 200 {
		fmt.Errorf("non 200 response from device: %d, %s", resp.StatusCode, string(b))
	}
	return nil
}

func postData(url, action, serviceType, body string) (*http.Response, error) {
	payload := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"utf-8\"?><s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\"><s:Body><u:%s xmlns:u=\"%s\">%s</u:%s></s:Body></s:Envelope>",
		action, serviceType, body, action,
	)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err)
	}
	req.Header.Add("SOAPACTION", "\""+serviceType+"#"+action+"\"")
	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	return client.Do(req)
}

type belkinListener struct {
	URN       string
	Responses *[]ScanResponse
}

func (t belkinListener) Response(m gossdp.ResponseMessage) {
	// example response
	// urn:Belkin:device:insight:1
	//{MaxAge:86400 SearchType:urn:Belkin:device:insight:1 DeviceId:Insight-1_0-231550K1200093 Usn:uuid:Insight-1_0-231550K1200093::urn:Belkin:device:insight:1 Location:http://10.22.22.1:49152/setup.xml Server:Unspecified, UPnP/1.0, Unspecified RawResponse:0xc208072120 Urn:urn:Belkin:device:insight:1}

	//urn:Belkin:service:basicevent:1
	//{MaxAge:86400 SearchType:urn:Belkin:service:basicevent:1 DeviceId:Insight-1_0-231550K1200093 Usn:uuid:Insight-1_0-231550K1200093::urn:Belkin:service:basicevent:1 Location:http://10.22.22.1:49152/setup.xml Server:Unspecified, UPnP/1.0, Unspecified RawResponse:0xc208072120 Urn:urn:Belkin:service:basicevent:1}

	if m.SearchType != t.URN {
		return
	}

	*t.Responses = append(*t.Responses, ScanResponse{
		MaxAge:     m.MaxAge,
		SearchType: m.SearchType,
		DeviceID:   m.DeviceId,
		USN:        m.Usn,
		Location:   m.Location,
		Server:     m.Server,
		URN:        m.Urn,
	})
}
func (l belkinListener) Tracef(fmt string, args ...interface{}) {}
func (l belkinListener) Infof(fmt string, args ...interface{})  {}
func (l belkinListener) Warnf(fmt string, args ...interface{})  {}
func (l belkinListener) Errorf(fmt string, args ...interface{}) {}

/*
Response from querying the location address of the insight device
<?xml version="1.0"?>
<root xmlns="urn:Belkin:device-1-0">
  <specVersion>
    <major>1</major>
    <minor>0</minor>
  </specVersion>
  <device>
    <deviceType>urn:Belkin:device:insight:1</deviceType>
    <friendlyName>WeMo Insight</friendlyName>
    <manufacturer>Belkin International Inc.</manufacturer>
    <manufacturerURL>http://www.belkin.com</manufacturerURL>
    <modelDescription>Belkin Insight 1.0</modelDescription>
    <modelName>Insight</modelName>
    <modelNumber>1.0</modelNumber>
    <modelURL>http://www.belkin.com/plugin/</modelURL>
    <serialNumber>231550K1200093</serialNumber>
    <UDN>uuid:Insight-1_0-231550K1200093</UDN>
    <UPC>123456789</UPC>
    <macAddress>94103ECFA7FA</macAddress>
    <firmwareVersion>WeMo_WW_2.00.9213.PVT-OWRT-InsightV2</firmwareVersion>
    <iconVersion>0|49152</iconVersion>
    <binaryState>0</binaryState>
    <iconList>
      <icon>
        <mimetype>jpg</mimetype>
        <width>100</width>
        <height>100</height>
        <depth>100</depth>
         <url>icon.jpg</url>
      </icon>
    </iconList>
    <serviceList>
      <service>
        <serviceType>urn:Belkin:service:WiFiSetup:1</serviceType>
        <serviceId>urn:Belkin:serviceId:WiFiSetup1</serviceId>
        <controlURL>/upnp/control/WiFiSetup1</controlURL>
        <eventSubURL>/upnp/event/WiFiSetup1</eventSubURL>
        <SCPDURL>/setupservice.xml</SCPDURL>
      </service>
      <service>
        <serviceType>urn:Belkin:service:timesync:1</serviceType>
        <serviceId>urn:Belkin:serviceId:timesync1</serviceId>
        <controlURL>/upnp/control/timesync1</controlURL>
        <eventSubURL>/upnp/event/timesync1</eventSubURL>
        <SCPDURL>/timesyncservice.xml</SCPDURL>
      </service>
      <service>
        <serviceType>urn:Belkin:service:basicevent:1</serviceType>
        <serviceId>urn:Belkin:serviceId:basicevent1</serviceId>
        <controlURL>/upnp/control/basicevent1</controlURL>
        <eventSubURL>/upnp/event/basicevent1</eventSubURL>
        <SCPDURL>/eventservice.xml</SCPDURL>
      </service>
      <service>
        <serviceType>urn:Belkin:service:firmwareupdate:1</serviceType>
        <serviceId>urn:Belkin:serviceId:firmwareupdate1</serviceId>
        <controlURL>/upnp/control/firmwareupdate1</controlURL>
        <eventSubURL>/upnp/event/firmwareupdate1</eventSubURL>
        <SCPDURL>/firmwareupdate.xml</SCPDURL>
      </service>
      <service>
        <serviceType>urn:Belkin:service:rules:1</serviceType>
        <serviceId>urn:Belkin:serviceId:rules1</serviceId>
        <controlURL>/upnp/control/rules1</controlURL>
        <eventSubURL>/upnp/event/rules1</eventSubURL>
        <SCPDURL>/rulesservice.xml</SCPDURL>
      </service>

      <service>
        <serviceType>urn:Belkin:service:metainfo:1</serviceType>
        <serviceId>urn:Belkin:serviceId:metainfo1</serviceId>
        <controlURL>/upnp/control/metainfo1</controlURL>
        <eventSubURL>/upnp/event/metainfo1</eventSubURL>
        <SCPDURL>/metainfoservice.xml</SCPDURL>
      </service>

      <service>
        <serviceType>urn:Belkin:service:remoteaccess:1</serviceType>
        <serviceId>urn:Belkin:serviceId:remoteaccess1</serviceId>
        <controlURL>/upnp/control/remoteaccess1</controlURL>
        <eventSubURL>/upnp/event/remoteaccess1</eventSubURL>
        <SCPDURL>/remoteaccess.xml</SCPDURL>
      </service>

      <service>
        <serviceType>urn:Belkin:service:deviceinfo:1</serviceType>
        <serviceId>urn:Belkin:serviceId:deviceinfo1</serviceId>
        <controlURL>/upnp/control/deviceinfo1</controlURL>
        <eventSubURL>/upnp/event/deviceinfo1</eventSubURL>
        <SCPDURL>/deviceinfoservice.xml</SCPDURL>
      </service>

      <service>
        <serviceType>urn:Belkin:service:insight:1</serviceType>
        <serviceId>urn:Belkin:serviceId:insight1</serviceId>
        <controlURL>/upnp/control/insight1</controlURL>
        <eventSubURL>/upnp/event/insight1</eventSubURL>
        <SCPDURL>/insightservice.xml</SCPDURL>
      </service>

      <service>
        <serviceType>urn:Belkin:service:smartsetup:1</serviceType>
        <serviceId>urn:Belkin:serviceId:smartsetup1</serviceId>
        <controlURL>/upnp/control/smartsetup1</controlURL>
        <eventSubURL>/upnp/event/smartsetup1</eventSubURL>
        <SCPDURL>/smartsetup.xml</SCPDURL>
      </service>

      <service>
        <serviceType>urn:Belkin:service:manufacture:1</serviceType>
        <serviceId>urn:Belkin:serviceId:manufacture1</serviceId>
        <controlURL>/upnp/control/manufacture1</controlURL>
        <eventSubURL>/upnp/event/manufacture1</eventSubURL>
        <SCPDURL>/manufacture.xml</SCPDURL>
      </service>

    </serviceList>
   <presentationURL>/pluginpres.html</presentationURL>
</device>
</root>
*/
