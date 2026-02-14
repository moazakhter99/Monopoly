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
