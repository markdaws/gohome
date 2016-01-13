package comm

import (
	"fmt"
	"testing"
)

type mockConnection struct {
	OpenCalled    bool
	OpenError     error
	OpenCallback  func()
	CloseCalled   bool
	ReadN         int
	ReadError     error
	WriteN        int
	WriteError    error
	Status_       ConnectionStatus
	PingCallback_ PingCallback
	PingCount     int
}

func (c *mockConnection) Open() error {
	c.OpenCalled = true
	c.Status_ = CSConnected
	if c.OpenCallback != nil {
		c.OpenCallback()
	}
	return c.OpenError
}

func (c *mockConnection) Close() {
	c.CloseCalled = true
}

func (c *mockConnection) Read(p []byte) (n int, err error) {
	return c.ReadN, c.ReadError
}

func (c *mockConnection) Write(p []byte) (n int, err error) {
	return c.WriteN, c.WriteError
}

func (c *mockConnection) SetPingCallback(h PingCallback) {
	c.PingCallback_ = h
}

func (c *mockConnection) PingCallback() PingCallback {
	return c.PingCallback_
}

func (c *mockConnection) Status() ConnectionStatus {
	return c.Status_
}

func (c *mockConnection) SetStatus(s ConnectionStatus) {
	c.Status_ = s
}

func TestNewConnectionPool(t *testing.T) {
	p := createPool(3, nil)
	c1 := p.Get().(*mockConnection)
	c2 := p.Get().(*mockConnection)
	c3 := p.Get().(*mockConnection)
	c4, _ := p.Get().(*mockConnection)

	if c1 == nil || c2 == nil || c3 == nil || c4 != nil {
		t.Errorf("Unexpected value, expected (!nil, !nil, !nil, nil) got: %v, %v, %v, %v", c1, c2, c3, c4)
	}

	if !c1.OpenCalled || !c2.OpenCalled || !c2.OpenCalled {
		t.Errorf("Open not called on all new connections\n")
	}

	if c1.Status_ != CSConnected || c2.Status_ != CSConnected || c3.Status_ != CSConnected {
		t.Errorf("Not all connections are marked as connected\n")
	}
}

func TestReleaseClosedConnectionNotReturnedToPool(t *testing.T) {

	c1 := &mockConnection{}
	c2 := &mockConnection{}
	callCount := 0
	p := createPool(1, func() Connection {
		callCount++
		if callCount == 1 {
			return c1
		} else {
			return c2
		}
	})

	c := p.Get().(*mockConnection)
	if c != c1 {
		fmt.Errorf("unexpected connection returned")
	}

	if p.Get() != nil {
		fmt.Errorf("Pool with no connections should return nil")
	}

	c.Status_ = CSClosed
	p.Release(c)

	// Readding to the pool is async
	wait := make(chan bool)
	p.SetNewConnectionCallback(func() {
		c = p.Get().(*mockConnection)
		if c != c2 {
			fmt.Errorf("Unexpected connection returned from pool")
		}
		close(wait)
	})
	<-wait
}

func TestConnectionReturnedToPool(t *testing.T) {
	p := createPool(1, nil)
	c1 := p.Get()
	c2 := p.Get()
	if c2 != nil {
		fmt.Errorf("Unexpected connection")
	}

	p.Release(c1)
	c3 := p.Get()
	if c3 != c1 {
		fmt.Errorf("Unexpected connection on second get")
	}
}

func createPool(count int, newFunc func() Connection) ConnectionPool {
	if newFunc == nil {
		newFunc = func() Connection {
			c := &mockConnection{}
			c.SetPingCallback(func() error {
				c.PingCount++
				return nil
			})
			return c
		}
	}
	p := NewConnectionPool(count, newFunc)
	return p
}
