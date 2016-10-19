package connectedbytcp

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-home-iot/gossdp"
	"github.com/go-home-iot/gouuid"
)

// CREDIT: information on how to communicate with this device came from:
// https://github.com/stockmopar/connectedbytcp

// ErrUnauthorized represents an error when the user tried to call an
// API but were not authorized to do so. This can occur if you try to
// call GetToken without pressing the "sync" button on your physical hub
// or if you try to make API calls without a valid token
var ErrUnauthorized = errors.New("unauthorized")

var rootCmd = "cmd=%s&data=%s&fmt=xml"

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

// Discover returns the address e.g. https://192.168.0.23 of the ConnectByTCP Hub if
// one was found on the network.
func Scan(waitTimeSeconds int) ([]ScanResponse, error) {
	URN := "urn:greenwavereality-com:service:gop:1"
	var responses []ScanResponse
	l := tcpListener{
		URN:       URN,
		Responses: &responses,
	}

	c, err := gossdp.NewSsdpClientWithLogger(l, l)
	if err != nil {
		return nil, fmt.Errorf("failed to start ssdp discovery client: %s", err)
	}

	defer c.Stop()
	go c.Start()
	err = c.ListenFor(URN)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %s", err)
	}

	time.Sleep(time.Duration(waitTimeSeconds) * time.Second)
	return responses, nil
}

// GetToken returns the security token required to make any API calls to the
// ConnectedByTCP hub. In order for this function to succeed, you must press
// the physical "sync" button on the hub before calling this function. Calling
// this function without first pressing the sync button will cause a
// ErrUnauthorized error to be returned.  The address field should be in the
// format "https://192.168.0.23" for example.
func GetToken(address string) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %s", err)
	}

	data := fmt.Sprintf("<gip><version>1</version><email>%s</email><password>%s</password></gip>", id, id)
	resp, err := postData(address, "GWRLogin", data)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("error fetching token: %s", err)
	}

	// example responses
	// OK - <gip><version>1</version><rc>200</rc><token>xyzaqlifpzoo7lao56xoy3m0pu3wsy1n4dnzobkj</token></gip>
	// Not synced - <gip><version>1</version><rc>404</rc></gip>
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading token response: %s", err)
	}

	rBody := string(b)
	if strings.Contains(rBody, "<rc>401</rc>") ||
		strings.Contains(rBody, "<rc>404</rc>") {
		return "", ErrUnauthorized
	}

	type response struct {
		Token string `xml:"token"`
	}
	var r response
	err = xml.Unmarshal([]byte(rBody), &r)
	if err != nil {
		return "", fmt.Errorf("error reading xml response: %s", err)
	}

	if r.Token == "" {
		return "", fmt.Errorf("token not found, empty")
	}
	return r.Token, nil
}

// Device represents the device information returned by the ConnectedByTCP hub
// in the response to the RoomGetCarousel call
type Device struct {
	DID        string  `xml:"did"`
	Known      float64 `xml:"known"`
	State      float64 `xml:"state"`
	Offline    float64 `xml:"offline"`
	Node       float64 `xml:"node"`
	Port       float64 `xml:"port"`
	NodeType   float64 `xml:"nidetype"`
	Name       string  `xml:"name"`
	ColorID    float64 `xml:"colorid"`
	Type       string  `xml:"type"`
	RangeMin   float64 `xml:"rangemin"`
	Power      float64 `xml:"power"`
	PowerAvg   float64 `xml:"poweravg"`
	Energy     float64 `xml:"energy"`
	Score      float64 `xml:"score"`
	ProductID  float64 `xml:"productid"`
	ProdBrand  string  `xml:"prodbrand"`
	ProdModel  string  `xml:"prodmodel"`
	ProdType   string  `xml:"prodtype"`
	ProdTypeID float64 `xml:"prodtypeid"`
	ClassID    float64 `xml:"classid"`
}

// Room represents the room information returned by the ConnectedByTCP hub in
// response to the RoomGetCarousel call
type Room struct {
	Name        string   `xml:"name"`
	Description string   `xml:"desc"`
	Known       float64  `xml:"known"`
	Type        float64  `xml:"type"`
	Color       string   `xml:"color"`
	ColorID     float64  `xml:"colorid"`
	Image       string   `xml:"img"`
	Power       float64  `xml:"power"`
	PowerAVG    float64  `xml:"poweravg"`
	Energy      float64  `xml:"energy"`
	Devices     []Device `xml:"device"`
}

// GIP represents the gip element returned by the ConnectedByTCP hub in response
// to the RoomGetCarousel call
type GIP struct {
	Rooms []Room `xml:"gwrcmd>gdata>gip>room"`
}

