package main

import (
	"flag"
	"fmt"
	"io"
	"net"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "3000", "the tunnel server port")
}

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tunnel server listening on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(clientConn net.Conn) {
	fmt.Println("New client connected")

	// Wait for a connection from the tunnel client
	serverConn, err := net.Dial("tcp", "127.0.0.1:4000") // Change this to the actual tunnel client address
	if err != nil {
		fmt.Println("Error connecting to tunnel client:", err)
		clientConn.Close()
		return
	}

	fmt.Println("Connected to tunnel client")
	go copyIO(clientConn, serverConn)
	go copyIO(serverConn, clientConn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(dest, src)
}
