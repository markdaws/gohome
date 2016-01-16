package gohome

import "testing"

type MockTrigger struct {
	ProcessFunc     func(Event) bool
	ProcessesEvents bool
	Fires           chan bool
}

func (m *MockTrigger) Type() string {
	return "mockTrigger"
}

func (m *MockTrigger) Ingredients() []Ingredient {
	return nil
}

func (m *MockTrigger) Name() string {
	return "Mock Trigger"
}

func (m *MockTrigger) Description() string {
	return "A Mock Trigger"
}

func (m *MockTrigger) New() Trigger {
	return &MockTrigger{}
}

func (m *MockTrigger) Init(done <-chan bool) (<-chan bool, bool) {
	return m.Fires, m.ProcessesEvents
}

func (m *MockTrigger) ProcessEvent(e Event) bool {
	return m.ProcessFunc(e)
}

type MockAction struct {
	ExecuteFunc func(*System) error
}

func (m *MockAction) Type() string {
	return "mockAction"
}

func (m *MockAction) Ingredients() []Ingredient {
	return nil
}

func (m *MockAction) Name() string {
	return "Mock Action"
}

func (m *MockAction) Description() string {
	return "A Mock Action"
}

func (m *MockAction) New() Action {
	return &MockAction{}
}

func (m *MockAction) Execute(s *System) error {
	return m.ExecuteFunc(s)
}

func TestNewRecipe(t *testing.T) {
	tr := &MockTrigger{}
	a := &MockAction{}
	s := NewSystem("n", "d")
	r, err := NewRecipe("n", "d", true, tr, a, s)

	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if r.ID == "" {
		t.Errorf("ID should not be empty")
	}

	if r.Name != "n" {
		t.Errorf("Name not set")
	}
	if r.Description != "d" {
		t.Errorf("Description not set")
	}
	if r.Trigger != tr {
		t.Errorf("Trigger not set")
	}
	if r.Action != a {
		t.Errorf("Action not set")
	}
	if r.Version != "1" {
		t.Errorf("Invalid version")
	}
	if r.system != s {
		t.Errorf("system not set")
	}
	if r.enabled != true {
		t.Errorf("enabled not set")
	}
}

func TestEnabled(t *testing.T) {
	r, _ := NewRecipe("n", "d", true, nil, nil, nil)
	if r.Enabled() != true {
		t.Errorf("enabled not set")
	}

	r.SetEnabled(false)
	if r.Enabled() != false {
		t.Errorf("SetEnabled did not set enabled state")
	}
}

func TestEventConsumerID(t *testing.T) {
	r, _ := NewRecipe("n", "d", true, nil, nil, nil)
	if r.EventConsumerID() != r.Name+" - "+r.ID {
		t.Errorf("EventConsumerID() returned unexpected ID")
	}
}

func TestTriggerFiringExecutesAction(t *testing.T) {

	executed := make(chan bool, 2)

	e := Event{}
	tr := &MockTrigger{
		ProcessesEvents: true,
		ProcessFunc: func(e Event) bool {
			return true
		},
	}
	a := &MockAction{
		ExecuteFunc: func(s *System) error {
			executed <- true
			return nil
		},
	}
	s := NewSystem("n", "d")
	r, _ := NewRecipe("n", "d", true, tr, a, s)

	c := r.StartConsumingEvents()
	c <- e
	c <- e

	// Wait until the action has executed
	<-executed
	<-executed
}

func TestClosingConsumerChannelStopsRecipe(t *testing.T) {
	executed := make(chan bool, 2)
	fires := make(chan bool)
	e := Event{}
	tr := &MockTrigger{
		ProcessesEvents: true,
		Fires:           fires,
		ProcessFunc: func(e Event) bool {
			return true
		},
	}
	a := &MockAction{
		ExecuteFunc: func(s *System) error {
			executed <- true
			return nil
		},
	}
	s := NewSystem("n", "d")
	r, _ := NewRecipe("n", "d", true, tr, a, s)

	c := r.StartConsumingEvents()
	c <- e
	<-executed

	close(c)
}

func TestTriggerFireNotOnProcessEventExecutesAction(t *testing.T) {
	executed := make(chan bool, 2)

	fires := make(chan bool)
	tr := &MockTrigger{
		ProcessesEvents: false,
		Fires:           fires,
	}
	a := &MockAction{
		ExecuteFunc: func(s *System) error {
			executed <- true
			return nil
		},
	}
	s := NewSystem("n", "d")
	r, _ := NewRecipe("n", "d", true, tr, a, s)

	_ = r.StartConsumingEvents()
	fires <- true
	<-executed
	fires <- true
	<-executed
}

func TestStartStopConsumingMultipleTimes(t *testing.T) {
	executed := make(chan bool, 2)

	e := Event{}
	tr := &MockTrigger{
		ProcessesEvents: true,
		ProcessFunc: func(e Event) bool {
			return true
		},
	}
	a := &MockAction{
		ExecuteFunc: func(s *System) error {
			executed <- true
			return nil
		},
	}
	s := NewSystem("n", "d")
	r, _ := NewRecipe("n", "d", true, tr, a, s)

	c := r.StartConsumingEvents()
	c <- e
	<-executed

	close(c)
	c = r.StartConsumingEvents()
	c <- e
	<-executed
}
