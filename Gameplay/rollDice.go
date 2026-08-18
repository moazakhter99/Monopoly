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
	client    *gameroom.NewClient
	room      *gameroom.GameRoom
	clientMap map[string]*gameroom.NewClient
	writeCh   chan<- models.WSMessage
}

func CreateRollDice(db db.DbOperations, room *gameroom.GameRoom) *RollDiceProc {
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
	logger.ZapLogger.Infow(models.ROLLDICE, "Dice Value", diceVal)

	targetMap[""] = response
	g.room.UpdateGameState(gameId, playerId, models.ROLLDICE)

	logger.ZapLogger.Infoln("Exit Roll Dice Player")
	return
}

func (g *RollDiceProc) Response(targetMap map[string]any, param map[string]string, readChan chan []byte) (err error) {
	logger.ZapLogger.Infoln("Enter Roll Dice Response")
	respMsg := targetMap[""]
	// go callMovePos(respMsg, g.client, param)

	resp, err := json.Marshal(respMsg)
	if err != nil {
		logger.ZapLogger.Errorw("JSON Error", "Error", err)
		return
	}

	go func() {
		
		diceValResp := respMsg.(models.RespDiceRoll)
		newPosReq := models.ReqMovePos{
			UpdateBy: diceValResp.DiceVal,
		}

		req, err := json.Marshal(newPosReq)
		if err != nil {
			logger.ZapLogger.Infow(models.ROLLDICE, "Json Error", err)
			return
		}

		err = g.client.Server.Write(models.MOVEPOS, req, param, nil)
		if err != nil {
			logger.ZapLogger.Errorw("WS Message Router", "Error", err)
		}
		
	}()

	for _, client := range g.room.GamePlayerMap[param["Game"]] {
		client.WriteMsg <- resp
	}

	logger.ZapLogger.Infoln("Exit Roll Dice Response")
	return
}
