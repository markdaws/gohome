package honeywell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Thermostat is an interface to a thermostat device
type Thermostat interface {
	// Connect connects to the mytotalconnectcomform service and authenticates the caller
	Connect(ctx context.Context, login, password string) error

	// FetchStatus fetches the current status of the device. NOTE: you must have called Connect()
	// first to authenticate the caller before calling this function, otherwise it will fail
	FetchStatus(ctx context.Context) (*Status, error)

	// CoolMode enables cool mode to the desired temp and duration. Pass 0 for period to run to next
	// schedule point
	CoolMode(ctx context.Context, temp float32, period time.Duration) error

	// HeatMode enables heat mode to the desired temp and duration. Pass 0 for period to run to next
	// schedule point
	HeatMode(ctx context.Context, temp float32, period time.Duration) error

	// FanMode switches between on and auto.  You can pass "on" or "auto" for mode
	FanMode(ctx context.Context, mode string) error

	// Cancel reverts the thermostat back to the schedule
	Cancel(ctx context.Context) error
}

// NewThermostat returns an initialzed thermostat instance.
// The deviceID has to be determined manually, log in to the mytotalconnectcomfort website,
// navigate to your device, then the URL will look something like
// https://mytotalconnectcomfort.com/portal/Device/CheckDataSession/123456, you need to copy the number
// that is in place of the 123456 and use that as your device ID.
func NewThermostat(deviceID int) Thermostat {
	return &thermostat{deviceID: deviceID}
}

type thermostat struct {
	deviceID int
	auth     *Auth
}

// Login logs in to the mytotalconnectcomfort service and returns an Auth instance
// that can then be used to fetch further information. The context object can be used
// to cancel a long running request. NOTE: Looks like the honeywell server returns a cookie
// that expires in 50years :/
func (t *thermostat) Connect(ctx context.Context, login, password string) error {

	form := url.Values{}
	form.Add("UserName", login)
	form.Add("Password", password)
	form.Add("timeOffset", "0")

	req, err := http.NewRequest("POST", "https://mytotalconnectcomfort.com/portal/", strings.NewReader(form.Encode()))
	if err != nil {
		return errors.Wrap(err, "unable to create request")
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req = req.WithContext(ctx)
	client := http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return errors.Wrap(err, "login request failed")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid respose code: %d", resp.StatusCode)
	}

	t.auth = &Auth{cookies: resp.Cookies()}
	return nil
}

