package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)

type BuyBlockProc struct {
	db      db.DbOperations
	client  *models.Client
	writeCh chan<- models.WSMessage
}

func CreateBuyBlock(db db.DbOperations) *BuyBlockProc {
	return &BuyBlockProc{
		db: db,
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

func (b *BuyBlockProc) Play(payload any) (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Play Buy Block")
	var updatedCash int
	var buy bool
	req := payload.(models.ReqBuyBlock)

	playerId := ""
	gameId := ""
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
		BlockId:  req.BlockId,
		Buy:      buy,
		Cash:     updatedCash,
	}

	err = b.db.UpdatePlayerCash(playerId, updatedCash, pos)
	if err != nil {
		logger.ZapLogger.Errorw(models.BUYBLOCK, "DB Error", err)
		return
	}

	resp, err = json.Marshal(response)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Play Buy Block")
	return
}
