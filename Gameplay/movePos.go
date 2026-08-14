package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"strconv"
)

type MovePosProc struct {
	db      db.DbOperations
	client  *models.Client
	writeCh chan<- models.WSMessage
}

func CreateMovePos(db db.DbOperations) *MovePosProc {
	return &MovePosProc{
		db: db,
	}
}

func (m *MovePosProc) Validate(reqMsg []byte) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Validate Move Pos")
	var req models.ReqMovePos
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "Validation Error", err)
		return
	}
	logger.ZapLogger.Infoln("Exit Validate Move Pos")
	return req, err
}

func (m *MovePosProc) Play(payload any) (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Play Move Pos")
	var response models.RespMovePos
	req := payload.(models.ReqMovePos)

	playerId := ""
	gameId := ""

	cash, pos, err := m.db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}
	logger.ZapLogger.Infow(models.MOVEPOS, "Current Position", pos)

	newPos := updatePos(req.UpdateBy, pos)
	// err = db.UpdatePlayerPos(playerId, newPos)
	err = m.db.UpdatePlayerCash(playerId, cash, newPos)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
	}

	logger.ZapLogger.Infow(models.MOVEPOS, "Player Id", playerId, "New Pos", newPos)
	block, err := m.db.GetBlockState(newPos, gameId)
	if err != nil {
		logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
		return
	}

	switch block.OwnerId {
	case "":
		// Buy or Action Card
		if block.Type == models.SPECIALCARD {
			var status string
			response = models.RespMovePos{
				BlockId:  block.BlockId,
				NewPos:   newPos,
				Type:     block.Type,
				BlockName: block.BlockName,
			}
			if block.BlockName == models.COMMUNITYCHEST || block.BlockName == models.CHANCE {
				response.CardNo = getCardNo()
				// cardInfo, err := db.GetBlockInfo(block.BlockName, resp.CardNo)
				// if err != nil {
				// 	logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
				// 	return
				// }
				// resp.CardInfo = cardInfo
				status = models.ACTIONCARD + "_" + strconv.Itoa(response.CardNo) + "_" + block.BlockName
			}

			status = models.ACTIONCARD + "_" + block.BlockName
			err = m.db.UpdatePlayerStatus(playerId, status)
			if err != nil {
				logger.ZapLogger.Errorw(models.MOVEPOS, "DB Error", err)
				return
			}
			// go callActionCard(client, readCh)

		} else {
			response = models.RespMovePos{
				BlockId:  block.BlockId,
				NewPos:   newPos,
				Type:     block.Type,
			}

		}
	}

	resp, err = json.Marshal(response)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Play Move Pos")
	return
}
