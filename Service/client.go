package service

import (
	models "Monopoly/Models"
	// service "Monopoly/Service"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)




type Client struct {
	logger *zap.SugaredLogger
	Conn *websocket.Conn
	ReadMsg chan models.WSMessage
	WriteMsg chan models.WSMessage
	ErrMsg chan string
	PlayerId string
	// gameHub service.GameHubProcessor
	
}

func CreateNewClient(playerId string, conn *websocket.Conn, logger *zap.SugaredLogger) *Client {
	return &Client{
		logger: logger,
		ReadMsg: make(chan models.WSMessage),
		WriteMsg: make(chan models.WSMessage),
		ErrMsg: make(chan string),
		PlayerId: playerId,
		Conn: conn,
		// gameHub: hub,
	}
}


func (c *Client) ReadMessage() {
	c.logger.Infoln("Enter Read WS Message")
	var message models.WSMessage
	defer func ()  {
		c.ErrMsg <- c.PlayerId
	}() 

	for {

		err := c.Conn.ReadJSON(&message)
		if err != nil {
			c.logger.Errorw("Client Read Message", "Error", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Errorw("Unexpected Closer", "Error", err)
			}
			break
		}

		c.ReadMsg <- message
		hub := GameHub{
			logger: c.logger,
		}
		hub.ProcessEvent(message)
	}
	c.logger.Infoln("Exit Read WS Message")
}

func (c *Client) WriteMessage() {
	c.logger.Infoln("Enter Write WS Message")
	defer func ()  {
		c.ErrMsg <- c.PlayerId
	}()

	for {
		message, ok := <- c.WriteMsg
		if !ok {
			wsError := models.WsError{
				Message: "",
				WsError: websocket.CloseMessage,
			}
			err := c.Conn.WriteJSON(wsError)
			if err != nil {

			}

			break
		}
		
		if err := c.Conn.WriteJSON(message); err != nil {
			c.logger.Errorw("Chan Write Message", "Error", err)
		}

	}

	c.logger.Infoln("Exit Write WS Message")	
}