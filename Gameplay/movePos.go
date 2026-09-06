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

	case playerId:
		// Already owned by the player 
		response = models.RespMovePos{
			BlockId:  block.BlockId,
			NewPos:   newPos,
			State:    models.OWNED,
			Type:     block.Type,
			OwnerId:  playerId,
		}
		logger.ZapLogger.Infow(models.CHANGEPLAYER, "Current Player", playerId)

	default:
		// Pay rent
		response = models.RespMovePos{
			BlockId:  block.BlockId,
			NewPos:   newPos,
			State:    models.SOLD,
			Type:     block.Type,
			OwnerId:  block.OwnerId,
			Price: block.Price,
		}
	}

	targetMap[""] = response
	m.room.UpdateGameState(gameId, playerId, models.MOVEPOS)

	logger.ZapLogger.Infoln("Exit Play Move Pos")
	return
}

// Response implements [Game].
func (m *MovePosProc) Response(targetMap map[string]any, reqParam map[string]string, readChan chan []byte) (err error) {
	logger.ZapLogger.Infoln("Enter Move Pos Response")

	gameId := reqParam["Game"]
	playerId := reqParam["Player"]
	clientList := m.room.GetClientListByGame(gameId)
	respMsg := targetMap[""]
	resp, err := json.Marshal(respMsg)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	logger.ZapLogger.Infow(models.MOVEPOS, "Game", gameId, "Clinet Count", len(clientList))
	wsMessage := models.WSMessage{
		Type: models.ROLLDICE,
		Payload: resp,
	}

	wsResp, err := json.Marshal(wsMessage)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	go func() {
		movePosResp := respMsg.(models.RespMovePos)
		cl, ok := clientList[playerId]
		if !ok {
			logger.ZapLogger.Errorf("Client Not Found: %v", playerId)
			return
		} 

		ownerId := movePosResp.OwnerId
		if playerId == "" {
			return

		} else if playerId == ownerId {
			changePlayerReq := models.Request{}
			req, err := json.Marshal(changePlayerReq)
			if err != nil {
				return
			}

			err = cl.Server.Write(models.CHANGEPLAYER, req, reqParam, readChan)
			if err != nil {
				logger.ZapLogger.Errorw("WS Message Router", "Error", err)
			}

		} else if playerId != ownerId {
			calRentReq := &models.ReqCalculateRent{
				BlockId: movePosResp.BlockId,
				BlockType: movePosResp.Type,
				Price: movePosResp.Price,
				OwnerId: ownerId,
			}
			req, err := json.Marshal(calRentReq)
			if err != nil {
				logger.ZapLogger.Infow(models.MOVEPOS, "Json Error", err)
				return
			}

			err = cl.Server.Write(models.CALCULATERENT, req, reqParam, readChan)
			if err != nil {
				logger.ZapLogger.Errorw("WS Message Router", "Error", err)
			}
		}
	}()

	for _, client := range clientList {
		client.WriteMsg <- wsResp
	}

	logger.ZapLogger.Infoln("Exit Move Pos Response")
	return
}

