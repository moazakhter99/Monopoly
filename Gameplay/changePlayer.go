package gameplay

import (
	db "Monopoly/DB"
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"errors"
)

type ChangePlayerProc struct {
	db   db.DbOperations
	room gameroom.Room
}

func CreateChangePlayer(db db.DbOperations, room gameroom.Room) *ChangePlayerProc {
	return &ChangePlayerProc{
		db:   db,
		room: room,
	}
}

func (c *ChangePlayerProc) Validate(reqMsg []byte, param map[string]string) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Validate Change Player")
	var req models.Request
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "Validation Error", err)
		return
	}
	if param["Player"] != c.room.GetCurrentPlayer(param["Game"]) {
		logger.ZapLogger.Errorf("Player %v is not playing", param["Player"])
		logger.ZapLogger.Infoln("Exit Validate Move Pos")
		return nil, errors.New("Not Playing")
	}

	logger.ZapLogger.Infoln("Exit Validate Change Player")
	return req, err
}

// Might not work as the resposens are targeted
func (c *ChangePlayerProc) Play(payload any, param map[string]string) (targetMap map[string]any, err error) {
	logger.ZapLogger.Infoln("Enter Play Change Player")
	// req := payload.(models.Request)

	gameId := param["Game"]
	playerId := param["Player"]
	targetMap = make(map[string]any, 2)

	seq, count, err := c.db.GetPlayerSeqAndCount(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.CHANGEPLAYER, "Curr Player", playerId, "CurrSeq", seq, "currCount", count)
	nextPlayerId, err := c.db.GetNextPlayer(gameId, nextSeq(seq, count))
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.CHANGEPLAYER, "Next Player", nextPlayerId)

	nextPlayer := models.RespChangePlayer{
		NextPlayer: nextPlayerId,
		Playing:    true,
	}
	targetMap[nextPlayerId] = nextPlayer

	currPlayer := models.RespChangePlayer{
		NextPlayer: nextPlayerId,
		Playing:    false,
	}
	targetMap[playerId] = currPlayer
	c.room.UpdateGameState(gameId, nextPlayerId, models.CHANGEPLAYER)

	logger.ZapLogger.Infoln("Exit Play Change Player")
	return
}

// Response implements [Game].
func (c *ChangePlayerProc) Response(targetMap map[string]any, reqParam map[string]string, readChan chan []byte) (err error) {
	logger.ZapLogger.Infoln("Enter Change Player Response")

	gameId := reqParam["Game"]
	clientList := c.room.GetClientListByGame(gameId)

	for id, respMsg := range targetMap {
		client, ok := clientList[id]
		if ok {
			resp, err := json.Marshal(respMsg)
			if err != nil {
				logger.ZapLogger.Errorw("JSON Error", "Error", err)
				return err
			}
			logger.ZapLogger.Infow(models.CHANGEPLAYER, "Payload", string(resp))
			logger.ZapLogger.Infow(models.CHANGEPLAYER, "Game", gameId, "Clinet Count", len(clientList))
			wsMessage := models.WSMessage{
				Type: models.CHANGEPLAYER,
				Payload: resp,
			}

			wsResp, err := json.Marshal(wsMessage)
			if err != nil {
				logger.ZapLogger.Errorw("JSON Error", "Error", err)
				return err
			}
			client.WriteMsg <- wsResp

		}
	}

	logger.ZapLogger.Infoln("Exit Change Player Response")
	return
}
