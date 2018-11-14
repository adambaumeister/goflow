package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type testPacket struct {
	Field int32
}

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:9999")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	tp := testPacket{64}
	err = binary.Write(conn, binary.LittleEndian, tp)
	if err != nil {
		fmt.Println("err:", err)
	}
}
