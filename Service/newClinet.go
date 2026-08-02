package service

import (
	models "Monopoly/Models"
	"Monopoly/router"

	"github.com/gorilla/websocket"
)



type NewClient struct {
	ReadMsg chan models.WSMessage
	WriteMsg chan models.WSMessage
	Conn *websocket.Conn
	Server	router.GpServer
}


func CreateOtherClinet(r router.Router) *NewClient {
	s := router.NewServer(r)
	return &NewClient{
		ReadMsg: make(chan models.WSMessage),
		WriteMsg: make(chan models.WSMessage),
		// Conn: conn,
		Server: s,
	}
}


func (c *NewClient) ReadMessage() {

	var message models.WSMessage
	c.Conn.ReadJSON(message)

	err := c.Server.Write(message.Type, message.Payload)
	if err != nil {
		return
	}

}

func (c *NewClient) WriteMessage() {

	var message *models.WSMessage
	err := c.Server.ReadJson(message)
	if err != nil {
		return
	}
	
	c.Conn.WriteJSON(message)

}