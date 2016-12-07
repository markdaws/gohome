package gohome

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/go-yaml/yaml"
	"github.com/markdaws/gohome/clock"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

// Automation represents an automation instance. Each piece of automation has a trigger which is a set
// of conditions which when evaluating to true cause the automation actions to execute
type Automation struct {
	ID      string
	Name    string
	Enabled bool
	Trigger Trigger
	Actions *CommandGroup
	evtbus.Consumer
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

// helper type to deserialize the yaml in to our internal object model
type automationIntermediate struct {
	Name    string `yaml:"name"`
	Enabled *bool  `yaml:"enabled"`
	Trigger *struct {
		Time *struct {
			At   string `yaml:"at"`
			Days string `yaml:"days"`
		} `yaml:"time"`
	} `yaml:"trigger"`
	Actions []struct {
		Scene *struct {
			ID string `yaml:"id"`
		} `yaml:"scene"`
		LightZone *struct {
			ID    *string                `yaml:"id"`
			Attrs map[string]interface{} `yaml:"attrs"`
		} `yaml:"light_zone"`
		//TODO: Other feature types
	} `yaml:"actions"`
}

// LoadAutomation loads all of the automation files from the specified path
func LoadAutomation(sys *System, path string) ([]*Automation, error) {

	var autos []*Automation
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate automation files: %s", err)
	}

	for _, file := range files {
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

		autos = append(autos, auto)
	}
	return autos, nil
}

// NewAutomation creates a new automation instance
func NewAutomation(sys *System, config string) (*Automation, error) {

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

	actions, err := parseActions(sys, auto)
	if err != nil {
		return nil, err
	}

	trigger, err := parseTrigger(sys, auto, actions)
	if err != nil {
		return nil, err
	}

	return &Automation{
		ID:      sys.NewID(),
		Name:    auto.Name,
		Enabled: *auto.Enabled,
		Trigger: trigger,
		Actions: actions,
	}, nil
}

func parseActions(sys *System, auto automationIntermediate) (*CommandGroup, error) {

	cmdGroup := CommandGroup{Desc: auto.Name}

	for _, action := range auto.Actions {
		if action.Scene != nil {
			scene, ok := sys.Scenes[action.Scene.ID]
			if !ok {
				return nil, fmt.Errorf("invalid scene ID: %s", action.Scene.ID)
			}

			cmdGroup.Cmds = append(cmdGroup.Cmds, &cmd.SceneSet{
				ID:        sys.NewID(),
				SceneID:   scene.ID,
				SceneName: scene.Name,
			})

		} else {
			return nil, fmt.Errorf("unsupported action type")
		}
	}

	return &cmdGroup, nil
}

func parseTrigger(sys *System, auto automationIntermediate, actions *CommandGroup) (Trigger, error) {
	if auto.Trigger.Time != nil {
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
			//Offset - don't support right now
			At:   at,
			Mode: mode,
			Days: days,
			Time: clock.SystemTime{},
			Triggered: func() {
				sys.Services.CmdProcessor.Enqueue(*actions)
			},
		}
		return timeTrigger, nil
	} else {
		return nil, fmt.Errorf("unsupported trigger type")
	}
}
