package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Invalid input: [host] [port]")
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Listening on", conn.LocalAddr())

	reply := make([]byte, 1024)

	_, addr, err = conn.ReadFromUDP(reply)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.WriteTo([]byte("hola mundo"), addr)
	if err != nil {
		log.Fatal(err)
	}
}
