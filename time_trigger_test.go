package gohome_test

import (
	"testing"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/stretchr/testify/require"
)

type MockSystem struct {
	CommandGroup *gohome.CommandGroup
}

func (s *MockSystem) Scene(ID string) *gohome.Scene {
	return &gohome.Scene{
		ID:   "45678",
		Name: "Mock Scene",
	}
}
func (s *MockSystem) CmdEnqueue(g gohome.CommandGroup) error {
	s.CommandGroup = &g
	return nil
}
func (s *MockSystem) NewID() string {
	return "12345"
}

type MockTime struct {
	now   time.Time
	after func(time.Duration) <-chan time.Time
}

func (mt MockTime) Now() time.Time {
	return mt.now
}

func (mt MockTime) After(d time.Duration) <-chan time.Time {
	return mt.after(d)
}

func TestSunrise(t *testing.T) {
	t.Parallel()

	ms := &MockSystem{}
	mt := MockTime{
		now:   time.Date(2016, time.December, 5, 10, 0, 0, 0, time.UTC),
		after: func(d time.Duration) <-chan time.Time { return time.After(d) },
	}

	trigger := &gohome.TimeTrigger{
		Time:        mt,
		Scener:      ms,
		NewIDer:     ms,
		CmdEnqueuer: ms,
		Name:        "test trigger",
		Mode:        gohome.TimeTriggerModeSunrise,
		Offset:      0,
		Days:        gohome.TimeTriggerDaysMon | gohome.TimeTriggerDaysFri,
		SceneID:     "12345",
	}

	ch := make(chan evtbus.Event)
	trigger.StartConsuming(ch)

	// Send sunrise event, check offset
	ch <- &gohome.SunriseEvt{}

	// Send sunset - nothing should happen
	time.Sleep(time.Second * 1)

	require.NotNil(t, ms.CommandGroup)

	// Make sure all of the required values are set
	sceneCmd, ok := ms.CommandGroup.Cmds[0].(*cmd.SceneSet)
	require.True(t, ok)

	require.Equal(t, "12345", sceneCmd.ID)
	require.Equal(t, "45678", sceneCmd.SceneID)
	require.Equal(t, "Mock Scene", sceneCmd.SceneName)
}

func TestSunset(t *testing.T) {
	t.Parallel()

	ms := &MockSystem{}
	mt := MockTime{
		now:   time.Date(2016, time.December, 5, 10, 0, 0, 0, time.UTC),
		after: func(d time.Duration) <-chan time.Time { return time.After(d) },
	}

	trigger := &gohome.TimeTrigger{
		Time:        mt,
		Scener:      ms,
		NewIDer:     ms,
		CmdEnqueuer: ms,
		Name:        "test trigger",
		Mode:        gohome.TimeTriggerModeSunset,
		Offset:      0,
		Days:        gohome.TimeTriggerDaysMon | gohome.TimeTriggerDaysFri,
		SceneID:     "12345",
	}

	ch := make(chan evtbus.Event)
	trigger.StartConsuming(ch)

	// Send sunrise event, check offset
	ch <- &gohome.SunsetEvt{}

	// Send sunset - nothing should happen
	time.Sleep(time.Second * 1)

	require.NotNil(t, ms.CommandGroup)

	// Make sure all of the required values are set
	sceneCmd, ok := ms.CommandGroup.Cmds[0].(*cmd.SceneSet)
	require.True(t, ok)

	require.Equal(t, "12345", sceneCmd.ID)
	require.Equal(t, "45678", sceneCmd.SceneID)
	require.Equal(t, "Mock Scene", sceneCmd.SceneName)
}

func TestExactWithDate(t *testing.T) {
	t.Parallel()

	ms := &MockSystem{}
	mt := MockTime{
		now:   time.Date(2016, time.December, 5, 10, 0, 0, 0, time.UTC),
		after: func(d time.Duration) <-chan time.Time { return time.After(d) },
	}

	trigger := &gohome.TimeTrigger{
		Time:        mt,
		Scener:      ms,
		NewIDer:     ms,
		CmdEnqueuer: ms,
		Name:        "test trigger",
		Mode:        gohome.TimeTriggerModeExact,
		At:          mt.Now().Add(time.Second * 1),
		SceneID:     "12345",
	}

	ch := make(chan evtbus.Event)
	trigger.StartConsuming(ch)

	time.Sleep(time.Second * 2)

	require.NotNil(t, ms.CommandGroup)

	// Make sure all of the required values are set
	sceneCmd, ok := ms.CommandGroup.Cmds[0].(*cmd.SceneSet)
	require.True(t, ok)

	require.Equal(t, "12345", sceneCmd.ID)
	require.Equal(t, "45678", sceneCmd.SceneID)
	require.Equal(t, "Mock Scene", sceneCmd.SceneName)

}

func TestExactWithoutDate(t *testing.T) {
	t.Parallel()

	mt := MockTime{
		// This is a monday
		now: time.Date(2016, time.December, 5, 10, 10, 0, 0, time.UTC),
		after: func(d time.Duration) <-chan time.Time {
			// Return immediately
			c := make(chan time.Time, 1)
			c <- time.Now()
			return c
		},
	}
	ms := &MockSystem{}

	// If the time doesn't have a date, then it runs every day (for the days it was
	// specified) at the desired time
	nullTime := time.Time{}

	// This is the same as the current mock time
	at := time.Date(
		nullTime.Year(), nullTime.Month(), nullTime.Day(),
		mt.Now().Hour(), mt.Now().Minute(), mt.Now().Second(), 0, mt.Now().Location())

	// move in to the future
	at = at.Add(time.Second)

	evalCount := 0
	eval := func() {
		// This is called each time the trigger evaluates if it should run

		// If this is not the first time it is evaluating, then pretend we incremented
		// the current time by 23 hours to the next day (not 24 incase we jump past the
		// execution time), so that the test can run over multiple days. We only want the
		// test to execute on Monday and Friday
		if evalCount > 0 {
			mt.now = mt.now.Add(time.Hour * 23)
		}

		// Start on Monday, if eval is one, it means we ran on monday, so we should have got
		// an execute scene, then also on friday
		if evalCount == 1 || evalCount == 5 {
			require.NotNil(t, ms.CommandGroup)

			sceneCmd, ok := ms.CommandGroup.Cmds[0].(*cmd.SceneSet)
			require.True(t, ok)

			require.Equal(t, "12345", sceneCmd.ID)
			require.Equal(t, "45678", sceneCmd.SceneID)
			require.Equal(t, "Mock Scene", sceneCmd.SceneName)
		} else {
			require.Nil(t, ms.CommandGroup)
		}
		ms.CommandGroup = nil
		evalCount++
	}

	trigger := &gohome.TimeTrigger{
		Time:        &mt,
		Scener:      ms,
		NewIDer:     ms,
		CmdEnqueuer: ms,
		Name:        "test trigger - withoutdate",
		Mode:        gohome.TimeTriggerModeExact,
		Days:        gohome.TimeTriggerDaysMon | gohome.TimeTriggerDaysFri,
		At:          at,
		SceneID:     "12345",
		Evaluating:  eval,
	}

	ch := make(chan evtbus.Event)
	trigger.StartConsuming(ch)

	time.Sleep(time.Second * 2)
}
