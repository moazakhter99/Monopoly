package service

import (
	db "Monopoly/DB"

	"go.uber.org/zap"
)


type RequestProcessor struct {
	db db.DbOperations
	logger *zap.SugaredLogger
}

func CreateNewRequestProcessor(db db.DbOperations, logger *zap.SugaredLogger) RequestProcessor {
	return RequestProcessor{
		db: db,
		logger: logger,
	}
}


type Processor interface {
	Validate(data []byte) (req any, err error)
	ProcessMsg(req any) (resp []byte, err error)
}

type GameHubProcessor interface {
	// ReadMessage()
	// WriteMessage()
	ProcessEvent(message any)
}


