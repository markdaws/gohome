package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/log"
)

type GoHomeHubDevice struct {
	device
}

func (d *GoHomeHubDevice) ModelNumber() string {
	return "GoHomeHub"
}

func (d *GoHomeHubDevice) InitConnections() {
	//ci := *d.connectionInfo.(*comm.TelnetConnectionInfo)

	//TODO: Get connection info externally, one pool per controlled device?
	ci := comm.TelnetConnectionInfo{
		PoolSize: 2,
		Network:  "tcp",
		Address:  "192.168.0.24:5577",
	}
	createConnection := func() comm.Connection {
		conn := comm.NewTelnetConnection(ci)
		/*
			conn.SetPingCallback(func() error {
				if _, err := conn.Write([]byte("#PING\r\n")); err != nil {
					return fmt.Errorf("%s ping failed: %s", d, err)
				}
				return nil
			})*/
		return conn
	}
	ps := ci.PoolSize
	log.V("%s init connections, pool size %d", d, ps)
	d.pool = comm.NewConnectionPool(d.name, ps, createConnection)
	log.V("%s connected", d)

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

				var rV, gV, bV byte
				lvl := c.Level.Value
				if lvl == 0 {
					if (c.Level.R == 0) && (c.Level.G == 0) && (c.Level.B == 0) {
						rV = 0
						gV = 0
						bV = 0
					} else {
						rV = c.Level.R
						gV = c.Level.G
						bV = c.Level.B
					}
				} else {
					rV = byte((lvl / 100) * 255)
					gV = rV
					bV = rV
				}

				b := []byte{0x31, rV, gV, bV, 0x00, 0xf0, 0x0f}
				var t int = 0
				for _, v := range b {
					t += int(v)
				}
				cs := t & 0xff
				b = append(b, byte(cs))

				conn, err := d.Connect()
				if err != nil {
					return fmt.Errorf("StringCommand - error connecting %s", err)
				}

				defer func() {
					d.ReleaseConnection(conn)
				}()
				_, err = conn.Write(b)
				if err != nil {
					fmt.Printf("ERROR SENDING %s\n", err)
				} else {
				}
				return err
			},
		}, nil
	default:
		fmt.Println(z.Controller)
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
