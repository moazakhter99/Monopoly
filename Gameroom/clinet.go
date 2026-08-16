package gameroom

import (
	models "Monopoly/Models"
	"Monopoly/logger"
	"Monopoly/router"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type NewClient struct {
	PlayerId  string
	GameId    string
	ReadMsg   chan []byte
	WriteMsg  chan []byte
	Conn      *websocket.Conn
	Server    router.GpServer
	ConnClose func()
}

type ClinetProcessor interface {
	// UpgradeClinet(playerId string, conn *websocket.Conn, logger *zap.SugaredLogger)
	ReadMessage()
	WriteMessage()
	UpgradeClinet(conn *websocket.Conn, PlayerId, GameId string) *NewClient
}

func CreateOtherClinet(r router.Router) *NewClient {
	s := router.NewServer(r)
	return &NewClient{
		ReadMsg:  make(chan []byte),
		WriteMsg: make(chan []byte),
		// Conn: conn,
		Server: s,
	}
}

func (c *NewClient) UpgradeClinet(conn *websocket.Conn, playerId, gameId string) *NewClient {
	c.Conn = conn
	c.PlayerId = playerId
	c.GameId = gameId
	return c
}

func (c *NewClient) ReadMessage() {
	logger.ZapLogger.Infoln("Enter Read WS Message ")
	defer func() {
		// Complete deregister
		// c.Conn.Close()
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

		clientDetail := models.Client{
			PlayerId: c.PlayerId,
			GameId:   "",
		}
		message.Client = &clientDetail
		logger.ZapLogger.Infow("Read", "Msg", message)
		ReqParam["Player"] = c.PlayerId
		ReqParam["Game"] = c.GameId

		err = c.Server.Write(message.Type, message.Payload, ReqParam, c.ReadMsg)
		if err != nil {
			// Write this error
			logger.ZapLogger.Errorw("WS Message Router", "Error", err)
			return
		}

	}
	logger.ZapLogger.Infoln("Exit Read WS Message")

}

func (c *NewClient) WriteMessage() {
	logger.ZapLogger.Infoln("Enter Write WS Message ")
	defer func() {
		// Complete deregister
		// c.Conn.Close()
	}()

	for {
		var message models.WSMessage
		msgByte := c.Server.Read()
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
		logger.ZapLogger.Infoln("Write Msg: ", message)
		if err := c.Conn.WriteJSON(message); err != nil {
			logger.ZapLogger.Errorw("Chan Write Message", "Error", err)
		}

	}

	logger.ZapLogger.Infoln("Exit Write WS Message")

}
