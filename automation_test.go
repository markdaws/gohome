package gohome_test

import (
	"testing"
	"time"

	"github.com/markdaws/gohome"
	"github.com/stretchr/testify/require"
)

// Time with year + time
// Time with only time
// No days specified, should default to every day
// Actions...
func TestTimeTriggerNoDate(t *testing.T) {
	t.Parallel()

	config := `
name: Test
trigger:
  time:
    at: '13:59:30'
    days: mon|fri
actions:
  - scene:
      id: 12345
`
	/*
	   	  - light_zone:
	         id: 12345
	         attrs:
	           onoff: on
	           brightness: 98
	*/

	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)

	trigger := auto.Trigger.(*gohome.TimeTrigger)
	require.Equal(t, gohome.TimeTriggerModeExact, trigger.Mode)
	require.Equal(t, gohome.TimeTriggerDaysMon|gohome.TimeTriggerDaysFri, trigger.Days)
	require.Equal(t, time.Date(0, 1, 1, 13, 59, 30, 0, time.Now().Location()), trigger.At)
}

func TestTimeTriggerWithDate(t *testing.T) {
	t.Parallel()

	config := `
name: Test
trigger:
  time:
    at: '2016/11/19 13:59:30'
    days: mon|fri
actions:
  - scene:
      id: 12345
`
	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)

	trigger := auto.Trigger.(*gohome.TimeTrigger)
	require.Equal(t, gohome.TimeTriggerModeExact, trigger.Mode)
	require.Equal(t, gohome.TimeTriggerDaysMon|gohome.TimeTriggerDaysFri, trigger.Days)
	require.Equal(t, time.Date(2016, 11, 19, 13, 59, 30, 0, time.Now().Location()), trigger.At)
}

func TestTimeTriggerSunrise(t *testing.T) {
	t.Parallel()

	config := `
name: Test
trigger:
  time:
    at: sunrise
    days: mon|fri
actions:
  - scene:
      id: 12345
`
	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)

	trigger := auto.Trigger.(*gohome.TimeTrigger)
	require.Equal(t, gohome.TimeTriggerModeSunrise, trigger.Mode)
	require.Equal(t, gohome.TimeTriggerDaysMon|gohome.TimeTriggerDaysFri, trigger.Days)
}

func TestTimeTriggerSunset(t *testing.T) {
	t.Parallel()

	config := `
name: Test
trigger:
  time:
    at: sunset
    days: mon|fri
actions:
  - scene:
      id: 12345
`
	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)

	trigger := auto.Trigger.(*gohome.TimeTrigger)
	require.Equal(t, gohome.TimeTriggerModeSunset, trigger.Mode)
	require.Equal(t, gohome.TimeTriggerDaysMon|gohome.TimeTriggerDaysFri, trigger.Days)
}

func TestMissingNameField(t *testing.T) {
	t.Parallel()

	config := `
trigger:
  time:
    at: '03:59:30'
    days: mon|tues|wed|thurs|fri|sat|sun
actions:
  - scene:
      id: 12345
`

	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	_, err := gohome.NewAutomation(sys, config)
	require.NotNil(t, err)
}

func TestMissingActionsKey(t *testing.T) {
	t.Parallel()

	config := `
name: test
trigger:
  time:
    at: '03:59:30'
    days: mon|tues|wed|thurs|fri|sat|sun
`
	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	_, err := gohome.NewAutomation(sys, config)
	require.NotNil(t, err)
}

func TestMissingTriggerKey(t *testing.T) {
	t.Parallel()

	config := `
name: test
actions:
  - scene:
      id: 12345
`

	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	_, err := gohome.NewAutomation(sys, config)
	require.NotNil(t, err)
}

func TestWithoutEnabledDefaultsToTrue(t *testing.T) {
	t.Parallel()

	config := `
name: test
trigger:
  time:
    at: '03:59:30'
    days: mon|tues|wed|thurs|fri|sat|sun
actions:
  - scene:
      id: 12345
`

	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)
	require.True(t, auto.Enabled)
}

func TestEnabledFalse(t *testing.T) {
	t.Parallel()

	config := `
name: test
enabled: false
trigger:
  time:
    at: '03:59:30'
    days: mon|tues|wed|thurs|fri|sat|sun
actions:
  - scene:
      id: 12345
`

	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)
	require.False(t, auto.Enabled)
}

func TestEnabledTrue(t *testing.T) {
	t.Parallel()

	config := `
name: test
enabled: true
trigger:
  time:
    at: '03:59:30'
    days: mon|tues|wed|thurs|fri|sat|sun
actions:
  - scene:
      id: 12345
`

	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)
	require.True(t, auto.Enabled)
}
