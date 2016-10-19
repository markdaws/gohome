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
