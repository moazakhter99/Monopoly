package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)



type RollDiceProc struct {
	db 		db.DbOperations
	client 	*models.Client
	writeCh chan<- models.WSMessage

}

func CreateRollDice(db db.DbOperations) *RollDiceProc {
	return &RollDiceProc{
		db: db,
	}
}


func (g *RollDiceProc) Validate(reqMsg []byte) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Roll Dice Validate")
	var req models.ReqRolDice
	logger.ZapLogger.Infoln("Msg: ", string(reqMsg))
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.ROLLDICE, "Validation Error", err)
		return
	}
	logger.ZapLogger.Infoln("Exit Roll Dice Validate")
	return req, err
}


func (g *RollDiceProc) Play(payload any) (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Roll Dice Player")
	var diceVal int
	req := payload.(models.ReqRolDice)

	if req.Roll {
		diceVal = diceRoll()
	}
	
	response := models.RespDiceRoll{
		DiceVal:  diceVal,
	}

	resp, err = json.Marshal(response)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Roll Dice Player")
	return
}

// func (g *RollDiceProc) StateManagement(reqType, state string) (err error) {
// 	logger.ZapLogger.Infoln("Enter Roll Dice State Management")

// 	logger.ZapLogger.Infoln("Exit Roll Dice State Management")
// 	return
// }

// func (g *RollDiceProc) Response(respMsg any, clientMap map[string]any) (response []byte, err error) {
// 	logger.ZapLogger.Infoln("Enter Roll Dice Response")

// 	resp := respMsg.(models.RespDiceRoll)
// 	response, err = json.Marshal(resp)
// 	if err != nil {
// 		logger.ZapLogger.Errorw(models.ROLLDICE, "JSON Error", err)
// 		return
// 	}

// 	logger.ZapLogger.Infoln("Exit Roll Dice Response")
// 	return
// }