package gameplay

import (
	db "Monopoly/DB"
	gameroom "Monopoly/Gameroom"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)

type RollDiceProc struct {
	db        db.DbOperations
	room      gameroom.Room
}

func CreateRollDice(db db.DbOperations, room gameroom.Room) *RollDiceProc {
	return &RollDiceProc{
		db:   db,
		room: room,
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

func (g *RollDiceProc) Play(payload any, param map[string]string) (targetMap map[string]any, err error) {
	logger.ZapLogger.Infoln("Enter Roll Dice Player")
	var diceVal int
	targetMap = make(map[string]any)
	req := payload.(models.ReqRolDice)
	gameId := param["Game"]
	playerId := param["Player"]

	if req.Roll {
		diceVal = diceRoll()
	}

	response := models.RespDiceRoll{
		DiceVal: diceVal,
	}
	logger.ZapLogger.Infow(models.ROLLDICE, "Game", gameId, "Player", playerId, "Dice Value", diceVal)

	targetMap[""] = response
	g.room.UpdateGameState(gameId, playerId, models.ROLLDICE)

	logger.ZapLogger.Infoln("Exit Roll Dice Player")
	return
}

func (g *RollDiceProc) Response(targetMap map[string]any, param map[string]string, readChan chan []byte) (err error) {
	logger.ZapLogger.Infoln("Enter Roll Dice Response")

	gameId := param["Game"]
	playerId := param["Player"]
	clientList := g.room.GetClientListByGame(gameId)
	respMsg := targetMap[""]
	resp, err := json.Marshal(respMsg)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	logger.ZapLogger.Infow(models.ROLLDICE, "Game", gameId, "Clinet Count", len(clientList))
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
		
		diceValResp := respMsg.(models.RespDiceRoll)
		cl, ok := clientList[playerId]
		if !ok {
			logger.ZapLogger.Errorf("Client Not Found: %v", playerId)
			return
		} 

		newPosReq := models.ReqMovePos{
			UpdateBy: diceValResp.DiceVal,
		}

		req, err := json.Marshal(newPosReq)
		if err != nil {
			logger.ZapLogger.Infow(models.ROLLDICE, "Json Error", err)
			return
		}

		err = cl.Server.Write(models.MOVEPOS, req, param, readChan)
		if err != nil {
			logger.ZapLogger.Errorw("WS Message Router", "Error", err)
		}
		
	}()

	for key, client := range clientList {
		client.WriteMsg <- wsResp
	}

	logger.ZapLogger.Infoln("Exit Roll Dice Response")
	return
}
