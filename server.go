package main

import (
	"net"
)

func Listen() {
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

		conn.Write([]byte("Hi!"))
		conn.Close()
	}
}

