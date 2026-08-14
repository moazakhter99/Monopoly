package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)



type JailProc struct {
	db      db.DbOperations
	client  *models.Client
	writeCh chan<- models.WSMessage

}

func CreateJail(db db.DbOperations) *JailProc {
	return &JailProc{
		db: db,
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

func (j *JailProc) Play(payload any) (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Jail Play")
	req := payload.(models.ReqJail)
	var updatedCash int
	var inJail bool

	playerId := j.client.PlayerId
	gameId := j.client.GameId

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
		Cash: updatedCash,
		InJail: inJail,
	}

	// go callChangePlayer(client, readCh)

	resp, err = json.Marshal(response)
	if err != nil {
		logger.ZapLogger.Errorw(models.JAIL, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Jail Play")
	return
}