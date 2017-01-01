package gohome

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/go-yaml/yaml"
	"github.com/markdaws/gohome/pkg/attr"
	"github.com/markdaws/gohome/pkg/clock"
	"github.com/markdaws/gohome/pkg/cmd"
	"github.com/markdaws/gohome/pkg/feature"
	"github.com/markdaws/gohome/pkg/log"

	"github.com/Machiel/slugify"
)

// automationSys provides access to system information, just enough for the automation type to function
type automationSys interface {
	NewID() string
	SceneByID(ID string) *Scene
	FeaturesByType(featureType string) map[string]*feature.Feature
	FeatureByID(ID string) *feature.Feature
	FeatureByAID(AID string) *feature.Feature
}

// Automation represents an automation instance. Each piece of automation has a trigger which is a set
// of conditions which when evaluating to true cause the automation actions to execute
type Automation struct {
	Name    string
	TempID  string
	Enabled bool
	Trigger Trigger
	evtbus.Consumer
	Triggered func(actions *CommandGroup)
}

func (a *Automation) ConsumerName() string {
	return fmt.Sprintf("automation - %s", a.Name)
}

func (a *Automation) StartConsuming(ch chan evtbus.Event) {
	a.Trigger.StartConsuming(ch)
}

func (a *Automation) StopConsuming() {
	a.Trigger.StopConsuming()
}

type condition struct {
	AttrLocalID *string     `yaml:"attr"`
	Op          *string     `yaml:"op"`
	Value       interface{} `yaml:"value"`
	feature     *feature.Feature
	And         []*condition `yaml:"and"`
	Or          []*condition `yaml:"or"`
}

func (c *condition) Evaluate(e *FeatureAttrsChangedEvt) bool {
	if c.feature.ID != e.FeatureID {
		return false
	}

	attribute, ok := e.Attrs[*c.AttrLocalID]
	if !ok {
		return false
	}

	switch attribute.DataType {
	case attr.DTInt32:
		a := attribute.Value.(int32)
		b := toInt32(c.Value)
		if b == nil {
			return false
		}
		switch *c.Op {
		case "<":
			return a < *b
		case ">":
			return a > *b
		case "==":
			return a == *b
		case "!=":
			return a != *b
		case "<=":
			return a <= *b
		case ">=":
			return a >= *b
		}

	case attr.DTFloat32:
		a := attribute.Value.(float32)
		b := toFloat32(c.Value)
		if b == nil {
			return false
		}
		switch *c.Op {
		case "<":
			return a < *b
		case ">":
			return a > *b
		case "==":
			return a == *b
		case "!=":
			return a != *b
		case "<=":
			return a <= *b
		case ">=":
			return a >= *b
		}

	case attr.DTString:
		a := attribute.Value.(string)
		b := c.Value.(string)
		switch *c.Op {
		case "<":
			return a < b
		case ">":
			return a > b
		case "==":
			return a == b
		case "!=":
			return a != b
		case "<=":
			return a <= b
		case ">=":
			return a >= b
		}

	case attr.DTBool:
		switch *c.Op {
		case "==":
			return attribute.Value == c.Value
		case "!=":
			return attribute.Value != c.Value
		}
	}

	return false
}

// helper type to deserialize the yaml in to our internal object model
type automationIntermediate struct {
	Name    string `yaml:"name"`
	Enabled *bool  `yaml:"enabled"`
	Trigger *struct {
		Time *struct {
			At   string `yaml:"at"`
			Days string `yaml:"days"`
		} `yaml:"time"`
		Feature *struct {
			ID        *string    `yaml:"id"`
			AID       *string    `yaml:"aid"`
			Condition *condition `yaml:"condition"`
			Count     int        `yaml:"count"`
			Duration  int        `yaml:"duration"`
		} `yaml:"feature"`
	} `yaml:"trigger"`
	Actions []struct {
		Scene *struct {
			ID string `yaml:"id"`
		} `yaml:"scene"`
		LightZone *struct {
			ID         *string  `yaml:"id"`
			AID        *string  `yaml:"aid"`
			OnOff      *string  `yaml:"on_off"`
			Brightness *float64 `yaml:"brightness"`
		} `yaml:"light_zone"`
		Outlet *struct {
			ID    *string `yaml:"id"`
			AID   *string `yaml:"aid"`
			OnOff *string `yaml:"on_off"`
		} `yaml:"outlet"`
		Switch *struct {
			ID    *string `yaml:"id"`
			AID   *string `yaml:"aid"`
			OnOff *string `yaml:"on_off"`
		} `yaml:"switch"`
		WindowTreatment *struct {
			ID         *string  `yaml:"id"`
			AID        *string  `yaml:"aid"`
			OpenClosed *string  `yaml:"open_closed"`
			Offset     *float64 `yaml:"offset"`
		} `yaml:"window_treatment"`
		HeatZone *struct {
			ID         *string  `yaml:"id"`
			AID        *string  `yaml:"aid"`
			TargetTemp *float64 `yaml:"target_temp"`
		} `yaml:"heat_zone"`
	} `yaml:"actions"`
}

