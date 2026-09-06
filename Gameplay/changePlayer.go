package gameplay

import (
	db "Monopoly/DB"
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
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

func (c *ChangePlayerProc) Validate(reqMsg []byte) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Validate Change Player")
	var req models.Request
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "Validation Error", err)
		return
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
	logger.ZapLogger.Infow(models.CHANGEPLAYER, "CurrSeq", seq, "currCount", count)
	nextPlayerId, err := c.db.GetNextPlayer(gameId, nextSeq(seq, count))
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "DB Error", err)
		return
	}

	nextPlayer := models.RespChangePlayer{
		NextPlayer: nextPlayerId,
		Playing:    true,
	}
	nextResp, err := json.Marshal(nextPlayer)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "JSON Error", err)
		return
	}
	targetMap[nextPlayerId] = nextResp

	currPlayer := models.RespChangePlayer{
		NextPlayer: nextPlayerId,
		Playing:    false,
	}
	currResp, err := json.Marshal(currPlayer)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "JSON Error", err)
		return
	}
	targetMap[playerId] = currResp

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
