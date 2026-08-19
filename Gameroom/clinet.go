package gameroom

import (
	models "Monopoly/Models"
	"Monopoly/logger"
	"Monopoly/router"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type ClinetProcessor interface {
	ReadMessage()
	WriteMessage()
}

type NewClient struct {
	PlayerId  string
	GameId    string
	ReadMsg   chan []byte
	WriteMsg  chan []byte
	Conn      *websocket.Conn
	Server    router.GpServer
	ConnClose func()
	room      Room
}

func CreateNewOtherClinet(r router.Router, conn *websocket.Conn, room Room, playerId, gameId string) *NewClient {
	s := router.NewServer(r)
	c := NewClient{
		ReadMsg:  make(chan []byte),
		WriteMsg: make(chan []byte),
		Conn:     conn,
		Server:   s,
		PlayerId: playerId,
		GameId:   gameId,
		room:     room,
	}
	c.room.AddPlayer(gameId, c)
	return &c
}

func (c *NewClient) ReadMessage() {
	logger.ZapLogger.Infoln("Enter Reading WS Message for ", c.PlayerId)
	defer func() {
		c.room.RemovePlayer(c.PlayerId, c.GameId)
		c.Conn.Close()
	}()

	for {
		var message models.WSMessage
		ReqParam := make(map[string]string, 2)

		err := c.Conn.ReadJSON(&message)
		if err != nil {
			logger.ZapLogger.Errorw("Client Read Message", "Error", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.ZapLogger.Errorw("Unexpected Closer", "Error", err)
			}
			break
		}

		logger.ZapLogger.Infow("Read", "Msg", message)
		ReqParam["Player"] = c.PlayerId
		ReqParam["Game"] = c.GameId

		err = c.Server.Write(message.Type, message.Payload, ReqParam, nil)
		if err != nil {
			// Write this error
			logger.ZapLogger.Errorw("WS Message Router", "Error", err)
			return
		}

	}
	logger.ZapLogger.Infoln("Exit Reading WS Message")

}

func (c *NewClient) WriteMessage() {
	logger.ZapLogger.Infoln("Enter Writing WS Message for ", c.PlayerId)
	defer func() {
		c.room.RemovePlayer(c.PlayerId, c.GameId)
		c.Conn.Close()
	}()

	for {
		var message models.WSMessage
		msgByte, ok := <-c.WriteMsg
		if !ok {
			wsError := models.WsError{
				Message: "Error While Reding Response",
				WsError: websocket.CloseMessage,
			}
			err := c.Conn.WriteJSON(wsError)
			if err != nil {
				// Write this error
				break
			}
			return
		}

		logger.ZapLogger.Infoln("Write Msg: ", string(msgByte))
		err := json.Unmarshal(msgByte, &message)
		if err != nil {
			wsError := models.WsError{
				Message: err.Error(),
				WsError: websocket.CloseMessage,
			}
			err := c.Conn.WriteJSON(wsError)
			if err != nil {
				// Write this error
				break
			}
			return
		}
		if err := c.Conn.WriteJSON(message); err != nil {
			logger.ZapLogger.Errorw("Chan Write Message", "Error", err)
		}

	}

	logger.ZapLogger.Infoln("Exit Write WS Message")

}
