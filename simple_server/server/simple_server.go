package main

import (
	"fmt"
	"net"
	"os"
)

const PORT = "15354"

func main() {
	ln, err := net.Listen("tcp", "[2600:2d00:ffb0:3b3:42a8:308::]:"+PORT)
	if err != nil {
		fmt.Println("error in listening on socket:", err)
		os.Exit(1)
	}
	fmt.Println("listening on socket!")
	conn, err := ln.Accept()
	defer conn.Close()
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
	n, err = conn.Write(readBuf[0:n])
	if err != nil {
		fmt.Println("unable to write to peer:", err)
	}
}
