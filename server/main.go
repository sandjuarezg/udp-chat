package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var addrs []*net.UDPAddr

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
	wg := sync.WaitGroup{}

	n, addr, err := conn.ReadFromUDP(reply)
	if err != nil {
		log.Fatal(err)
	}

	for {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			_, err = conn.WriteTo([]byte(fmt.Sprintf(" - Welcome %s - \n", reply[:n-1])), addr)
			if err != nil {
				log.Fatal(err)
			}

			addrs = append(addrs, addr)

			addrAux, nAux, err := handleRequest(conn, addr, reply)
			if err != nil {
				log.Fatal(err)
			}

			addr = addrAux
			n = nAux
		}(&wg)

		wg.Wait()
	}
}

// handleRequest Handle client request
//  @param1 (conn): connection between client and server
//  @param2 (addr): address of a UDP end point
//  @param3 (reply): buffer of reply
//
//  @return1 (addrAux): address aux of a UDP end point
//  @return2 (n): number of characters in buffer
//  @return3 (err): error variable
func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, reply []byte) (addrAux *net.UDPAddr, n int, err error) {
	for {
		n, addrAux, err = conn.ReadFromUDP(reply)
		if err != nil {
			log.Fatal(err)
		}

		if !addrExists(addrAux) {
			return
		}

		_, err = conn.WriteTo([]byte(fmt.Sprintf("-> %s", reply[:n])), addrAux)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// addrExists Check if address exists
//  @param1 (addr): address of a UDP end point
//
//  @return1 (ban): ban variable
func addrExists(addr *net.UDPAddr) (ban bool) {
	for _, element := range addrs {
		if element.String() == addr.String() {
			ban = true
			return
		}
	}

	return
}
