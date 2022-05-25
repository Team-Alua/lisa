package cecie

import (
	"net"
	"bytes"
	"errors"
	"fmt"
	"time"
	"io"
	"io/fs"
	"encoding/binary"
	"github.com/Team-Alua/lisa/internal/client"
)
func (c * Connection) Connect(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, 5 * time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c * Connection) Close() error {
	return c.conn.Close()
}

func (c * Connection) SendTarget(content * client.ContentTarget) error {
	var buffer [56]byte
	copy(buffer[0:16], content.TitleId)
	copy(buffer[16:48], content.DirectoryName)
	binary.LittleEndian.PutUint64(buffer[48:], content.Blocks)

	_, err := io.CopyN(c.conn, bytes.NewReader(buffer[:]), 56)
	return err
}

func (c * Connection) SendFileHeader(path string, size uint64) error {
	var buffer [136]byte
	copy(buffer[0:128], path)
	binary.LittleEndian.PutUint64(buffer[128:], size)

	_, err := io.CopyN(c.conn, bytes.NewReader(buffer[:]), 136)
	if err != nil {
		return err
	}

	return nil
}

func (c * Connection) SendFile(f fs.File) error {
	_, err := io.Copy(c.conn, f)
	if err != nil {
		return err
	}

	return nil
}

func (c * Connection) SendZipFile(f io.ReadCloser) error {
	_, err := io.Copy(c.conn, f)
	if err != nil {
		return err
	}

	return nil
}

func (c * Connection) ReceiveZipFile(f io.Writer, size int64) error {
	n, err := io.CopyN(f, c.conn, size)
	if err != nil {
		return err
	}
	fmt.Println("Dumped bytes:", n)
	return nil
}

func (c * Connection) CheckOkay() error {
	var r [64]byte

	_, err := io.ReadFull(c.conn, r[:])
	if err != nil {
		return err
	}
	s := string(bytes.Trim(r[:], "\x00"))

	if s != "ok" {
		return errors.New(s)
	}

	return nil
}

func (c * Connection) SendCommand(name string) error {
	var command [32]byte
	for i, char := range []byte(name) {
		command[i] = char
	}

	reader := bytes.NewReader(command[:])

	_, err := io.CopyN(c.conn, reader, 32)
	return err
}
