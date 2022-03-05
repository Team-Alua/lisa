package main

import (
	"net"
	"bufio"
	"encoding/json"
	"fmt"
)

type SaveClient struct {
	Conn net.Conn
	Id string
	TitleId string
	DirName string
}

type SaveClientRequest struct {
	Name string `json:"name"`
	Arg1 string `json:"arg1,omitempty"` 
	Arg2 string `json:"arg2,omitempty"`
	Arg3 string `json:"arg3,omitempty"`
}

type SaveClientResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func (sc * SaveClient) ReadRequest() (SaveClientRequest, error) {
	jsonStr, err := bufio.NewReader(sc.Conn).ReadString('\n')
	req := SaveClientRequest{}
	if err != nil {
		return req, err
	}
	err = json.Unmarshal([]byte(jsonStr),&req) 
	return req, err
}

func (sc * SaveClient) SendError(msg string) error {
	resp := SaveClientResponse{Type: "Error", Msg: msg}
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Fprint(sc.Conn, string(b) + "\n")
	return nil

}

func (sc * SaveClient) GetSlotName() string {
	return sc.TitleId + "-" + sc.DirName 
}
