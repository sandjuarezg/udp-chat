package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// User structure
type user struct {
	addr *net.UDPAddr // address of UDP
	name string       // name of user
}

var users []user

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Invalid input: [port]")
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%s", os.Args[1]))
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
	var mess string

	for {

		n, addr, err := conn.ReadFromUDP(reply)
		if err != nil {
			log.Fatal(err)
		}

		name, exists := addrExists(addr)

		if exists {

			if string(reply[:n-1]) == "EXIT" {

				for n, element := range users {

					if element.addr.String() == addr.String() {
						users = append(users[:n], users[n+1:]...)
					}

					_, err = conn.WriteTo([]byte(fmt.Sprintf(" - Bye %s - \n", name)), element.addr)
					if err != nil {
						log.Fatal(err)
					}

				}

				fmt.Println(name, "offline")

			} else {

				for _, element := range users {

					if element.addr.String() != addr.String() {

						t := time.Now()
						_, err = conn.WriteTo([]byte(fmt.Sprintf("%s (%d:%d): %s", name, t.Hour(), t.Minute(), reply[:n])), element.addr)
						if err != nil {
							log.Fatal(err)
						}

					}

				}
			}

		} else {

			u := user{addr: addr, name: string(reply[:n-1])}
			users = append(users, u)

			fmt.Println(u.name, "connected")

			mess = fmt.Sprintf(" - %s connected - \n", u.name)
			mess += fmt.Sprintf(" - %d connected users - \n", len(users))

			for _, element := range users {

				_, err = conn.WriteTo([]byte(mess), element.addr)
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
//  @return1 (name): name of user
//  @return2 (exists): exists variable
func addrExists(addr *net.UDPAddr) (name string, exists bool) {
	for _, element := range users {

		if element.addr.String() == addr.String() {

			name = element.name
			exists = true
			return

		}

	}

	return
}
