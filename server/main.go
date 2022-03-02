package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// Socket structure of connections
type socket struct {
	conn net.UDPConn // connection of UDP
	addr net.UDPAddr // address of UDP
}

var sockets []socket

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

	for {
		reply := make([]byte, 1024)

		_, addr, err := conn.ReadFromUDP(reply)
		if err != nil {
			log.Fatal(err)
		}

		handleRequest(conn, addr)

		// fmt.Println(runtime.NumGoroutine())
	}
}

// handleRequest Handle client request
//  @param1 (conn): connection between client and server
//  @param2 (addr): address of a UDP end point
//
//  @return1 (err): error variable
func handleRequest(conn *net.UDPConn, addr *net.UDPAddr) (err error) {
	reply := make([]byte, 1024)

	mess := fmt.Sprintln(" - Welcome to chat - ")
	mess += fmt.Sprint("Enter your name: ")

	// write message
	_, err = conn.WriteTo([]byte(mess), addr)
	if err != nil {
		log.Fatal(err)
	}

	// read user name
	res := bufio.NewReader(conn)
	n, err := res.Read(reply)
	if err != nil {
		log.Fatal(err)
	}
	name := reply[:n-1]

	sockets = append(sockets, socket{conn: *conn, addr: *addr})
	fmt.Printf("%s connected\n", name)

	mess = fmt.Sprintf(" - %s connected - \n", name)
	mess += fmt.Sprintf(" - %d connected users - \n", len(sockets))

	// // write message to all connections
	for _, element := range sockets {
		_, err = element.conn.WriteToUDP([]byte(mess), &element.addr)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		reply = make([]byte, 1024)

		// read text to chat
		n, err = res.Read(reply)
		if err != nil {
			if err == io.EOF {
				// remove connection from chat
				for n, element := range sockets {
					if conn == &element.conn {
						sockets = append(sockets[:n], sockets[n+1:]...)
					}

					mess = fmt.Sprintf(" - Bye %s - \n", name)
					mess += fmt.Sprintf(" - %d connected users - \n", len(sockets)-1)

					_, err = element.conn.WriteToUDP([]byte(mess), &element.addr)
					if err != nil {
						log.Fatal(err)
					}
				}

				fmt.Printf("%s offline\n", name)

				break
			} else {
				log.Fatal(err)
			}
		}

		if string(reply[:n]) == "\n" {
			continue
		}

		//  write message to all connections
		for _, element := range sockets {
			if element.conn != *conn {
				_, err = element.conn.WriteToUDP([]byte(fmt.Sprintf("%s (%s): %s", name, time.Now().Format(time.RFC822Z), reply[:n])), &element.addr)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

	}

	return
}
