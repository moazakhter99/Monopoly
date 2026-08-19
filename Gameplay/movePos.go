package gameplay

import (
	db "Monopoly/DB"
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"strconv"
)

type MovePosProc struct {
	db   db.DbOperations
	room gameroom.Room
}

func CreateMovePos(db db.DbOperations, room gameroom.Room) *MovePosProc {
	return &MovePosProc{
		db:   db,
		room: room,
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

func (m *MovePosProc) Play(payload any, param map[string]string) (targetMap map[string]any, err error) {
	logger.ZapLogger.Infoln("Enter Play Move Pos")
	var response models.RespMovePos
	req := payload.(models.ReqMovePos)
	targetMap = make(map[string]any, 2)

	gameId := param["Game"]
	playerId := param["Player"]

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
				BlockId:   block.BlockId,
				NewPos:    newPos,
				Type:      block.Type,
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

		} else {
			response = models.RespMovePos{
				BlockId: block.BlockId,
				NewPos:  newPos,
				Type:    block.Type,
			}

		}
	}

	targetMap[""] = response
	m.room.UpdateGameState(gameId, playerId, models.ROLLDICE)

	logger.ZapLogger.Infoln("Exit Play Move Pos")
	return
}

// Response implements [Game].
func (m *MovePosProc) Response(targetMap map[string]any, reqParam map[string]string, readChan chan []byte) (err error) {
			// go callActionCard(client, readCh)
	panic("unimplemented")
}