// LoadAutomation loads all of the automation files from the specified path
func LoadAutomation(sys automationSys, path string) (map[string]*Automation, error) {

	autos := make(map[string]*Automation)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate automation files: %s", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}
		fullPath := path + "/" + file.Name()
		b, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.E("automation - failed to read contents of automation file: %s", fullPath)
			continue
		}

		auto, err := NewAutomation(sys, string(b))
		if err != nil {
			log.E("automation - failed to create automation: %s, %s", fullPath, err)
			continue
		}
		autos[auto.Name] = auto
	}
	return autos, nil
}

// NewAutomation creates a new automation instance
func NewAutomation(sys automationSys, config string) (*Automation, error) {

	var auto automationIntermediate
	err := yaml.Unmarshal([]byte(config), &auto)

	if err != nil {
		return nil, err
	}

	if auto.Name == "" {
		return nil, fmt.Errorf("missing name key, all automation must have a name")
	}

	if auto.Trigger == nil {
		return nil, fmt.Errorf("missing trigger key, trigger must be defined")
	}

	if len(auto.Actions) == 0 {
		return nil, fmt.Errorf("missing actions key, no actions specified")
	}

	// Defaults to true if not present
	if auto.Enabled == nil {
		auto.Enabled = new(bool)
		*auto.Enabled = true
	}

	// We parse the actions, but don't use them here, they are built just to make sure
	// there are no syntax issues, but we generate commands at the time the trigger
	// fires to make sure we have the latest state. For example if the user has written
	// automation to turn all lights off, if we generate the commands on load and then
	// add a new light, we would need to update the automation, by deferring the command
	// generation to the point of execution we mitigate this issue
	_, err = parseActions(sys, auto)
	if err != nil {
		return nil, err
	}

	finalAuto := &Automation{
		Name: auto.Name,

		// Automation doesn't have a permanent ID since we load them from files each time the
		// system starts, we give it a temp ID so that the client can reference it in API calls
		// but it is called TempID since you shouldn't store it for any reason
		TempID:  slugify.Slugify(auto.Name),
		Enabled: *auto.Enabled,
	}

	// This is called when the trigger triggers, we build the commands at this point
	triggered := func() {
		actions, err := parseActions(sys, auto)
		if err != nil {
			log.V("unable to build commands for automation: %s. %s", finalAuto.Name, err)
			return
		}

		if finalAuto.Triggered != nil {
			finalAuto.Triggered(actions)
		}
	}

	trigger, err := parseTrigger(sys, auto, triggered)
	if err != nil {
		return nil, err
	}
	finalAuto.Trigger = trigger

	return finalAuto, nil
}

