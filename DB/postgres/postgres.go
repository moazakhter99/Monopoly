package postgres

import (
	models "Monopoly/Models"
	"Monopoly/logger"
	"database/sql"
	// "fmt"
	// "github.com/spf13/viper"
)

type Postgres struct {
	DB *sql.DB
}

// GetBlockInfo implements [db.DbOperations].
func (p *Postgres) GetBlockInfo(blockType string, infoNo int) (cardInfo string, err error) {
	panic("unimplemented")
}

// GetBlockPrice implements [db.DbOperations].
func (p *Postgres) GetBlockPrice(blockId string) (price int, err error) {
	panic("unimplemented")
}

// GetPlayerCashPos implements [db.DbOperations].
func (p *Postgres) GetPlayerCashPos(playerId string) (cash int, pos int, err error) {
	panic("unimplemented")
}

// GetPlayerStatus implements [db.DbOperations].
func (p *Postgres) GetPlayerStatus(playerId string, gameId string) (status string, err error) {
	panic("unimplemented")
}

// GetPlayerStatusList implements [db.DbOperations].
func (p *Postgres) GetPlayerStatusList(playerId string, gameId string) (statusList []string, err error) {
	panic("unimplemented")
}

// GetCardOwnerCount implements [db.DbOperations].
func (p *Postgres) GetCardOwnerCount(playerId string, gameId string) (count int, err error) {
	panic("unimplemented")
}

// GetCardOwnershipStatus implements [db.DbOperations].
func (p *Postgres) GetCardOwnershipStatus(playerId string, blockId string) (status string, err error) {
	panic("unimplemented")
}

// GetBlockOwner implements [db.DbOperations].
func (p *Postgres) GetBlockOwner(blockID string, gameID string) (playerID string, err error) {
	panic("unimplemented")
}

// GetPlayerSeqAndCount implements [db.DbOperations].
func (p *Postgres) GetPlayerSeqAndCount(playerId string) (seq int, count int, err error) {
	panic("unimplemented")
}

// GetNextPlayer implements [db.DbOperations].
func (p *Postgres) GetNextPlayer(gameId string, seq int) (playerId string, err error) {
	panic("unimplemented")
}

// DeleteGetOutOfJailCard implements [db.DbOperations].
func (p *Postgres) DeleteGetOutOfJailCard(playerId string, GameId, cardId string) (err error) {
	panic("unimplemented")
}

// UpdatePlayerStatus implements [db.DbOperations].
func (p *Postgres) UpdatePlayerStatus(playerId string, count string) (err error) {
	panic("unimplemented")
}

// GetCardOwnership implements [db.DbOperations].
func (p *Postgres) GetCardOwnership(blockId string, gameId string) (playerId string, err error) {
	panic("unimplemented")
}

// UpdateGetOutOfJailCard implements [db.DbOperations].
func (p *Postgres) UpdateGetOutOfJailCard(playerId string, gameId string) (err error) {
	panic("unimplemented")
}

// GetBlockInfoByBlockType implements [db.DbOperations].
func (p *Postgres) GetBlockInfoByBlockType(blockType string) (jailInfo []models.Jail, err error) {
	panic("unimplemented")
}

// GetPosByBlockName implements [db.DbOperations].
func (p *Postgres) GetPosByBlockName(blockName string) (pos int, err error) {
	panic("unimplemented")
}

// GetCardAction implements [db.DbOperations].
func (p *Postgres) GetCardAction(cardNo string) (action string, err error) {
	panic("unimplemented")
}

// GetPlayerCash implements [db.DbOperations].
func (p *Postgres) GetPlayerCash(playerId string) (cash int, err error) {
	panic("unimplemented")
}

// UpdatePlayerCash implements [db.DbOperations].
func (p *Postgres) UpdatePlayerCash(playerId string, cash, pos int) (err error) {
	panic("unimplemented")
}

// GetBlockInfoById implements [db.DbOperations].
func (p *Postgres) GetBlockInfoById(blockId string) (block *models.Block, err error) {
	block = &models.Block{
		BlockId: "block11",
		Price:   700,
	}
	return
}

// UpdatePlayerCard implements [db.DbOperations].
func (p *Postgres) UpdatePlayerCard(playerId string, gameId string, blockId, status string) (err error) {
	return
}

// UpdatePlayerPos implements [db.DbOperations].
func (p *Postgres) UpdatePlayerPos(playerId string, position int) (err error) {
	return
}

// GetBlockState implements [db.DbOperations].
func (p *Postgres) GetBlockState(position int, gameId string) (block *models.Block, er error) {
	return
}

// GetPlayerInfoById implements [db.DbOperations].
func (p *Postgres) GetPlayerInfoById(playerId string) (player *models.Player, err error) {
	player = &models.Player{
		PlayerId: playerId,
		Name:     "Moaz_" + playerId,
		Pos:      5,
		GameId:   "ABC",
		Cash:     1500,
		Seq:      1,
	}
	return
}

func OpenDatabase() (*Postgres, error) {

	logger.ZapLogger.Infoln("Inside Open DB connection")

	// host := viper.GetString("DB.DB_HOSTNAME")
	// port := viper.GetString("DB.DB_PORT")
	// user := viper.GetString("DB.DB_USER")
	// password := viper.GetString("DB.DB_PASSWORD")
	// dbname := viper.GetString("DB.DB_NAME")
	// psqlSSLMode := viper.GetString("DB.PSQL_SSL_MODE")

	// psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, psqlSSLMode)
	// db, err := sql.Open("postgres", psqlconn)
	// if err != nil {
	// 	logger.ZapLogger.Fatalw("Database Open", "Err", err)
	// 	return nil, err
	// }

	// Testing
	db := &sql.DB{}
	data := &Postgres{DB: db}

	return data, nil
}

func (p *Postgres) Ping() (err error) {
	logger.ZapLogger.Infow("Enter Postgres Ping")

	err = p.DB.Ping()
	if err != nil {
		logger.ZapLogger.Fatalw("Connecting to the Database", "Err", err)
		return err
	}
	logger.ZapLogger.Infoln("Database Connection Successfully Done")

	logger.ZapLogger.Infoln("Exit Postgres Ping")
	return
}

func (p *Postgres) InsertGame(gameId, matchId string) (err error) {
	logger.ZapLogger.Infoln("Enter Insert Game")

	logger.ZapLogger.Infoln("Exit Insert Game")
	return
}

func (p *Postgres) InsertPlayer(player *models.Player, gameId string) (err error) {
	logger.ZapLogger.Infoln("Enter InsertPlayer DB")

	logger.ZapLogger.Infoln("Exit InsertPlayer DB")
	return
}

func (p *Postgres) GetGameFromMatchId(matchId string) (gameId string, err error) {
	logger.ZapLogger.Infoln("Enter GetgameFromMatchId DB")

	logger.ZapLogger.Infoln("Exit GetgameFromMatchId DB")
	return
}