// RoomGetCarousel return the room and device information from the ConnectedByTCP
// hub.  Call this to get room and bulb information, such as device IDs, names etc.
// Address should be in the form "https://192.168.0.23" for example and the token
// value should be the token you retried from the GetToken function call
func RoomGetCarousel(address, token string) (*GIP, error) {
	data := fmt.Sprintf("<gwrcmds><gwrcmd><gcmd>RoomGetCarousel</gcmd><gdata><gip><version>1</version><token>%s</token><fields>name,control,power,product,class,realtype,status</fields></gip></gdata></gwrcmd></gwrcmds>", token)
	resp, err := postData(address, "GWRBatch", data)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching carousel: %s", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	// Example response from the hub
	// "<gwrcmds><gwrcmd><gcmd>RoomGetCarousel</gcmd><gdata><gip><version>1</version><rc>200</rc><room><rid>0</rid><name>Kitchen</name><desc></desc><known>1</known><type>0</type><color>000000</color><colorid>0</colorid><img>images/black.png</img><power>0</power><poweravg>0</poweravg><energy>0</energy><device><did>216438039298518643</did><known>1</known><lock>0</lock><state>0</state><level>100</level><node>64</node><port>0</port><nodetype>16386</nodetype><name>Bulb1</name><desc>LED</desc><colorid>0</colorid><type>multilevel</type><rangemin>0</rangemin><rangemax>99</rangemax><power>0</power><poweravg>0</poweravg><energy>0</energy><score>0</score><productid>1</productid><prodbrand>TCP</prodbrand><prodmodel>LED A19 11W</prodmodel><prodtype>LED</prodtype><prodtypeid>78</prodtypeid><classid>2</classid><class></class><subclassid>1</subclassid><subclass></subclass><other><rcgroup></rcgroup><manufacturer>TCP</manufacturer><capability>productinfo,identify,meter_power,switch_binary,switch_multilevel</capability><bulbpower>11</bulbpower></other></device></room></gip></gdata></gwrcmd></gwrcmds>"
	var g GIP
	err = xml.Unmarshal(b, &g)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response xml: %s", err)
	}
	return &g, nil
}

// VerifyConnection connects to the hub and tries to perform an API call to verify that the
// supplied parameters are correct. If no error is returned then the call was successful and
// the address and token are valid values
func VerifyConnection(address, token string) error {
	//See if we can get some status successfully
	_, err := RoomGetCarousel(address, token)
	return err
}

// TurnOn turns the bulb on
func TurnOn(hubAddress, zoneAddress, token string) error {
	return SetLevel(hubAddress, zoneAddress, token, 1)
}

// TurnOff turns the bulb off
func TurnOff(hubAddress, zoneAddress, token string) error {
	return SetLevel(hubAddress, zoneAddress, token, 0)
}

// SetLevel sets the bulb to the specified level
func SetLevel(hubAddress, zoneAddress, token string, level int32) error {
	if level == 0 {
		return setLevel(hubAddress, zoneAddress, token, 0)
	} else if level == 1 {
		return setLevel(hubAddress, zoneAddress, token, 1)
	} else {
		// 0 -> off, 1 -> on, if the light was set to 0 then you have to set a 1 first
		// before trying to set any other level
		err := setLevel(hubAddress, zoneAddress, token, 1)
		if err != nil {
			return err
		}
		return setLevel(hubAddress, zoneAddress, token, level)
	}
}

func setLevel(hubAddress, zoneAddress, token string, level int32) error {
	var data string
	if level == 0 || level == 1 {
		data = "<gip><version>1</version><token>%s</token><did>%s</did><value>%d</value></gip>"
	} else {
		data = "<gip><version>1</version><token>%s</token><did>%s</did><value>%d</value><type>level</type></gip>"
	}
	data = fmt.Sprintf(data, token, zoneAddress, level)
	data = fmt.Sprintf("<gwrcmds><gwrcmd><gcmd>DeviceSendCommand</gcmd><gdata>%s</gdata></gwrcmd></gwrcmds>", data)

	resp, err := postData(hubAddress, "GWRBatch", data)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("failed to send command: %s", err)
	}

	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error sending command: %s", string(b))
		}
	}
	return nil
}

func postData(address, command, data string) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr}
	cmd := fmt.Sprintf(rootCmd, command, data)
	return client.Post(address+"/gwr/gpo.php", "text/xml; charset=\"utf-8\"", bytes.NewReader([]byte(cmd)))
}

type tcpListener struct {
	URN       string
	Responses *[]ScanResponse
}

func (t tcpListener) Response(m gossdp.ResponseMessage) {
	// example response
	// {MaxAge:7200 SearchType:urn:greenwavereality-com:service:gop:1 DeviceId:71403833960916 Usn:uuid:71403833960916::urn:greenwavereality-com:service:gop:1 Location:https://192.168.0.23 Server:linux UPnP/1.1 Apollo3/3.0.74 RawResponse:0xc2080305a0 Urn:urn:greenwavereality-com:service:gop:1}
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
	})
}
func (l tcpListener) Tracef(fmt string, args ...interface{}) {}
func (l tcpListener) Infof(fmt string, args ...interface{})  {}
func (l tcpListener) Warnf(fmt string, args ...interface{})  {}
func (l tcpListener) Errorf(fmt string, args ...interface{}) {}
