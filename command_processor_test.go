package gohome_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/zone"
	"github.com/stretchr/testify/require"
)

type mockBuilder struct {
	WaitGroup *sync.WaitGroup
}

func (b *mockBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	return &cmd.Func{
		Func: func() error {
			fmt.Println("exec mock func")
			b.WaitGroup.Done()
			panic("panic worker")
			return nil
		},
		Friendly: "mock func",
	}, nil
}
func (b *mockBuilder) ID() string {
	return "mockBuilder"
}

func makeTestSystem(cp gohome.CommandProcessor, b *mockBuilder) *gohome.System {
	s := gohome.NewSystem("mock system", "", cp, 1)

	s.Extensions.CmdBuilders["mock"] = b

	d, _ := gohome.NewDevice("abcd", "", "1", "mock dev", "", nil, false, b, nil, nil)
	s.AddDevice(*d)

	z := &zone.Zone{Name: "z1", Address: "1", ID: "z1", DeviceID: d.ID}
	s.AddZone(z)
	return s
}

func TestWorkersShouldRestartAfterPanic(t *testing.T) {
	cp := gohome.NewCommandProcessor(2, 100)

	var wg sync.WaitGroup
	wg.Add(10)

	cmdBuilder := &mockBuilder{&wg}
	s := makeTestSystem(cp, cmdBuilder)
	cp.SetSystem(s)
	cp.Start()

	// Mock commands will panic when executed, but the command processor should
	// keep processing them as the workers restart
	for i := 0; i < 10; i++ {
		cp.Enqueue(gohome.NewCommandGroup("mock group", &cmd.ZoneTurnOn{ZoneID: "z1"}))
	}

	// Wait here until all 10 requests are processed, if something goes wrong we will be stuck here
	// and the test will fail
	wg.Wait()
	fmt.Println("all workers done")
}

func TestEnqueueToFullChannelReturnsError(t *testing.T) {
	// If we try to enqueue a job and the command processor channel is
	// full, then we should get an error

	// Have 0 workers to simulate queue backing up
	cp := gohome.NewCommandProcessor(0, 1)
	cp.Start();

	// Queue holds up to 1 command
	err := cp.Enqueue(gohome.NewCommandGroup("mock group", &cmd.ZoneTurnOn{ZoneID: "z1"}))
	require.Nil(t, err)

	// Should get an error this time and the enqueue should not block on trying to
	// add to the channel
	err = cp.Enqueue(gohome.NewCommandGroup("mock group", &cmd.ZoneTurnOn{ZoneID: "z1"}))
	require.NotNil(t, err)
}
