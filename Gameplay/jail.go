package gameplay

import (
	db "Monopoly/DB"
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)

type JailProc struct {
	db   db.DbOperations
	room gameroom.Room
}

func CreateJail(db db.DbOperations, room gameroom.Room) *JailProc {
	return &JailProc{
		db:   db,
		room: room,
	}
}

func (j *JailProc) Validate(reqMsg []byte) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Validate Jail")
	var req models.ReqJail
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.JAIL, "Validation Error", err)
		return
	}
	logger.ZapLogger.Infoln("Exit Validate Jail")
	return req, err
}

func (j *JailProc) Play(payload any, param map[string]string) (targetMap map[string]any, err error) {
	logger.ZapLogger.Infoln("Enter Jail Play")
	req := payload.(models.ReqJail)
	targetMap = make(map[string]any, 2)
	var updatedCash int
	var inJail bool

	gameId := param["Game"]
	playerId := param["Player"]

	if !req.InJail {
		return
	}

	cash, pos, err := j.db.GetPlayerCashPos(playerId)
	if err != nil {
		logger.ZapLogger.Errorw(models.ACTIONCARD, "DB Error", err)
		return
	}

	switch req.JailId {

	case "Jail1":
		updatedCash = cash - 500
		logger.ZapLogger.Infow(models.JAIL, "Updated Cash", updatedCash)
		err = j.db.UpdatePlayerCash(playerId, updatedCash, pos)
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

		inJail = false

		err = j.db.UpdatePlayerStatus(playerId, "")
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

	case "Jail2":
		err = j.db.DeleteGetOutOfJailCard(playerId, gameId, "Special Card0")
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}
		inJail = false
		updatedCash = cash

		err = j.db.UpdatePlayerStatus(playerId, "")
		if err != nil {
			logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
			return
		}

	case "Jail3":
		// The player is blocked at the front end and they are also checked at backend
		inJail = req.InJail
		updatedCash = cash

	}

	response := models.RespJail{
		Cash:   updatedCash,
		InJail: inJail,
	}

	targetMap[""] = response
	j.room.UpdateGameState(gameId, playerId, models.ROLLDICE)

	logger.ZapLogger.Infoln("Exit Jail Play")
	return
}

// Response implements [Game].
func (j *JailProc) Response(targetMap map[string]any, reqParam map[string]string, readChan chan []byte) (err error) {
	logger.ZapLogger.Infoln("Enter Jail Response")

	gameId := reqParam["Game"]
	playerId := reqParam["Player"]
	clientList := j.room.GetClientListByGame(gameId)
	respMsg := targetMap[""]
	resp, err := json.Marshal(respMsg)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	logger.ZapLogger.Infow(models.JAIL, "Game", gameId, "Clinet Count", len(clientList))
	wsMessage := models.WSMessage{
		Type: models.JAIL,
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

	logger.ZapLogger.Infoln("Exit Jail Response")
	return
}
