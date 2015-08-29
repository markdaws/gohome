package gohome

import (
	"fmt"
	"io"
	"testing"

	"github.com/markdaws/gohome"
)

type FakeReader struct {
	count int
}

//TODO: Pass in data stream
func (r *FakeReader) Read(p []byte) (n int, err error) {
	var data = []string{
		"GNET>\r\n",
		"GNET>",
		"something else \r\n",
		"GNET>\r\n",
		"GNET> ~OUTPUT,2,1,25.00\r\n",
		"~OUTPUT,2,1,0.00\r\n",
		"~DEVICE,1,9,3\r\n",
		"~OUTPUT,16,1,0.00\r\n",
		"~DEVICE,1,9,4\r\n",
		"\r\n",
		"just some garbage    ",
		"~ERROR 1\r\n",
		"GNET>\r\n",
		"GNET> ~DEVICE,1,7,3\r\n",
		"~OUTPUT,16,1,0.00\r\n",
		"~DEVICE,1,7,4\r\n",
	}

	if r.count >= len(data) {
		return 0, io.EOF
	}

	bts := []byte(data[r.count])
	for i, b := range bts {
		p[i] = b
	}
	n = len(bts)
	err = nil
	r.count++
	return
}

func TestStream(t *testing.T) {
	s := &gohome.System{
		Identifiable: gohome.Identifiable{
			ID:          "Sys1",
			Name:        "Test System",
			Description: "Test System description",
		},
		Devices: map[string]*gohome.Device{},
		Zones: map[string]*gohome.Zone{
			"2": &gohome.Zone{
				Identifiable: gohome.Identifiable{
					ID:          "2",
					Name:        "Kitchen Lights",
					Description: "Main lights in the kitchen",
				},
			},
			"16": &gohome.Zone{
				Identifiable: gohome.Identifiable{
					ID:          "16",
					Name:        "Bathroom Lights",
					Description: "Main lights in the bathroom",
				},
			},
		},
	}
	d := &gohome.Device{
		Identifiable: gohome.Identifiable{
			ID:          "1",
			Name:        "Main device",
			Description: "Testing main device"},
		System:     s,
		Connection: nil,
	}
	s.Devices[d.ID] = d

	fc, dc := d.StartProducingEvents()
	go func() {
		for {
			select {
			case e := <-fc:
				fmt.Println(e.String())
			}
		}
	}()

	go func() {
		gohome.Stream(d, &FakeReader{})
	}()

	<-dc
}
