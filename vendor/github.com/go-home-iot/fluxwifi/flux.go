package fluxwifi

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// CREDIT: All the knowledge of how to control this product came from:
// https://github.com/beville/flux_led

// State represents the current state of the bulb
type State struct {
	// Power - 0 -> off, 1 -> on, 2-> unknown
	Power int

	// Mode - ww|color|custom|unknown
	Mode string

	// R - red value if in color mode
	R byte

	// G - green value if in color mode
	G byte

	// B - blue value if in color mode
	B byte
}

// SetLevel changes the RGB values of the bulb.
func SetLevel(r, g, b byte, w io.Writer) error {
	return write(w, []byte{0x31, r, g, b, 0x00, 0xf0, 0x0f})
}

// TurnOn turns the bulb on
func TurnOn(w io.Writer) error {
	return write(w, []byte{0x71, 0x23, 0x0f})
}

// TurnOff turns the bulb off
func TurnOff(w io.Writer) error {
	return write(w, []byte{0x71, 0x24, 0x0f})
}

// FetchState returns the current state of the bulb
func FetchState(conn net.Conn) (*State, error) {
	err := write(conn, []byte{0x81, 0x8a, 0x8b})
	if err != nil {
		return nil, err
	}

	resp := make([]byte, 100)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	n, err := conn.Read(resp)
	if err != nil {
		return nil, err
	}

	// Seems like flux sends back multiple responses if they are backed up, sometimes not all
	// complete, so just grab the last 14 bytes
	if n > 14 {
		resp = resp[n-14 : n]
		if len(resp) != 14 {
			return nil, fmt.Errorf("unknown response from get state")
		}
	}

	state := &State{}
	power := resp[2]
	switch power {
	case 0x23:
		state.Power = 1
	case 0x24:
		state.Power = 0
	default:
		state.Power = 2
	}

	pattern := resp[3]
	wwLevel := resp[9]
	mode := "unknown"
	switch pattern {
	case 0x61, 0x62:
		if wwLevel != 0 {
			mode = "ww"
		} else {
			mode = "color"
		}
	case 0x60:
		mode = "custom"
	}
	state.Mode = mode

	switch mode {
	case "color":
		state.R = resp[6]
		state.G = resp[7]
		state.B = resp[8]
	}
	return state, nil
}

func write(w io.Writer, b []byte) error {
	var t int
	for _, v := range b {
		t += int(v)
	}
	cs := t & 0xff
	b = append(b, byte(cs))
	_, err := w.Write(b)
	return err
}

// Scan scans the local network looking for Flux WIFI bulbs. waitTimeSeconds
// determines how long the function waits until it stops listening for
// responses on the network
func Scan(waitTimeSeconds int) ([]BulbInfo, error) {
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 48899,
	})

	if socket != nil {
		defer socket.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("failed trying to create a socket: %s", err)
	}

	done := make(chan bool)
	var infos []BulbInfo

	magicScanConst := "HF-A11ASSISTHREAD"
	go func() {
		readDeadline := time.Now().Add(time.Second * time.Duration(waitTimeSeconds))
		for {
			socket.SetReadDeadline(readDeadline)
			buff := make([]byte, 4096)
			n, err := socket.Read(buff)
			if err != nil {
				close(done)
				break
			}

			resp := string(buff[:n])

			// Ignore request on the network meant for the bulbs
			if resp == magicScanConst {
				continue
			}

			parts := strings.Split(resp, ",")
			if len(parts) != 3 {
				//should be 3 parts, ignore
				continue
			}

			infos = append(infos, BulbInfo{
				IP:    parts[0] + ":5577",
				ID:    parts[1],
				Model: parts[2],
			})
		}
	}()

	// Send message that will make bulbs respond
	_, err = socket.WriteToUDP([]byte(magicScanConst), &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 48899,
	})
	if err != nil {
		return nil, fmt.Errorf("error trying to request bulb information: %s", err)
	}

	<-done
	return infos, nil
}
