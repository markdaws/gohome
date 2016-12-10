package gohome_test

import (
	"testing"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/feature"
	"github.com/stretchr/testify/require"
)

/*
type MockSystem struct {
}

func (s *MockSystem) NewID() string {
	return ""
}
func (s *MockSystem) SceneByID(ID string) *Scene {
	return nil
}
func (s *MockSystem) FeaturesByType(ft string) map[string]*feature.Feature {
}
func (s *MockSystem) FeatureByID(ID string) *feature.Feature {
}
func (s *MockSystem) FeatureByAID(AID string) *feature.Feature {
}*/

func TestFeatureTriggerCount(t *testing.T) {
	t.Parallel()

	config := `
name: Test
trigger:
  feature:
    count: 3
    duration: 3000
    id: 12345
    condition:
      attr: 'openclose'
      op: '=='
      value: 2
actions:
  - light_zone:
      on_off: 'on'
`

	sys := gohome.NewSystem("test system")
	sensor := feature.NewSensor("12345", attr.NewOpenClose("openclose", nil))
	sys.AddFeature(sensor)
	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)

	// Simulate firing the sensor 3 times in less that 3 seconds, should trigger the event
	ch := make(chan evtbus.Event)
	auto.StartConsuming(ch)

	wasTriggered := false
	auto.Triggered = func(actions *gohome.CommandGroup) {
		wasTriggered = true
	}

	//Update
	openclosed := sensor.Attrs["openclose"].Clone()
	openclosed.Value = attr.OpenCloseOpen

	evt := &gohome.FeatureAttrsChangedEvt{
		FeatureID: "12345",
		Attrs:     feature.NewAttrs(openclosed),
	}

	require.False(t, wasTriggered)
	ch <- evt
	time.Sleep(100 * time.Millisecond)
	require.False(t, wasTriggered)

	ch <- evt
	time.Sleep(100 * time.Millisecond)
	require.False(t, wasTriggered)

	ch <- evt
	time.Sleep(100 * time.Millisecond)

	// Should trigger after 3 events inside the required time
	require.True(t, wasTriggered)
}

func TestFeatureTriggerCountExpiredDuration(t *testing.T) {
	// Make sure that if we trigger events but not within the desired time the trigger doesn't fire
	t.Parallel()

	config := `
name: Test
trigger:
  feature:
    count: 3
    duration: 1000
    id: 12345
    condition:
      attr: 'openclose'
      op: '=='
      value: 2
actions:
  - light_zone:
      on_off: 'on'
`

	sys := gohome.NewSystem("test system")
	sensor := feature.NewSensor("12345", attr.NewOpenClose("openclose", nil))
	sys.AddFeature(sensor)
	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)

	// Simulate firing the sensor 3 times in less that 3 seconds, should trigger the event
	ch := make(chan evtbus.Event)
	auto.StartConsuming(ch)

	wasTriggered := false
	auto.Triggered = func(actions *gohome.CommandGroup) {
		wasTriggered = true
	}

	//Update
	openclosed := sensor.Attrs["openclose"].Clone()
	openclosed.Value = attr.OpenCloseOpen

	evt := &gohome.FeatureAttrsChangedEvt{
		FeatureID: "12345",
		Attrs:     feature.NewAttrs(openclosed),
	}

	require.False(t, wasTriggered)
	ch <- evt
	time.Sleep(100 * time.Millisecond)
	require.False(t, wasTriggered)

	ch <- evt
	time.Sleep(100 * time.Millisecond)
	require.False(t, wasTriggered)

	//Make sure the last event is past the duration so the trigger shouldn't fire
	time.Sleep(time.Second * 2)
	ch <- evt
	time.Sleep(100 * time.Millisecond)

	// Should not have triggered
	require.False(t, wasTriggered)

	// Make sure it will trigger after a failure
	ch <- evt
	time.Sleep(100 * time.Millisecond)
	require.False(t, wasTriggered)

	ch <- evt
	time.Sleep(100 * time.Millisecond)
	require.True(t, wasTriggered)
}

func TestSensorTrigger(t *testing.T) {
	t.Parallel()

	config := `
name: Test
trigger:
  sensor:
    condition:
      attr: 'on_off'
      op: 'eq'
      value: 'off'
actions:
  - scene:
      id: 12345
`

	sys := gohome.NewSystem("test system")
	s1 := &gohome.Scene{ID: "12345"}
	sys.AddScene(s1)

	auto, err := gohome.NewAutomation(sys, config)
	require.Nil(t, err)
	_ = auto
}

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

func TestTimeWithNoDaysDefaultsToEveryDay(t *testing.T) {
	t.Parallel()

	config := `
name: Test
trigger:
  time:
    at: '2016/11/19 13:59:30'
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
	require.Equal(t,
		gohome.TimeTriggerDaysSun|gohome.TimeTriggerDaysMon|gohome.TimeTriggerDaysTues|
			gohome.TimeTriggerDaysWed|gohome.TimeTriggerDaysThurs|gohome.TimeTriggerDaysFri|
			gohome.TimeTriggerDaysSat, trigger.Days)
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
