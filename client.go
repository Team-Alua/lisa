package main

import (

	"encoding/binary"
	"archive/zip"
	"net"
	"encoding/hex"
	"time"
	"math/rand"
	"fmt"
	"io"
	"os"
	"github.com/Team-Alua/lisa/internal/reservation"
	"github.com/Team-Alua/lisa/internal/client"
	"github.com/Team-Alua/lisa/internal/command"
	"github.com/Team-Alua/lisa/internal/validator"
	"github.com/Team-Alua/lisa/internal/cecie"

	"github.com/goccy/go-yaml"
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

func generateRandomFileName() string {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	var byteFileName [16]byte
	
	for i := 0; i < 16; i++ {
		byteFileName[i] = byte(generator.Intn(16))
	}
	return hex.EncodeToString(byteFileName[:])

}

func uploadZip(clientConn net.Conn, zipPath string) error {
	info, err := os.Stat(zipPath)
	if err != nil {
		return err
	}
	zipSize := info.Size()
	err = binary.Write(clientConn, binary.LittleEndian, &zipSize)
	if err != nil {
		return err
	}
	f, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_ , err = io.Copy(clientConn, f)
	return err
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


func handleCommandExecution(clientConn net.Conn, cecieIpPort string, rc * zip.ReadCloser, rr chan*reservation.Request) (string, error) {
	cRc, err := rc.Open("content.yml")
	if err != nil {
		return "", err
	}
	decoder := yaml.NewDecoder(cRc)
	var content client.Content
	if err = decoder.Decode(&content); err != nil {
		return "", err
	}

	responseChan := make(chan reservation.Response) 
	rs := reservation.Slot{Id: content.Target.TitleId + "-" + content.Target.DirectoryName, Out: responseChan}
	myRR := reservation.Request{Type: reservation.Add, Value: &rs} 
	rr <- &myRR

	for isReady := false; !isReady; {
		response := <- rs.Out
		switch(response.Type) {
		case reservation.Ready:
			fmt.Println("Ready!")
			// Put into slot array
			isReady = true
		case reservation.NotReady:
			fmt.Println("Not ready!")
			// Put into queue
		}
	}

	defer func() {
		removedRR := reservation.Request{Type: reservation.Remove, Value: &rs}
		rr <- &removedRR
	}()

	cc := cecie.Connection{}
	if err := cc.Connect(cecieIpPort); err != nil {
		fmt.Println(err)
		return "", err
	}
	defer cc.Close()


	zipName, err := command.Execute(&cc, &content, rc)
	
	if err != nil {
		if zipName != "" {
			// os.Remove(zipName)
		}
		return "", err
	}

	return zipName, nil
}

func handleRequest(clientConn net.Conn, rr chan*reservation.Request, cecieIpPort string) {
	defer clientConn.Close()

	zipName, err := downloadZip(clientConn)
	if err != nil {
		fmt.Println(err)
		return
	}
	rc, err :=  zip.OpenReader(zipName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(zipName)

	defer rc.Close()

	err = validator.CheckZip(rc)
	if err != nil {
		fmt.Println(err)
		return
	} 
	outZip, err := handleCommandExecution(clientConn, cecieIpPort,  rc, rr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(outZip)
	err = uploadZip(clientConn, outZip)
	if err != nil {
		fmt.Println(err)
	}
	// Send back response
}


