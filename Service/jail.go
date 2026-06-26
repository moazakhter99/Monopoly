package service

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)




type JailReq struct {
	db db.DbOperations
}


func CreateJailReq(db db.DbOperations) *JailReq {
	return &JailReq{
		db: db,
	}
}


func (p *JailReq) Validate(data []byte) (req any, err error) {
	logger.ZapLogger.Infoln("Enter Jail Validation")
	var request models.ReqJail
	err = json.Unmarshal(data, &request)
	if err != nil {
		logger.ZapLogger.Errorw(models.JAIL, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Jail Validation")
	return request, nil
}


func (p *JailReq) ProcessMsg(body any) (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Jail Processor")

	var infoList []*models.BlockInfo
	req := body.(models.ReqBlockInfo)

	jailInfo, err := p.db.GetBlockInfoByBlockType(models.JAIL)
	if err != nil {
		logger.ZapLogger.Errorw(models.JAIL, "DB Error", err)
		return
	}

	for _, info := range jailInfo {
		blockInfo := models.BlockInfo{
			InfoId: info.InfoId,
			Info: info.Info,
		}

		infoList = append(infoList, &blockInfo)
	}

	jailResp := models.RespBlockInfo{
		MsgId: req.MsgId,
		BlockType: models.JAIL,
		BlockInfo: infoList,
	}

	resp, err = json.Marshal(jailResp)
	if err != nil {
		logger.ZapLogger.Errorw(models.JAIL, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Jail Processor")
	return
}