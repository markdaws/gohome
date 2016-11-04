package belkin

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-home-iot/gossdp"
)

// Device contains information about a device that has been found on the network
type Device struct {
	Scan             ScanResponse
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
	Timeout          time.Duration
}

type root struct {
	Device *Device `xml:"device"`
}

// Load fetches all of the device specific information and updates the calling struct. The timeout
// parameter specifies how long to wait to connect and get a response before giving up
func (d *Device) Load(timeout time.Duration) error {
	// Example response
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

	client := http.Client{Timeout: timeout}
	resp, err := client.Get(d.Scan.Location)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("error fetching device info: %s", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response from device: %s", err)
	}

	var root root
	root.Device = d
	err = xml.Unmarshal(b, &root)
	if err != nil {
		return err
	}

	return nil
}

// TurnOn turns on the device.
func (d *Device) TurnOn(timeout time.Duration) error {
	location := parseLocation(d.Scan.Location)
	_, err := sendSOAP(
		timeout,
		location,
		"urn:Belkin:service:basicevent:1",
		"/upnp/control/basicevent1",
		"SetBinaryState",
		"<BinaryState>1</BinaryState>",
	)
	return err
}

// TurnOff turns off the device.
func (d *Device) TurnOff(timeout time.Duration) error {
	location := parseLocation(d.Scan.Location)
	_, err := sendSOAP(
		timeout,
		location,
		"urn:Belkin:service:basicevent:1",
		"/upnp/control/basicevent1",
		"SetBinaryState",
		"<BinaryState>0</BinaryState>",
	)
	return err
}

type attribute struct {
	Name  string `xml:"name"`
	Value int    `xml:"value"`
}

// FetchAttributes fetches the attributes of the device, such as switch state, sensor state etc
func (d *Device) FetchAttributes(timeout time.Duration) (*DeviceAttributes, error) {
	if d.Scan.SearchType != DTMaker {
		return nil, ErrUnsupportedAction
	}

	location := parseLocation(d.Scan.Location)
	body, err := sendSOAP(
		timeout,
		location,
		"urn:Belkin:service:deviceevent:1",
		"/upnp/control/deviceevent1",
		"GetAttributes",
		"",
	)

	if err != nil {
		return nil, err
	}

	// Response is double encoded
	body = html.UnescapeString(html.UnescapeString(body))

	/* Response looks like:
		<s:Envelope
	    	xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
		    s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
			<s:Body>
			    <u:GetAttributesResponse xmlns:u="urn:Belkin:service:deviceevent:1">
	    		<attributeList>
	    			<attribute>
					    <name>Switch</name>
						<value>1</value>
					</attribute>
					<attribute>
					    <name>Sensor</name>
	    				<value>1</value>
					</attribute>
					<attribute>
		    			<name>SwitchMode</name>
			    		<value>0</value>
					</attribute>
					<attribute>
				    	<name>SensorPresent</name>
					    <value>1</value>
					</attribute>
					</attributeList>
	    		</u:GetAttributesResponse>
			</s:Body>
		</s:Envelope>*/

	r := regexp.MustCompile(`(<attributeList>.*</attributeList>)`)
	attrList := r.FindStringSubmatch(body)
	if attrList == nil || len(attrList) == 0 {
		return nil, errors.New("attributeList element not found in response from device")
	}

	attrs := ParseAttributeList(attrList[0])
	return attrs, nil
}

// ParseBinaryState parses the BinaryState element that is received
// from some WeMo devices. It is expected that the body string looks like:
// <BinaryState>1|1477978435|0|0|0|1168438|0|100|0|0</BinaryState>
func ParseBinaryState(body string) *BinaryState {
	body = strings.TrimSuffix(strings.TrimPrefix(body, "<BinaryState>"), "</BinaryState>")
	if body == "" {
		return nil
	}

	values := strings.Split(body, "|")
	if len(values) < 2 {
		return nil
	}

	onOff, err := strconv.Atoi(values[0])
	if err != nil {
		onOff = 0
	}
	onSince, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		onSince = 0
	}
	state := &BinaryState{
		OnOff:   onOff,
		OnSince: onSince,
	}
	return state
}

