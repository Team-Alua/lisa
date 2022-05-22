package main

import (
	"net"
	"encoding/binary"
	"math/rand"
	"time"
	"encoding/hex"
	"os"
	"io"
	"fmt"
	"github.com/Team-Alua/lisa/internal/reservation"
	"github.com/Team-Alua/lisa/internal/validator"
)



func ListenForReservationRequests(r <-chan *reservation.Request) {
	rs := reservation.System{}
	rs.Initialize()

	for {
		rr := <-r
		if rr.Type == reservation.Add {
			o := rr.Value.Out
			if !rs.Add(rr.Value) {
				o <- reservation.Response{Type:reservation.NotReady, Msg: ""}
				rs.SortQueue()
			} else {
				o <- reservation.Response{Type:reservation.Ready, Msg: ""}
			}
		} else if rr.Type == reservation.Remove {
			if !rs.RemoveFromSlot(rr.Value) {
				rs.RemoveFromQueue(rr.Value)
				continue
			}

			c := rs.ChooseReservationFromQueue()
			// No candidates
			if c == nil {
				continue
			}

			rs.RemoveFromQueue(c)
			if !rs.Add(c) {
				panic("Something terrible happened")
			}
			c.Out <- reservation.Response{Type:reservation.Ready, Msg: ""}
			rs.SortQueue()
		}
	}

}

func Listen() {
	rr := make(chan *reservation.Request)

	go ListenForReservationRequests(rr)
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

		go handleRequest(conn, rr)
	}
}

func generateRandomFileName() string {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	var byteFileName [16]byte
	
	for i := 0; i < 16; i++ {
		byteFileName[i] = byte(generator.Intn(16))
	}
	return hex.EncodeToString(byteFileName[:])

}

func downloadZip(clientConn net.Conn) (string, error) {

	var zipSize int64
	err := binary.Read(clientConn, binary.LittleEndian, &zipSize)
	if err != nil {
		return "", err
	}
	zipName := generateRandomFileName() + ".zip"
	f, err := os.OpenFile(zipName, os.O_CREATE | os.O_WRONLY, 0744)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.CopyN(f, clientConn, zipSize)
	if err != nil {
		return "", err
	}

	return zipName, nil
}

func handleRequest(clientConn net.Conn, rr chan*reservation.Request) {
	defer clientConn.Close()

	zipName, err := downloadZip(clientConn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// defer os.Remove(zipName)

	err = validator.CheckZip(zipName)
	if err != nil {
		fmt.Println(err)
		return
	} 

	// Execute Command specified in zip
	// Send back response
}

