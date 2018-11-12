package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:9999")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	conn.Write([]byte("Hello from client!"))
	conn.Close()
}
