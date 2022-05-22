package main

import (
	"net"
	"github.com/Team-Alua/lisa/internal/reservation"
)


func main() {
	rr := make(chan *reservation.Request)
	go ListenForReservationRequests(rr)

	cecieIpPort := "10.0.0.5:8766"


	// Listen for incoming connections
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handleRequest(conn, rr, cecieIpPort)
	}
}

