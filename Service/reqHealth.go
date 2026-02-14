package service

import (
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)



func (p *RequestProcessor) Health() (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Health Processor")

	var dbStatus string
	dbStatus = "SUCCESS"

	err = p.db.Ping()
	if err != nil {
		dbStatus = "FAILED"
	}
	Msg := "Game Server is Running"
	serviceStatus := "SUCCESS"

	respBody := models.Health{
		Msg: Msg,
		DBStatus: dbStatus,
		ServiceStatus: serviceStatus,
	}

	resp, err = json.Marshal(respBody)

	logger.ZapLogger.Infoln("Exit Health Processor")

	return
}