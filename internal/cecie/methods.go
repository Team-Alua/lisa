package cecie

import (
	"net"
	"bytes"
	"errors"
	"time"
	"io"
	"io/fs"
	"archive/zip"
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

func (c * Connection) ReceiveFileHeader() (string, uint64, error) {
	var rawBuffer bytes.Buffer

	_, err := io.CopyN(&rawBuffer, c.conn, 136)
	if err != nil {
		return "", 0, err
	}

	buffer := rawBuffer.Bytes()




	fileNameBuffer := buffer[:128] 

	fileNameEnd := bytes.IndexByte(fileNameBuffer, byte('\x00'))

	if fileNameEnd == -1 {
		fileNameEnd = 128
	}

	size := binary.LittleEndian.Uint64(buffer[128:])

	return string(fileNameBuffer[:fileNameEnd]), size,  nil
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
	_, err := io.CopyN(f, c.conn, size)
	if err != nil {
		return err
	}
	return nil
}

func (c * Connection) ReceiveContainerDump(zw * zip.Writer, rootFolder string) error {
	var fileCountBuffer bytes.Buffer
	_, err := io.CopyN(&fileCountBuffer, c.conn, 4)
	if err != nil {
		return err
	}
	fileCount := binary.LittleEndian.Uint32(fileCountBuffer.Bytes())

	for i := uint32(0); i < fileCount; i++ {
		filePath, size, err := c.ReceiveFileHeader()
		if err != nil {
			return err
		}

		fw, err := zw.Create(rootFolder + filePath)
		if err != nil {
			return err
		}
		_, err = io.CopyN(fw, c.conn, int64(size))

		if err != nil {
			return err
		}
	}
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
