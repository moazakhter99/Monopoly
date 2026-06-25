package service

import (
	models "Monopoly/Models"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)




type Client struct {
	PlayerId string
	GameId	 string
	logger *zap.SugaredLogger
	Conn *websocket.Conn
	ReadMsg chan models.WSMessage
	WriteMsg chan models.WSMessage
	ErrMsg chan string
	gameHub *GameHub
	
}

func CreateNewClient(playerId, gameId string, conn *websocket.Conn, logger *zap.SugaredLogger, hub *GameHub) *Client {
	return &Client{
		PlayerId: playerId,
		GameId: gameId,
		logger: logger,
		ReadMsg: make(chan models.WSMessage),
		WriteMsg: make(chan models.WSMessage),
		ErrMsg: make(chan string),
		Conn: conn,
		gameHub: hub,
	}
}


func (c *Client) ReadMessage() {
	c.logger.Infoln("Enter Read WS Message")
	var message models.WSMessage
	defer func ()  {
		c.gameHub.Deregister <- c
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

		clientDetail := models.Client{
			PlayerId: c.PlayerId,
			GameId: c.GameId,
		}
		message.Client = &clientDetail		

		c.gameHub.ReadMsg <- message

	}
	c.logger.Infoln("Exit Read WS Message")
}

func (c *Client) WriteMessage() {
	c.logger.Infoln("Enter Write WS Message")
	defer func ()  {
		c.gameHub.Deregister <- c
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