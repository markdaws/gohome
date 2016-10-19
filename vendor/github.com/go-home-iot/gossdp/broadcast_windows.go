/*
 * Copyright (c) 2015, fromkeith
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice, this
 *   list of conditions and the following disclaimer in the documentation and/or
 *   other materials provided with the distribution.
 *
 * * Neither the name of the fromkeith nor the names of its
 *   contributors may be used to endorse or promote products derived from
 *   this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
 * ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package gossdp


import (
    "encoding/binary"
    "fmt"
    "syscall"
    "unsafe"
    "net"
)


type theSocket struct {
    socket                  syscall.Handle
    readBytes               []byte
}

func (ts theSocket) IsValid() bool {
    return ts.socket != 0
}

func (s *Ssdp) createSocket() error {
    // create the socket
    var err error
    s.socket.socket, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
    if err != nil {
        return err
    }
    // make sure we can reuse it / share it
    if err := syscall.SetsockoptInt(s.socket.socket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil{
        syscall.Closesocket(s.socket.socket)
        s.socket.socket = 0
        return err
    }
    // going to broadcast
    if err := syscall.SetsockoptInt(s.socket.socket, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1); err != nil{
        syscall.Closesocket(s.socket.socket)
        s.socket.socket = 0
        return err
    }
    // bind it to the ssdp port
    lsa := &syscall.SockaddrInet4{Port: 1900, Addr: [4]byte{0, 0, 0, 0}}
    err = syscall.Bind(s.socket.socket, lsa)
    if err != nil {
        syscall.Closesocket(s.socket.socket)
        s.socket.socket = 0
        return err
    }
    iter, err := net.Interfaces()
    if err != nil {
        syscall.Closesocket(s.socket.socket)
        s.socket.socket = 0
        return err
    }
    wasFound := false
    for i := range iter {
        if iter[i].Flags & net.FlagMulticast == 0 {
            continue
        }
        addrs, err := iter[i].Addrs()
        if err != nil {
            continue
        }
        for k := range addrs {
            as4 := addrs[k].(*net.IPAddr).IP.To4()
            // join the multicast group
            mreq := &syscall.IPMreq{Multiaddr: [4]byte{239, 255, 255, 250}, Interface: [4]byte{as4[0], as4[1], as4[2], as4[3]}}
            if err := syscall.SetsockoptIPMreq(s.socket.socket, syscall.IPPROTO_IP, syscall.IP_ADD_MEMBERSHIP, mreq); err != nil {
                syscall.Closesocket(s.socket.socket)
                s.socket.socket = 0
                return err
            }
            wasFound = true
        }
    }
    // if we couldn't join a group, fall back to just 0.0.0.0
    if !wasFound {
        mreq := &syscall.IPMreq{Multiaddr: [4]byte{239, 255, 255, 250}, Interface: [4]byte{0, 0, 0, 0}}
        if err := syscall.SetsockoptIPMreq(s.socket.socket, syscall.IPPROTO_IP, syscall.IP_ADD_MEMBERSHIP, mreq); err != nil {
            syscall.Closesocket(s.socket.socket)
            s.socket.socket = 0
            return err
        }
    }

    s.socket.readBytes = make([]byte, 2048)

    return nil
}


func (s *Ssdp) closeSocket() {
    syscall.Closesocket(s.socket.socket)
    s.socket.socket = 0
}


func (s *Ssdp) read() ([]byte, string, error) {
    bufs := syscall.WSABuf{
        Len: 2048,
        Buf: &s.socket.readBytes[0],
    }
    var n, flags uint32
    var asIp4 syscall.RawSockaddrInet4
    fromAny := (*syscall.RawSockaddrAny) (unsafe.Pointer(&asIp4))
    fromSize := int32(unsafe.Sizeof(asIp4))
    err := syscall.WSARecvFrom(s.socket.socket, &bufs, 1, &n, &flags, fromAny, &fromSize, nil, nil)
    if err != nil {
        return nil, "", err
    }
    if n > 0 {
        // need to convert the port bytes ordering
        portBytes := make([]byte, 6)
        binary.BigEndian.PutUint16(portBytes, asIp4.Port)
        port := binary.LittleEndian.Uint16(portBytes)
        // set the address
        src := fmt.Sprintf("%d.%d.%d.%d:%d", asIp4.Addr[0], asIp4.Addr[1], asIp4.Addr[2], asIp4.Addr[3], port)
        //s.logger.Infof("Message: %s", string(readBytes[0:n]))
        return s.socket.readBytes[0:n], src, nil
    }
    return nil, "", nil
}



func (s *Ssdp) write(msg writeMessage) error {
    bufs := syscall.WSABuf{
        Len: uint32(len(msg.message)),
        Buf: &msg.message[0],
    }
    as4 := msg.to.IP.To4()
    to := &syscall.SockaddrInet4{
        Port: msg.to.Port,
        Addr: [4]byte{as4[0], as4[1], as4[2], as4[3]},
    }
    msgLen := uint32(len(msg.message))
    err := syscall.WSASendto(s.socket.socket, &bufs, 1, &msgLen, 0, to, nil, nil)
    return err
}