func parseActions(sys automationSys, auto automationIntermediate) (*CommandGroup, error) {

	cmdGroup := CommandGroup{Desc: auto.Name}

	for _, action := range auto.Actions {
		if action.Scene != nil {
			scene := sys.SceneByID(action.Scene.ID)
			if scene == nil {
				return nil, fmt.Errorf("invalid scene ID: %s", action.Scene.ID)
			}

			cmdGroup.Cmds = append(cmdGroup.Cmds, &cmd.SceneSet{
				ID:        sys.NewID(),
				SceneID:   scene.ID,
				SceneName: scene.Name,
			})
		} else if action.LightZone != nil {
			lz := action.LightZone
			if lz.ID == nil && lz.AID == nil {
				// The user did not specify an ID, so we apply the attributes to all light zones
				lightZones := sys.FeaturesByType(feature.FTLightZone)
				if len(lightZones) == 0 {
					continue
				}

				for _, zn := range lightZones {
					command := buildLightZoneCommand(zn, lz.OnOff, lz.Brightness)
					if command == nil {
						continue
					}
					cmdGroup.Cmds = append(cmdGroup.Cmds, command)
				}
			} else {
				zn, err := getFeature(sys, lz.ID, lz.AID)
				if err != nil {
					return nil, err
				}

				command := buildLightZoneCommand(zn, lz.OnOff, lz.Brightness)
				// command might not apply to this particular zone
				if command == nil {
					continue
				}
				cmdGroup.Cmds = append(cmdGroup.Cmds, command)
			}
		} else if action.WindowTreatment != nil {
			if action.WindowTreatment.ID == nil && action.WindowTreatment.AID == nil {
				// The user did not specify an ID, so we apply the attributes to all window treatments
				treatments := sys.FeaturesByType(feature.FTWindowTreatment)
				if len(treatments) == 0 {
					continue
				}

				for _, wt := range treatments {
					command := buildWindowTreatmentCommand(wt, action.WindowTreatment.OpenClosed, action.WindowTreatment.Offset)
					if command == nil {
						continue
					}
					cmdGroup.Cmds = append(cmdGroup.Cmds, command)
				}
			} else {
				wt, err := getFeature(sys, action.WindowTreatment.ID, action.WindowTreatment.AID)
				if err != nil {
					return nil, err
				}

				command := buildWindowTreatmentCommand(wt, action.WindowTreatment.OpenClosed, action.WindowTreatment.Offset)
				if command == nil {
					continue
				}
				cmdGroup.Cmds = append(cmdGroup.Cmds, command)
			}
		} else if action.Outlet != nil {
			if action.Outlet.ID == nil && action.Outlet.AID == nil {
				outlets := sys.FeaturesByType(feature.FTOutlet)
				if len(outlets) == 0 {
					continue
				}

				for _, outlet := range outlets {
					command := buildOutletCommand(outlet, action.Outlet.OnOff)
					if command == nil {
						continue
					}
					cmdGroup.Cmds = append(cmdGroup.Cmds, command)
				}
			} else {
				outlet, err := getFeature(sys, action.Outlet.ID, action.Outlet.AID)
				if err != nil {
					return nil, err
				}
				command := buildOutletCommand(outlet, action.Outlet.OnOff)
				if command == nil {
					continue
				}
				cmdGroup.Cmds = append(cmdGroup.Cmds, command)
			}
		} else if action.Switch != nil {
			if action.Switch.ID == nil && action.Switch.AID == nil {
				switches := sys.FeaturesByType(feature.FTSwitch)
				if len(switches) == 0 {
					continue
				}

				for _, sw := range switches {
					command := buildSwitchCommand(sw, action.Switch.OnOff)
					if command == nil {
						continue
					}
					cmdGroup.Cmds = append(cmdGroup.Cmds, command)
				}
			} else {
				sw, err := getFeature(sys, action.Switch.ID, action.Switch.AID)
				if err != nil {
					return nil, err
				}
				command := buildSwitchCommand(sw, action.Switch.OnOff)
				if command == nil {
					continue
				}
				cmdGroup.Cmds = append(cmdGroup.Cmds, command)
			}
		} else if action.HeatZone != nil {
			if action.HeatZone.ID == nil && action.HeatZone.AID == nil {
				zones := sys.FeaturesByType(feature.FTHeatZone)
				if len(zones) == 0 {
					continue
				}

				for _, hz := range zones {
					command := buildHeatZoneCommand(hz, action.HeatZone.TargetTemp)
					if command == nil {
						continue
					}
					cmdGroup.Cmds = append(cmdGroup.Cmds, command)
				}
			} else {
				hz, err := getFeature(sys, action.HeatZone.ID, action.HeatZone.AID)
				if err != nil {
					return nil, err
				}
				command := buildHeatZoneCommand(hz, action.HeatZone.TargetTemp)
				if command == nil {
					continue
				}
				cmdGroup.Cmds = append(cmdGroup.Cmds, command)
			}
		} else {
			return nil, fmt.Errorf("unsupported action type")
		}
	}

	return &cmdGroup, nil
}

func getFeature(sys automationSys, id, aid *string) (*feature.Feature, error) {
	if aid != nil {
		f := sys.FeatureByAID(*aid)
		if f == nil {
			return nil, fmt.Errorf("invalid automation ID: %s", *aid)
		}
		return f, nil
	}

	if id != nil {
		f := sys.FeatureByID(*id)
		if f == nil {
			return nil, fmt.Errorf("invalid ID: %s", *id)
		}
		return f, nil
	}

	return nil, fmt.Errorf("invalid automation, missing id and aid key, one must be present")
}

