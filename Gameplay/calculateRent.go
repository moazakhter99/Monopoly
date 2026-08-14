package gameplay

import (
	db "Monopoly/DB"
	models "Monopoly/Models"
	"Monopoly/logger"
	"encoding/json"
	"strconv"
)



type CalculateRentProc struct {
	db 		db.DbOperations
	client 	*models.Client
	writeCh chan<- models.WSMessage

}

func CreateCalculateRent(db db.DbOperations) *CalculateRentProc {
	return &CalculateRentProc{
		db: db,
	}
}

func (c *CalculateRentProc) Validate(reqMsg []byte) (payload any, err error) {
	logger.ZapLogger.Infoln("Enter Validate Calculate Rent")
	var req models.ReqCalculateRent
	err = json.Unmarshal(reqMsg, &req)
	if err != nil {
		logger.ZapLogger.Errorw(models.CALCULATERENT, "Validation Error", err)
		return
	}
	logger.ZapLogger.Infoln("Exit Validate Calculate Rent")
	return req, err
}

// Might not work as the resposens are targeted
func (c *CalculateRentProc) Play(payload any) ([]byte, error) {
	logger.ZapLogger.Infoln("Enter Play Calculate Rent")
	req := payload.(models.ReqCalculateRent)
	var rent int
	var err error
	targetMap := make(map[string][]byte, 2)

	playerId := c.client.PlayerId
	gameId := c.client.GameId
	baseRent := req.Price * 10 / 100

	logger.ZapLogger.Infoln(models.CALCULATERENT, "Calculating Rent")

	switch req.BlockType {

	case models.UTILITY:

		cardCount, err := c.db.GetCardOwnerCount(playerId, gameId)
		if err != nil {
			logger.ZapLogger.Errorw(models.CALCULATERENT, "DB Error", err)
			return nil, err
		}
		rent = utilityRentCalculation(baseRent, cardCount)

	case models.CITY:

		status, err := c.db.GetCardOwnershipStatus(req.OwnerId, req.BlockId)
		if err != nil {
			logger.ZapLogger.Errorw(models.CALCULATERENT, "DB ERROR", err)
			return nil, err
		}
		rent = cityRentCalculation(baseRent, status)

	default:

	}

	status := models.PAYRENT + "-" + strconv.Itoa(rent) + "-" + req.OwnerId
	err = c.db.UpdatePlayerStatus(playerId, status)
	if err != nil {
		logger.ZapLogger.Errorw(models.CALCULATERENT, "DB Error", err)
		return nil, err
	}

	playerResp := models.RespCalculateRent{
		BlockId: req.BlockId,
		BlockType: req.BlockType,
		OwnerId: req.OwnerId,
		RenterId: playerId,
		Rent: rent,
	}

	playerResponse, err := json.Marshal(playerResp)
	if err != nil {
		logger.ZapLogger.Infow(models.CALCULATERENT, "Resp for", playerId, "JSON err", err)
		return nil, err
	}

	logger.ZapLogger.Infow(models.CALCULATERENT, "Client", playerId, "Resp Body", string(playerResponse))

	targetMap[playerId] = playerResponse

	ownerResp := models.RespCalculateRent{
		BlockId: req.BlockId,
		BlockType: req.BlockType,
		RenterId: playerId,
		OwnerId: req.OwnerId,
		Rent: rent,
	}

	ownerResponse, err := json.Marshal(ownerResp)
	if err != nil {
		logger.ZapLogger.Infow(models.CALCULATERENT, "Resp for", req.OwnerId, "JSON err", err)
		return nil, err
	}
	logger.ZapLogger.Infow(models.CALCULATERENT, "Client", req.OwnerId, "Resp Body", string(ownerResponse))

	targetMap[req.OwnerId] = ownerResponse

	logger.ZapLogger.Infoln("Exit Play Calculate Rent")
	return playerResponse, nil
}