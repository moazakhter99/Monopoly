package gameplay

import (
	db "Monopoly/DB"
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)

type BuyBlockProc struct {
	db   db.DbOperations
	room gameroom.Room
}

func CreateBuyBlock(db db.DbOperations, room gameroom.Room) *BuyBlockProc {
	return &BuyBlockProc{
		db:   db,
		room: room,
	}
}

func (b *BuyBlockProc) Validate(reqMsg []byte) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Validate Buy Block")
	var req models.ReqBuyBlock
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "Validation Error", err)
		return
	}
	logger.ZapLogger.Infoln("Exit Validate Buy Block")
	return req, err
}

func (b *BuyBlockProc) Play(payload any, param map[string]string) (targetMap map[string]any, err error) {
	logger.ZapLogger.Infoln("Enter Play Buy Block")
	var updatedCash int
	var buy bool
	req := payload.(models.ReqBuyBlock)
	targetMap = make(map[string]any, 2)

	gameId := param["Game"]
	playerId := param["Player"]
	buy = false

	cash, pos, err := b.db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}

	price, err := b.db.GetBlockPrice(req.BlockId)
	// block, err := db.GetBlockState(req.Pos, gameId)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}

	if checkPlayerCash(cash, price) {
		updatedCash = cash - price
		buy = true
	} else {
		logger.ZapLogger.Infow("Player Low on Cash", "Player Id", playerId, "Cash", cash)
		buy = false

	}

	// Calculate Status for (colour, house, hotel)
	status := models.BASE

	err = b.db.UpdatePlayerCard(playerId, gameId, req.BlockId, status)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}

	logger.ZapLogger.Infow(models.BUYBLOCK, "Game Id", gameId, "Player Id", playerId, "Block Id", req.BlockId)
	response := models.RespBuyBlock{
		BlockId: req.BlockId,
		Buy:     buy,
		Cash:    updatedCash,
	}

	err = b.db.UpdatePlayerCash(playerId, updatedCash, pos)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}

	targetMap[""] = response
	b.room.UpdateGameState(gameId, playerId, models.BUYBLOCK)

	logger.ZapLogger.Infoln("Exit Play Buy Block")
	return
}

// Response implements [Game].
func (b *BuyBlockProc) Response(targetMap map[string]any, reqParam map[string]string, readChan chan []byte) (err error) {
	logger.ZapLogger.Infoln("Enter Buy Block Response")

	gameId := reqParam["Game"]
	playerId := reqParam["Player"]
	clientList := b.room.GetClientListByGame(gameId)
	respMsg := targetMap[""]
	resp, err := json.Marshal(respMsg)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	logger.ZapLogger.Infow(models.BUYBLOCK, "Game", gameId, "Clinet Count", len(clientList))
	wsMessage := models.WSMessage{
		Type: models.BUYBLOCK,
		Payload: resp,
	}

	wsResp, err := json.Marshal(wsMessage)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	go func() {
		cl, ok := clientList[playerId]
		if ok {
			changePlayerReq := models.Request{}
			req, err := json.Marshal(changePlayerReq)
			if err != nil {
				return
			}

			err = cl.Server.Write(models.CHANGEPLAYER, req, reqParam, readChan)
			if err != nil {
				logger.ZapLogger.Errorw("WS Message Router", "Error", err)
				return
			}
		}
	}()

	for _, client := range clientList {
		client.WriteMsg <- wsResp
	}

	logger.ZapLogger.Infoln("Exit Buy Block Response")
	return
}
