package service

import (
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"time"

)


func (p *RequestProcessor) CreateGame(data []byte) (resp []byte, err error) {

// Validate
	var req models.ReqCreateGame
	err = json.Unmarshal(data, &req)
	if err != nil {
		logger.ZapLogger.Errorw("validation Error", "Error", err)
		return
	}
	p.logger = p.logger.With(
				"MsgId", req.MsgId, 
				"GameId", req.GameId,
			)
	p.logger.Infow("Request",
		"body", string(data),
	)
//ProcMsg
	player := req.Player
	var generalResp *models.GeneralResp
	var respCreateGame models.RespCreateGame

	p.logger.Infof("Match Id: %v", req.MatchId)
	generalResp = &models.GeneralResp{
		Message: "Game Created",
		Code: "200",
		Status: "SUCCESS",

	}
	respCreateGame.GameId = req.GameId
	respCreateGame.MsgId = req.MsgId
	respCreateGame.Timestamp = time.Now().Format(time.DateTime)

	err = p.db.InsertGame(req.GameId, req.MatchId)
	if err != nil {
		p.logger.Errorw("DB Error", "Error", err)
		return	
	}

	p.logger.Infow("Player Info", "PlayerID", player.PlayerId, "Name", player.Name, "Position", player.Pos)
	err = p.db.InsertPlayer(player, req.GameId)
	if err != nil {
		p.logger.Errorw("DB Error", "Error", err)
		return
	}

	respCreateGame.GeneralResp = generalResp
	resp, err = json.Marshal(respCreateGame)
	if err != nil {
		return 
	}


	return
} 	