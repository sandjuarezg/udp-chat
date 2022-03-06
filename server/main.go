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

	n, addr, err := conn.ReadFromUDP(reply)
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, err = conn.WriteTo([]byte(fmt.Sprintf(" - Welcome %s - \n", reply[:n-1])), addr)
		if err != nil {
			log.Fatal(err)
		}

		addrAux, nAux, err := handleRequest(conn, addr, reply)
		if err != nil {
			log.Fatal(err)
		}

		addr = addrAux
		n = nAux
	}
}

// handleRequest Handle client request
//  @param1 (conn): connection between client and server
//  @param2 (addr): address of a UDP end point
//  @param3 (reply): buffer of reply
//
//  @return1 (err): error variable
func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, reply []byte) (addrAux *net.UDPAddr, n int, err error) {
	for {
		n, addrAux, err = conn.ReadFromUDP(reply)
		if err != nil {
			log.Fatal(err)
		}

		if addr.String() != addrAux.String() {
			return
		}

		_, err = conn.WriteTo([]byte(fmt.Sprintf("-> %s", reply[:n])), addrAux)
		if err != nil {
			log.Fatal(err)
		}
	}
}
