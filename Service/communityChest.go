package service

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
)




type CommunityChestReq struct {
	db db.DbOperations
}

func CreateCommunityChestReq(db db.DbOperations) *CommunityChestReq {
	return &CommunityChestReq{
		db: db,
	}
}


func (p *CommunityChestReq) Validate(data []byte) (req any, err error) {
	logger.ZapLogger.Infoln("Enter Community Chest Validation")
	var request models.ReqBlockInfo
	err = json.Unmarshal(data, &request)
	if err != nil {
		logger.ZapLogger.Errorw(models.COMMUNITYCHEST, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Community Chest Validation")
	return
}

func (p *CommunityChestReq) ProcessMsg(body any) (resp []byte, err error) {
	logger.ZapLogger.Infoln("Enter Community Chest Processor")

	req := body.(models.ReqBlockInfo)
	var infoList []*models.BlockInfo

	communityChestInfo, err := p.db.GetBlockInfoByBlockType(models.COMMUNITYCHEST)
	if err != nil {
		logger.ZapLogger.Errorw(models.COMMUNITYCHEST, "DB Error", err)
		return
	}

	for _, info := range communityChestInfo {
		blockInfo := models.BlockInfo{
			InfoId: info.InfoId,
			Info: info.Info,
		}

		infoList = append(infoList, &blockInfo)
	}

	communityChestResp := models.RespBlockInfo{
		MsgId: req.MsgId,
		BlockType: models.COMMUNITYCHEST,
		BlockInfo: infoList,
	}

	resp, err = json.Marshal(communityChestResp)
	if err != nil {
		logger.ZapLogger.Errorw(models.COMMUNITYCHEST, "JSON Error", err)
		return
	}

	logger.ZapLogger.Infoln("Enter Community Chest Processor")
	return
}