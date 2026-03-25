package service

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type GameReq struct {
	db db.DbOperations
	logger *zap.SugaredLogger
}

func CreateGameReq(db db.DbOperations, logger *zap.SugaredLogger) *GameReq {
	return &GameReq{
		db: db,
		logger: logger,
	}
}


func (p *GameReq) Validate(data []byte) (req any, err error) {
	logger.ZapLogger.Infoln("Enter CreateGame Validation")
	var request models.ReqCreateGame
	err = json.Unmarshal(data, &request)
	if err != nil {
		logger.ZapLogger.Errorw("validation Error", "Error", err)
		return
	}
	p.logger = p.logger.With(
				"MsgId", request.MsgId, 
				"GameId", request.GameId,
			)
	p.logger.Infow("Request",
		"body", string(data),
	)

	logger.ZapLogger.Infoln("Exit CreateGame Validation")
	return &request, err
}


func (p *GameReq) ProcessMsg(body any) (resp []byte, err error) {

	req := body.(*models.ReqCreateGame)

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