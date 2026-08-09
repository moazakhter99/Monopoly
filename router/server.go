package router

import (
	"Monopoly/logger"
	"encoding/json"
	"errors"
)


type GpServer interface {
	Write(action string, msg []byte) (err error)
	Read() (msg []byte)
	WriteJson(action string, msg any) (err error)
	ReadJson(msg any) (err error)

}

type server struct {
	route Router
	action	string
	msg []byte
}

type Request struct {
	Action string
	Msg		[]byte
	JsonMsg	any
}

func NewServer(r Router, ) *server {
	return &server{
		route: r,
	}
}

func (s server) Write(action string, msg []byte) (err error) {
	req := Request{
		Action: action,
		Msg: msg,
	}
	if err = s.validate(action); err != nil {
		return
	}
	go s.route.router(req)
	return
}

func (s server) Read() (msg []byte) {
	return s.route.readResponse()
}

func (s server) WriteJson(action string, msg any) (err error) {
	if err = s.validate(action); err != nil {
		return
	}
	// msgByte, err := json.Marshal(msg)
	// if err != nil {
	// 	return
	// }
	req := Request{
		Action: action,
		JsonMsg: msg,
	}
	go s.route.router(req)

	return
}

func (s server) ReadJson(msg any) (err error) {
	msgByte := s.route.readResponse()
	logger.ZapLogger.Infoln("Write Msg: ", string(msgByte))
	err = json.Unmarshal(msgByte, msg)
	if err != nil {
		return
	}
	return
}

func (s server) validate(action string) (err error) {

	if action == "" {
		return errors.New("action cannot be emtpy")
	}
	if !s.route.inMap(action) {
		return errors.New("action not available: " + action)
	}
	return
}