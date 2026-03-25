package service

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"

	"go.uber.org/zap"
)



type JoinGame struct {
	db db.DbOperations
	logger *zap.SugaredLogger
}

func CreateJoinGameReq(db db.DbOperations, logger *zap.SugaredLogger) *JoinGame {
	return &JoinGame{
		db: db,
		logger: logger,
	}
}


func (p *JoinGame) Validate(data []byte) (req any, err error) {
	logger.ZapLogger.Infoln("Enter JoinGame Validation")
	var request models.ReqJoinGame
	err = json.Unmarshal(data, &request)
	if err != nil {
		logger.ZapLogger.Errorw("Validation Error", "Error", err)
		return
	}
	p.logger = p.logger.With(
				"MsgId", request.MsgId, 
			)
	p.logger.Infow("Request",
		"body", string(data),
	)

	logger.ZapLogger.Infoln("Exit JoinGame Validation")
	return &request, err
}


func (p *JoinGame) ProcessMsg(body any) (resp []byte, err error) {

	req := body.(*models.ReqJoinGame)
	player := req.Player
	var generalResp *models.GeneralResp
	var respJoinGame models.RespJoinGame

	p.logger.Infof("Match Id: %v", req.MatchId)
	respJoinGame.Joined = "TRUE"
	generalResp = &models.GeneralResp{
		Message: "",
		Code: "200",
		Status: "SUCCESS",

	}

	gameId, err := p.db.GetGameFromMatchId(req.MatchId)
	if err != nil {
		logger.ZapLogger.Errorw("DB Error", "Error", err)
		p.logger.Infow("Could not found the Game with this Match Id", "MatchId", req.MatchId)

		respJoinGame.Joined = "FALSE"
		respJoinGame.GeneralResp = generalResp
		resp, err = json.Marshal(respJoinGame)
		if err != nil {
			return nil, err
		}

		return
	}

	p.logger.Infow("Player Joined", "PlayerID", player.PlayerId, "Name", player.Name, "Position", player.Pos)
	err = p.db.InsertPlayer(player, gameId)
	if err != nil {
		p.logger.Errorw("DB Error", "Error", err)
		return
	}

	respJoinGame.GeneralResp = generalResp
	resp, err = json.Marshal(respJoinGame)
	if err != nil {
		return nil, err
	}
	return
}