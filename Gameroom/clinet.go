package gameroom

import (
	models "Monopoly/Models"
	"Monopoly/logger"
	"Monopoly/router"
	"encoding/json"

	"github.com/gorilla/websocket"
)



type NewClient struct {
	PlayerId string
	ReadMsg chan models.WSMessage
	WriteMsg chan models.WSMessage
	Conn *websocket.Conn
	Server	router.GpServer
	ConnClose func()
}


type ClinetProcessor interface {
	// UpgradeClinet(playerId string, conn *websocket.Conn, logger *zap.SugaredLogger)
	ReadMessage()
	WriteMessage()
	UpgradeClinet(conn *websocket.Conn, PlayerId, GameId string)
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

func (c *NewClient) UpgradeClinet(conn *websocket.Conn, PlayerId, GameId string) {
	c.Conn = conn
	c.PlayerId = PlayerId
}

func (c *NewClient) ReadMessage() {
	logger.ZapLogger.Infoln("Enter Read WS Message ")
	defer func ()  {
		// Complete deregister
		// c.Conn.Close()
	}()

	for {
		var message models.WSMessage

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
			GameId: "",
		}
		message.Client = &clientDetail
		logger.ZapLogger.Infow("Read", "Msg", message)	

		err = c.Server.Write(message.Type, message.Payload)
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
	defer func ()  {
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