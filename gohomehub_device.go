package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
)

type GoHomeHubDevice struct {
	device
}

func (d *GoHomeHubDevice) ModelNumber() string {
	return "GoHomeHub"
}

func (d *GoHomeHubDevice) InitConnections() {
}

func (d *GoHomeHubDevice) StartProducingEvents() (<-chan Event, <-chan bool) {
	return nil, nil
}

func (d *GoHomeHubDevice) Authenticate(c comm.Connection) error {
	return nil
}

func (d *GoHomeHubDevice) BuildCommand(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneSetLevel:
		return d.buildZoneSetLevelCommand(command)
	case *cmd.ButtonPress:
		//TODO: Phantom buttons?
		return nil, fmt.Errorf("goHomeHubDevice ButtonPressCommand not supported")
	case *cmd.ButtonRelease:
		return nil, fmt.Errorf("goHomeHubDevice ButtonReleaseCommand not supported")
	case *cmd.SceneSet:
		//TODO: Does this make sense, what does a scene mean in terms of this virtual hub?
	default:
		return nil, fmt.Errorf("goHomeHubDevice build commands not supported")
	}

	return nil, fmt.Errorf("goHomeHubDevice unsupported command")
}

//TODO: Level should be a type with value,r,g,b, not just one value
func (d *GoHomeHubDevice) buildZoneSetLevelCommand(c *cmd.ZoneSetLevel) (*cmd.Func, error) {

	z, ok := d.System().Zones[c.ZoneGlobalID]
	if !ok {
		return nil, fmt.Errorf("unknown zone ID %s", c.ZoneGlobalID)
	}

	switch z.Controller {
	case ZCFluxWIFI:
		return &cmd.Func{
			Func: func() error {
				return fmt.Errorf("not implmented ghh")
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported zone controller")
	}
}

/*
def __init__(self, ipaddr, port=5577):
self.ipaddr = ipaddr
self.port = port
self.__is_on = False

self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
self.socket.connect((self.ipaddr, self.port))

def setRgb(self, r,g,b, persist=True):
if persist:
msg = bytearray([0x31])
else:
msg = bytearray([0x41])
msg.append(r)
msg.append(g)
msg.append(b)
msg.append(0x00)
msg.append(0xf0)
msg.append(0x0f)
self.__write(msg)

def __writeRaw(self, bytes):
self.socket.send(bytes)

def __write(self, bytes):
# calculate checksum of byte array and add to end
csum = sum(bytes) & 0xFF
bytes.append(csum)
#print "-------------",utils.dump_bytes(bytes)
self.__writeRaw(bytes)
#time.sleep(.4)
*/

/*
//TODO: Add credit to github repo
- For wifi bulbs can control directly
- Capabilities supports zigbee etc
- Need to add bulb information
- Is a zone a bulb?

zone is a controllable unit
*/
