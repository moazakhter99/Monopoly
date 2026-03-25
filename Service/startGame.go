package service

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"

	"go.uber.org/zap"
)

type StartGame struct {
	db db.DbOperations
	logger *zap.SugaredLogger
}

func CreateStartGameReq(db db.DbOperations, logger *zap.SugaredLogger) *StartGame {
	return &StartGame{
		db: db,
		logger: logger,
	}
}


func (p *StartGame) Validate(data []byte) (req any, err error) {
	logger.ZapLogger.Infoln("Enter StartGame Validation")
	var request models.RespStartGame
	err = json.Unmarshal(data, &request)
	if err != nil {
		p.logger.Errorln("Validation Error", "Err", err)
		return
	}

	p.logger = p.logger.With(
				"MsgId", request.MsgId, 
				"GameId", request.GameId,
			)
	p.logger.Infow("Request",
		"body", string(data),
	)

	logger.ZapLogger.Infoln("Exit StartGame Validation")
	return &request, err
}

func (p *StartGame) ProcessMsg(body any) (resp []byte, err error) {
	p.logger.Infoln("Enter StartGame")


	p.logger.Infoln("Exit StartGame")
	return
}