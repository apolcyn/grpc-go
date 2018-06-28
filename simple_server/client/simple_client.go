package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

const PORT = "15354"

func main() {
	s, err := syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	defer syscall.Close(s)
	if err != nill {
		fmt.Println("error creating socket:", err)
		os.Exit(1)
	}
	sa := &syscall.SockaddrInet6{Port: port, ZoneId: uint32(zoneCache.index(zone))}
	copy(sa.addr[:], net.IPv6loopback)
	if err := syscall.Connect(s, sa); err != nil {
		fmt.Println("error connecting on socket:", err)
		os.Exit(1)
	}
	writeBuf := []byte{1, 2, 3, 4}
	n, err = syscall.Write(s, writeBuf)
	if err != nil {
		fmt.Println("error writing on socket:", err)
		os.Exit(1)
	}	
	readBuf := make([]byte, 100)
	n, err = syscall.Read(s, readBuf)
	if err != nil {
		fmt.Println("error reading on socket:", err)
		os.Exit(1)
	}	
	fmt.Println("Received bytes:|%v|", readBuf[0:n])
	if err != nil {
		fmt.Println("unable to accept connection on socket:", err)
		os.Exit(1)
	}
	fmt.Println("accepted connection from:", conn.RemoteAddr())
	readBuf := make([]byte, 100)
	n, err := conn.Read(readBuf)
	if err != nil {
		fmt.Println("unable to read from conn:", err)
	}
	fmt.Println("read bytes from peer:|%v|", readBuf[0:n])
}