func buildOutletCommand(outlet *feature.Feature, onOffVal *string) cmd.Command {
	if onOffVal == nil {
		log.V("missing on_off value for outlet ID: %s", outlet.ID)
		return nil
	}

	onoff := feature.OutletCloneAttrs(outlet)

	// NOTE: If we get an error we just log it an move on, since we want to try to execute as much
	// of the automation as possible even if one parts fails.

	var val int32
	switch *onOffVal {
	case "on":
		val = attr.OnOffOn
	case "off":
		val = attr.OnOffOff
	default:
		log.V("unsupported value for on_off, must be either [on|off], outlet ID: %s, %s", outlet.ID, *onOffVal)
		return nil
	}

	onoff.Value = val
	return &cmd.FeatureSetAttrs{
		FeatureID:   outlet.ID,
		FeatureName: outlet.Name,
		Attrs:       feature.NewAttrs(onoff),
	}
}

func buildHeatZoneCommand(hz *feature.Feature, targetTempVal *float64) cmd.Command {
	if targetTempVal == nil {
		log.V("missing target_temp field on heat zone: %s", hz.ID)
		return nil
	}

	_, targetTemp := feature.HeatZoneCloneAttrs(hz)
	targetTemp.Value = int32(*targetTempVal)

	return &cmd.FeatureSetAttrs{
		FeatureID:   hz.ID,
		FeatureName: hz.Name,
		Attrs:       feature.NewAttrs(targetTemp),
	}
}

func buildSwitchCommand(sw *feature.Feature, onOffVal *string) cmd.Command {
	if onOffVal == nil {
		log.V("missing on_off value for switch ID: %s", sw.ID)
		return nil
	}

	onoff := feature.SwitchCloneAttrs(sw)

	// NOTE: If we get an error we just log it an move on, since we want to try to execute as much
	// of the automation as possible even if one parts fails.

	var val int32
	switch *onOffVal {
	case "on":
		val = attr.OnOffOn
	case "off":
		val = attr.OnOffOff
	default:
		log.V("unsupported value for on_off, must be either [on|off], switch ID: %s, %s", sw.ID, *onOffVal)
		return nil
	}

	onoff.Value = val
	return &cmd.FeatureSetAttrs{
		FeatureID:   sw.ID,
		FeatureName: sw.Name,
		Attrs:       feature.NewAttrs(onoff),
	}
}

func buildLightZoneCommand(zn *feature.Feature, onOffVal *string, brightnessVal *float64) cmd.Command {
	onoff, brightness, _ := feature.LightZoneCloneAttrs(zn)

	// NOTE: If we get an error we just log it an move on, since we want to try to execute as much
	// of the automation as possible even if one parts fails.

	if onOffVal != nil {
		switch *onOffVal {
		case "on":
			onoff.Value = attr.OnOffOn
		case "off":
			onoff.Value = attr.OnOffOff
		default:
			log.V("unsupported value for on_off, must be either [on|off], light zone ID: %s, %s", zn.ID, *onOffVal)
		}
	} else {
		onoff = nil
	}

	if brightnessVal != nil {
		brightness.Value = float32(*brightnessVal)
	} else {
		brightness = nil
	}

	return &cmd.FeatureSetAttrs{
		FeatureID:   zn.ID,
		FeatureName: zn.Name,
		Attrs:       feature.NewAttrs(onoff, brightness),
	}
}

func buildWindowTreatmentCommand(wt *feature.Feature, openClosedVal *string, offsetVal *float64) cmd.Command {
	openclosed, offset := feature.WindowTreatmentCloneAttrs(wt)

	if openClosedVal != nil {
		switch *openClosedVal {
		case "open":
			openclosed.Value = attr.OpenCloseOpen
		case "closed":
			openclosed.Value = attr.OpenCloseClosed
		default:
			log.V("unsupported value for open_closed, must be either [open|closed], window treatment ID: %s, %s",
				wt.ID, *openClosedVal)
		}
	} else {
		openclosed = nil
	}

	if offsetVal != nil {
		offset.Value = float32(*offsetVal)
	} else {
		offset = nil
	}

	return &cmd.FeatureSetAttrs{
		FeatureID:   wt.ID,
		FeatureName: wt.Name,
		Attrs:       feature.NewAttrs(openclosed, offset),
	}
}