// FetchStatus fetches the status of a device.  The deviceID has to be determined manually, log in
// to the mytotalconnectcomfort website, navigate to your device, then the URL will look something like
// https://mytotalconnectcomfort.com/portal/Device/CheckDataSession/123456, you need to copy the number
// that is in place of the 123456 and use that as your device ID.  The auth object comes from calling
// the Login() funtion previously
func (t *thermostat) FetchStatus(ctx context.Context) (*Status, error) {
	req, err := http.NewRequest(
		"GET",
		"https://mytotalconnectcomfort.com/portal/Device/CheckDataSession/"+strconv.Itoa(t.deviceID),
		nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	resp, err := t.doRequest(ctx, req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fetch failed, statuscode: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var status Status
	err = json.Unmarshal(b, &status)
	if err != nil {
		return nil, errors.Wrap(err, "invalid JSON")
	}

	return &status, nil
}

type state struct {
	CoolSetpoint   *float32 `json:"CoolSetpoint"`
	DeviceID       int      `json:"DeviceID"`
	FanMode        *int     `json:"FanMode", omitempty`
	HeatSetpoint   *float32 `json:"HeatSetpoint", omitempty`
	StatusCool     *int     `json:"StatusCool", omitempty`
	StatusHeat     *int     `json:"StatusHeat", omitempty`
	CoolNextPeriod *int     `json:"CoolNextPeriod", omitempty`
	HeatNextPeriod *int     `json:"HeatNextPeriod", omitempty`
}

// CoolMode sets the cool set point for the thermostat
func (t *thermostat) CoolMode(ctx context.Context, temp float32, period time.Duration) error {
	i := 1
	s := state{
		DeviceID:     t.deviceID,
		StatusCool:   &i,
		StatusHeat:   &i,
		CoolSetpoint: &temp,
	}

	if period != 0 {
		until := time.Now().Add(period)
		untilMin := int((until.Hour()*60 + until.Minute()) / 15)
		s.CoolNextPeriod = &untilMin
	}
	return t.setState(ctx, s)
}

// HeatMode sets the heat set point for the thermostat. Pass 0 for period if you want to
// set the temp until the next schedule period, period should be in 15 minute increments
func (t *thermostat) HeatMode(ctx context.Context, temp float32, period time.Duration) error {
	i := 1
	s := state{
		DeviceID:     t.deviceID,
		StatusCool:   &i,
		StatusHeat:   &i,
		HeatSetpoint: &temp,
	}

	if period != 0 {
		until := time.Now().Add(period)
		untilMin := int((until.Hour()*60 + until.Minute()) / 15)
		s.HeatNextPeriod = &untilMin
	}
	return t.setState(ctx, s)
}

// FanMode enables you to set mode=="on" or mode=="auto"
func (t *thermostat) FanMode(ctx context.Context, mode string) error {
	on := 1
	auto := 0

	s := state{DeviceID: t.deviceID}
	switch mode {
	case "on":
		s.FanMode = &on
	case "auto":
		s.FanMode = &auto
	default:
		return fmt.Errorf("invalid mode, permissable values 'on' or 'auto'")
	}

	return t.setState(ctx, s)
}

func (t *thermostat) Cancel(ctx context.Context) error {
	i := 0
	return t.setState(ctx, state{
		DeviceID:   t.deviceID,
		StatusCool: &i,
		StatusHeat: &i,
	})
}

func (t *thermostat) doRequest(ctx context.Context, r *http.Request) (*http.Response, error) {
	for _, c := range t.auth.cookies {
		// Server sends back multiple values for some cookies, with ones as blank, have to prune otherwise
		// we don't get valid values back
		if c.Value == "" {
			continue
		}
		r.AddCookie(c)
	}
	r.Header.Set("X-Requested-With", "XMLHttpRequest")

	r = r.WithContext(ctx)
	client := http.Client{}
	resp, err := client.Do(r)
	return resp, err
}

func (t *thermostat) setState(ctx context.Context, st state) error {

	b, err := json.Marshal(st)
	if err != nil {
		return errors.Wrap(err, "failed to serialize json request")
	}

	req, err := http.NewRequest(
		"POST",
		"https://mytotalconnectcomfort.com/portal/Device/SubmitControlScreenChanges",
		bytes.NewReader(b))

	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := t.doRequest(ctx, req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return errors.Wrap(err, "error making set request")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error setting state, status code: %d", resp.StatusCode)
	}

	b, err = ioutil.ReadAll(resp.Body)

	return nil
}

// Auth will be returned from a successful Login() call, you should use this
// to call functions such as FetchStatus() that require auth information
type Auth struct {
	cookies []*http.Cookie
}

// Status block returned from the server
type Status struct {
	Success           bool       `json:"success"`
	DeviceLive        bool       `json:"deviceLive"`
	CommunicationLost bool       `json:"communicationLost"`
	LatestData        LatestData `json:"latestData"`
}

// LatestData block returned from the server
type LatestData struct {
	HasFan                   bool    `json:"hasFan"`
	CanControlHumidification bool    `json:"canControlHumidification"`
	UIData                   UIData  `json:"uiData"`
	FanData                  FanData `json:"fanData"`
	DRData                   DRData  `json:"drData"`
}

// FanData block returned from the server
type FanData struct {
	FanMode                      float32 `json:"fanMode"`
	FanModeAutoAllowed           bool    `json:"fanModeAutoAllowed"`
	FanModeOnAllowed             bool    `json:"fanModeOnAllowed"`
	FanModeCirculateAllowed      bool    `json:"fanModeCirculateAllowed"`
	FanModeFollowScheduleAllowed bool    `json:"fanModeFollowScheduleAllowed"`
	FanIsRunning                 bool    `json:"fanIsRunning"`
}

// DRData block returned from the server
type DRData struct {
	CoolSetpLimit float32 `json:"CoolSetpLimit"`
	HeatSetpLimit float32 `json:"HeatSetpLimit"`
	Phase         float32 `json:"Phase"`
	OptOutable    bool    `json:"OptOutable"`
	DeltaCoolSP   float32 `json:"DeltaCoolSP"`
	DeltaHeatSP   float32 `json:"DeltaHeatSP"`
	Load          float32 `json:"Load"`
}

// UIData contains information returned from the honeywell service about a thermostat
type UIData struct {
	DispTemperature                  float32
	HeatSetpoint                     float32
	CoolSetpoint                     float32
	DisplayUnits                     string
	StatusHeat                       float32
	StatusCool                       float32
	HoldUntilCapable                 bool
	ScheduleCapable                  bool
	VacationHold                     float32
	DualSetpointStatus               bool
	HeatNextPeriod                   float32
	CoolNextPeriod                   float32
	HeatLowerSetptLimit              float32
	HeatUpperSetptLimit              float32
	CoolLowerSetptLimit              float32
	CoolUpperSetptLimit              float32
	ScheduleHeatSp                   float32
	ScheduleCoolSp                   float32
	SwitchAutoAllowed                bool
	SwitchCoolAllowed                bool
	SwitchOffAllowed                 bool
	SwitchHeatAllowed                bool
	SwitchEmergencyHeatAllowed       bool
	SystemSwitchPosition             float32
	Deadband                         float32
	IndoorHumidity                   float32
	DeviceID                         int
	Commercial                       bool
	DispTemperatureAvailable         bool
	IndoorHumiditySensorAvailable    bool
	IndoorHumiditySensorNotFault     bool
	VacationHoldUntilTime            float32
	TemporaryHoldUntilTime           float32
	IsInVacationHoldMode             bool
	VacationHoldCancelable           bool
	SetpointChangeAllowed            bool
	OutdoorTemperature               float32
	OutdoorHumidity                  float32
	OutdoorHumidityAvailable         bool
	OutdoorTemperatureAvailable      bool
	DispTemperatureStatus            float32
	IndoorHumidStatus                float32
	OutdoorTempStatus                float32
	OutdoorHumidStatus               float32
	OutdoorTemperatureSensorNotFault bool
	OutdoorHumiditySensorNotFault    bool
	CurrentSetpointStatus            float32
	EquipmentOutputStatus            float32
}

/*
Example response from the server

{
  "success": true,
  "deviceLive": true,
  "communicationLost": false,
  "latestData": {
    "uiData": {
      "DispTemperature": 65,
      "HeatSetpoint": 65,
      "CoolSetpoint": 0,
      "DisplayUnits": "F",
      "StatusHeat": 0,
      "StatusCool": 0,
      "HoldUntilCapable": true,
      "ScheduleCapable": true,
      "VacationHold": 0,
      "DualSetpointStatus": true,
      "HeatNextPeriod": 68,
      "CoolNextPeriod": 0,
      "HeatLowerSetptLimit": 40,
      "HeatUpperSetptLimit": 80,
      "CoolLowerSetptLimit": 50,
      "CoolUpperSetptLimit": 99,
      "ScheduleHeatSp": 65,
      "ScheduleCoolSp": 85,
      "SwitchAutoAllowed": false,
      "SwitchCoolAllowed": false,
      "SwitchOffAllowed": true,
      "SwitchHeatAllowed": true,
      "SwitchEmergencyHeatAllowed": false,
      "SystemSwitchPosition": 1,
      "Deadband": 0,
      "IndoorHumidity": 128,
      "DeviceID": 2196840,
      "Commercial": false,
      "DispTemperatureAvailable": true,
      "IndoorHumiditySensorAvailable": false,
      "IndoorHumiditySensorNotFault": true,
      "VacationHoldUntilTime": 0,
      "TemporaryHoldUntilTime": 0,
      "IsInVacationHoldMode": false,
      "VacationHoldCancelable": true,
      "SetpointChangeAllowed": true,
      "OutdoorTemperature": 128,
      "OutdoorHumidity": 128,
      "OutdoorHumidityAvailable": false,
      "OutdoorTemperatureAvailable": false,
      "DispTemperatureStatus": 0,
      "IndoorHumidStatus": 128,
      "OutdoorTempStatus": 128,
      "OutdoorHumidStatus": 128,
      "OutdoorTemperatureSensorNotFault": true,
      "OutdoorHumiditySensorNotFault": true,
      "CurrentSetpointStatus": 0,
      "EquipmentOutputStatus": 0
    },
    "fanData": {
      "fanMode": 0,
      "fanModeAutoAllowed": true,
      "fanModeOnAllowed": true,
      "fanModeCirculateAllowed": true,
      "fanModeFollowScheduleAllowed": false,
      "fanIsRunning": false
    },
    "hasFan": true,
    "canControlHumidification": false,
    "drData": {
      "CoolSetpLimit": 0,
      "HeatSetpLimit": 0,
      "Phase": -1,
      "OptOutable": false,
      "DeltaCoolSP": -0.01,
      "DeltaHeatSP": -0.01,
      "Load": 127.5
    }
  },
  "alerts": "\r\n\r\n"
}
*/
