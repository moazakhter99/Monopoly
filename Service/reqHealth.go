package service

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)



type HealthReq struct {
	db db.DbOperations
}


func CreateHealthReq(db db.DbOperations) *HealthReq {
	return &HealthReq{
		db: db,
	}
}


func (p *HealthReq) Validate(data []byte) (req any, err error) {
	logger.ZapLogger.Infoln("Enter Health Validation")

	logger.ZapLogger.Infoln("Exit Health Validation")
	return
}


func (p *HealthReq) ProcessMsg(req any) (resp []byte, err error) {
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