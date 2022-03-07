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
	var mess string

	for {
		n, addr, err := conn.ReadFromUDP(reply)
		if err != nil {
			log.Fatal(err)
		}

		name, ban := addrExists(addr)

		if ban {
			if string(reply[:n-1]) == "EXIT" {
				for n, element := range users {
					if element.addr.String() == addr.String() {
						users = append(users[:n], users[n+1:]...)
					}

					mess = fmt.Sprintf(" - Bye %s - \n", name)
					mess += fmt.Sprintf(" - %d connected users - \n", len(users))

					_, err = conn.WriteTo([]byte(mess), element.addr)
					if err != nil {
						log.Fatal(err)
					}
				}

				fmt.Println(name, "offline")
				continue
			}

			for _, element := range users {
				if element.addr.String() != addr.String() {
					_, err = conn.WriteTo([]byte(fmt.Sprintf("%s (%s): %s", name, time.Now().Format(time.RFC822Z), reply[:n])), element.addr)
					if err != nil {
						log.Fatal(err)
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
//  @return2 (ban): ban variable
func addrExists(addr *net.UDPAddr) (name string, ban bool) {
	for _, element := range users {
		if element.addr.String() == addr.String() {
			name = element.name
			ban = true
			return
		}
	}

	return
}
