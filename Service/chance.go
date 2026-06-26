package service

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)





type ChanceReq struct {
	db db.DbOperations
}

func CreateChanceReq(db db.DbOperations) *ChanceReq {
	return &ChanceReq{
		db: db,
	}
}


func (p *ChanceReq) Validate(data []byte) (req any, err error) {
	logger.ZapLogger.Infoln("Enter Chance Validation")
	var request models.ReqBlockInfo
	err = json.Unmarshal(data, &request)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANCE, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Chance Validation")
	return
}

func (p *ChanceReq) ProcessMsg(body any) (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Chance Processor")

	req := body.(models.ReqBlockInfo)
	var infoList []*models.BlockInfo

	ChanceInfo, err := p.db.GetBlockInfoByBlockType(models.CHANCE)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANCE, "DB Error", err)
		return
	}

	for _, info := range ChanceInfo {
		blockInfo := models.BlockInfo{
			InfoId: info.InfoId,
			Info: info.Info,
		}

		infoList = append(infoList, &blockInfo)
	}

	ChanceResp := models.RespBlockInfo{
		MsgId: req.MsgId,
		BlockType: models.CHANCE,
		BlockInfo: infoList,
	}

	resp, err = json.Marshal(ChanceResp)
	if err != nil {
		logger.ZapLogger.Errorw(models.CHANCE, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infoln("Enter Chance Processor")
	return
}