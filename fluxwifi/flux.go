// Package fluxwifi provides methods to control flux WIFI bulbs
package fluxwifi

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// CREDIT: All the knowledge of how to control this product came from:
// https://github.com/beville/flux_led

//TODO: Turn on/off
//TODO: Get current level

// SetLevel changes the RGB values of the bulb. //TODO: conn parameter
func SetLevel(r, b, g byte, conn net.Conn) error {
	bytes := []byte{0x31, r, g, b, 0x00, 0xf0, 0x0f}
	var t int = 0
	for _, v := range bytes {
		t += int(v)
	}
	cs := t & 0xff
	bytes = append(bytes, byte(cs))
	_, err := conn.Write(bytes)
	return err
}

// BulbInfo contains the information returned from scanning the local
// network for Flux WIFI bulbs
type BulbInfo struct {
	IP    string
	ID    string
	Model string
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
	infos := make([]BulbInfo, 0)

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
				IP:    parts[0],
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