func parseCondition(f *feature.Feature, c *condition) error {
	if c.AttrLocalID != nil {
		c.feature = f

		_, ok := f.Attrs[*c.AttrLocalID]
		if !ok {
			return fmt.Errorf("invalid attr key: %s", *c.AttrLocalID)
		}

		if c.Op == nil {
			return fmt.Errorf("missing op key")
		}

		if c.Value == nil {
			return fmt.Errorf("missing value key")
		}

		return nil
	} else {
		return fmt.Errorf("condition is missing an 'attr' key")
	}

	return nil
}

func parseTrigger(sys automationSys, auto automationIntermediate, triggered func()) (Trigger, error) {
	if auto.Trigger.Feature != nil {
		if auto.Trigger.Feature.Condition == nil {
			return nil, fmt.Errorf("feature trigger missing condition key")
		}

		ft, err := getFeature(sys, auto.Trigger.Feature.ID, auto.Trigger.Feature.AID)
		if err != nil {
			return nil, err
		}

		err = parseCondition(ft, auto.Trigger.Feature.Condition)
		return &FeatureTrigger{
			Count:     auto.Trigger.Feature.Count,
			Duration:  time.Duration(auto.Trigger.Feature.Duration) * time.Millisecond,
			Triggered: triggered,
			Condition: auto.Trigger.Feature.Condition,
		}, nil

	} else if auto.Trigger.Time != nil {
		t := auto.Trigger.Time

		var mode string
		var at time.Time
		switch t.At {
		case "sunrise":
			mode = TimeTriggerModeSunrise
		case "sunset":
			mode = TimeTriggerModeSunset
		default:
			mode = TimeTriggerModeExact

			// This is a time, we support just a time or a datetime:
			// YYYY/MM/DD HH:MM:SS
			// HH:MM:SS
			var err error
			zoneName, _ := time.Now().Zone()
			at, err = time.Parse("2006/01/02 15:04:05 MST", t.At+" "+zoneName)

			if err != nil {
				// try just time
				at, err = time.Parse("15:04:05 MST", t.At+" "+zoneName)

				if err != nil {
					return nil, fmt.Errorf("invalid time input: %s, must be either HH:MM:SS or yyyy/MM/dd HH:mm:ss", t.At)
				}
			}
		}

		var days uint32
		if strings.Index(t.Days, "sun") != -1 {
			days |= TimeTriggerDaysSun
		}
		if strings.Index(t.Days, "mon") != -1 {
			days |= TimeTriggerDaysMon
		}
		if strings.Index(t.Days, "tues") != -1 {
			days |= TimeTriggerDaysTues
		}
		if strings.Index(t.Days, "wed") != -1 {
			days |= TimeTriggerDaysWed
		}
		if strings.Index(t.Days, "thurs") != -1 {
			days |= TimeTriggerDaysThurs
		}
		if strings.Index(t.Days, "fri") != -1 {
			days |= TimeTriggerDaysFri
		}
		if strings.Index(t.Days, "sat") != -1 {
			days |= TimeTriggerDaysSat
		}
		if t.Days == "" {
			days |= TimeTriggerDaysSun | TimeTriggerDaysMon | TimeTriggerDaysTues | TimeTriggerDaysWed |
				TimeTriggerDaysThurs | TimeTriggerDaysFri | TimeTriggerDaysSat
		}

		timeTrigger := &TimeTrigger{
			Name:      auto.Name,
			At:        at,
			Mode:      mode,
			Days:      days,
			Time:      clock.SystemTime{},
			Triggered: triggered,
		}
		return timeTrigger, nil
	} else {
		return nil, fmt.Errorf("unsupported trigger type")
	}
}

// When we unmarshal the scripts, the yaml parser will either return float64 or int for numbers
// we need float32 so we have to try to cast it correctly
func toFloat32(val interface{}) *float32 {
	if f64, ok := val.(float64); ok {
		f32 := float32(f64)
		return &f32
	}

	if i, ok := val.(int); ok {
		f32 := float32(i)
		return &f32
	}

	return nil
}

// Unmarshalling the yaml, we get either int or float64, need to cast to float32. Will return
// nil if the cast is unsuccessful
func toInt32(val interface{}) *int32 {
	if f64, ok := val.(float64); ok {
		i32 := int32(f64)
		return &i32
	}

	if i, ok := val.(int); ok {
		i32 := int32(i)
		return &i32
	}

	return nil
}
