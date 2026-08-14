package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)



type ChangePlayerProc struct {
	db 		db.DbOperations
	client 	*models.Client
	writeCh chan<- models.WSMessage

}

func CreateChangePlayer(db db.DbOperations) *ChangePlayerProc {
	return &ChangePlayerProc{
		db: db,
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
func (c *ChangePlayerProc) Play(payload any) (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Play Change Player")
	// req := payload.(models.Request)
	targetMap := make(map[string][]byte, 2)

	playerId := c.client.PlayerId
	gameId := c.client.GameId

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
		Playing: true,
	}
	nextResp, err := json.Marshal(nextPlayer)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANGEPLAYER, "JSON Error", err)
		return
	}
	targetMap[nextPlayerId] = nextResp

	currPlayer := models.RespChangePlayer{
		NextPlayer: nextPlayerId,
		Playing: false,
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