// ParseAttributeList parses the xml attributeList response from the device e.g.
// <attributeList><attribute><name>...</name><value>...</value></attribute>...</attributeList>
// the body must be the open and close attributeList element.  If a valid input is not
// found the function will return nil
func ParseAttributeList(body string) *DeviceAttributes {

	attrs := struct {
		List []attribute `xml:"attribute"`
	}{}

	err := xml.Unmarshal([]byte(body), &attrs)
	if err != nil {
		return nil
	}

	if len(attrs.List) == 0 {
		return nil
	}

	var da DeviceAttributes
	for _, attr := range attrs.List {
		attr := attr
		switch attr.Name {
		case "Switch":
			da.Switch = &attr.Value
		case "Sensor":
			da.Sensor = &attr.Value
		case "SwitchMode":
			da.SwitchMode = &attr.Value
		case "SensorPresent":
			da.SensorPresent = &attr.Value
		}
	}
	return &da
}

// FetchBinaryState fetches the latest binary state value from the device
func (d *Device) FetchBinaryState(timeout time.Duration) (int, error) {
	// GetBinaryState always returns off for WeMo Maker, return error to callers
	if d.Scan.SearchType == DTMaker {
		return 0, ErrUnsupportedAction
	}

	location := parseLocation(d.Scan.Location)

	/* The response looks like below:
		<s:Envelope
		    xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
	    	s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
			<s:Body>
		    	<u:GetBinaryStateResponse xmlns:u="urn:Belkin:service:basicevent:1">
			    <BinaryState>8</BinaryState>
	    		</u:GetBinaryStateResponse>
		    </s:Body>
		</s:Envelope>*/
	body, err := sendSOAP(
		timeout,
		location,
		"urn:Belkin:service:basicevent:1",
		"/upnp/control/basicevent1",
		"GetBinaryState",
		"",
	)
	if err != nil {
		return 0, err
	}

	resp := struct {
		BinaryState int `xml:"Body>GetBinaryStateResponse>BinaryState"`
	}{}
	err = xml.Unmarshal([]byte(body), &resp)
	if err != nil {
		return 0, err
	}

	//Note 1 and 8 mean on
	return resp.BinaryState, nil
}

func sendSOAP(timeout time.Duration, location, serviceType, controlURL, action, body string) (string, error) {
	url := location + controlURL
	resp, err := postData(timeout, url, action, serviceType, body)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("error sending command: %s", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %s", err)
	}
	if resp.StatusCode != 200 {
		fmt.Errorf("non 200 response from device: %d, %s", resp.StatusCode, string(b))
	}
	return string(b), nil
}

func postData(timeout time.Duration, url, action, serviceType, body string) (*http.Response, error) {
	payload := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"utf-8\"?><s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\"><s:Body><u:%s xmlns:u=\"%s\">%s</u:%s></s:Body></s:Envelope>",
		action, serviceType, body, action,
	)

	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err)
	}
	req.Header.Add("SOAPACTION", "\""+serviceType+"#"+action+"\"")
	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	return client.Do(req)
}

func parseLocation(location string) string {
	return strings.Replace(location, "/setup.xml", "", -1)
}

type belkinListener struct {
	URN       string
	Responses *[]ScanResponse
}

func (l belkinListener) Response(m gossdp.ResponseMessage) {
	// example response
	// urn:Belkin:device:insight:1
	//{MaxAge:86400 SearchType:urn:Belkin:device:insight:1 DeviceId:Insight-1_0-231550K1200093 Usn:uuid:Insight-1_0-231550K1200093::urn:Belkin:device:insight:1 Location:http://10.22.22.1:49152/setup.xml Server:Unspecified, UPnP/1.0, Unspecified RawResponse:0xc208072120 Urn:urn:Belkin:device:insight:1}

	//urn:Belkin:service:basicevent:1
	//{MaxAge:86400 SearchType:urn:Belkin:service:basicevent:1 DeviceId:Insight-1_0-231550K1200093 Usn:uuid:Insight-1_0-231550K1200093::urn:Belkin:service:basicevent:1 Location:http://10.22.22.1:49152/setup.xml Server:Unspecified, UPnP/1.0, Unspecified RawResponse:0xc208072120 Urn:urn:Belkin:service:basicevent:1}

	if m.SearchType != l.URN {
		return
	}

	*l.Responses = append(*l.Responses, ScanResponse{
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
