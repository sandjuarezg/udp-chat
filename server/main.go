package main

import (
	"fmt"
	"log"
	"net"
	"os"
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

	for {
		n, addr, err := conn.ReadFromUDP(reply)
		if err != nil {
			log.Fatal(err)
		}

		if addrExists(addr) {
			for _, element := range addrs {
				if element.String() != addr.String() {
					_, err = conn.WriteTo([]byte(fmt.Sprintf("-> %s", reply[:n])), element)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		} else {
			addrs = append(addrs, addr)

			name := reply[:n-1]
			fmt.Println(string(name), "connected")

			mess := fmt.Sprintf(" - %s connected - \n", name)
			mess += fmt.Sprintf(" - %d connected users - \n", len(addrs))

			for _, element := range addrs {
				_, err = conn.WriteTo([]byte(mess), element)
				if err != nil {
					log.Fatal(err)
				}
			}
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
