package sqlLite

import (
	models "Monopoly/Models"
	"Monopoly/logger"
	"context"
	"database/sql"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type SqlLite struct {
	DB *sql.DB
}

func OpenDatabase() (db *SqlLite, err error) {
	logger.ZapLogger.Infoln("Creating Database")
	DB, err := sql.Open("sqlite", "./DB/sqlLite/sqlLiteDB/monopoly.db")
	if err != nil {
		logger.ZapLogger.Fatalw("Sqlite DB Open", "Error ", err)
	}

	db = &SqlLite{
		DB: DB,
	}

	sqlPath := "./monopolyDB.sql"
	sqlBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		logger.ZapLogger.Errorw("Create Sqlite DB", "Error ", err)
		return nil, err
	}

	_, err = db.DB.Exec(string(sqlBytes))
	if err != nil {
		logger.ZapLogger.Errorw("Create Tables Exec", "Error", err)
		return nil, err
	}

	logger.ZapLogger.Infow("Database Created")
	return
}

func (l *SqlLite) Ping() (err error) {
	logger.ZapLogger.Infow("Enter SqlLite Ping")

	err = l.DB.Ping()
	if err != nil {
		logger.ZapLogger.Errorw("Connecting to the Database", "Err", err)
		return err
	}
	logger.ZapLogger.Infoln("Database Connection Successfully Done")

	logger.ZapLogger.Infoln("Exit SqlLite Ping")
	return
}

func (l *SqlLite) InsertGame(gameId, matchId string) (err error) {
	logger.ZapLogger.Infoln("Enter InsertGame DB")

	query := `INSERT INTO game (game_id, match_id) VALUES (?, ?)`

	ctx, cancelF := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelF()
	txn, err := l.DB.BeginTx(ctx, nil)
	if err != nil {
		logger.ZapLogger.Errorw("Begin Transaction", "Error", err)
		txn.Rollback()
		return
	}

	_, err = txn.Exec(query, gameId, matchId)
	if err != nil {
		logger.ZapLogger.Errorw("DB Insert", "Error", err)
		txn.Rollback()
		return
	}

	err = txn.Commit()
	if err != nil {
		logger.ZapLogger.Errorw("DB Commit", "Error", err)
		txn.Rollback()
		return
	}

	logger.ZapLogger.Infoln("Exit InsertGame DB")
	return
}

func (l *SqlLite) InsertPlayer(player *models.Player, gameId string) (err error) {
	logger.ZapLogger.Infoln("Enter InsertPlayer DB")

	query := `INSERT INTO player (player_id, player_name, position, gameId, cash) VALUES (?, ?, ?, ?, ?)`

	ctx, cancelF := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelF()
	txn, err := l.DB.BeginTx(ctx, nil)
	if err != nil {
		logger.ZapLogger.Errorw("Begin Transaction", "Error", err)
		txn.Rollback()
		return
	}

	_, err = txn.Exec(query, player.PlayerId, player.Name, player.Pos, gameId, player.Cash)
	if err != nil {
		logger.ZapLogger.Errorw("DB Insert", "Error", err)
		txn.Rollback()
		return
	}

	err = txn.Commit()
	if err != nil {
		logger.ZapLogger.Errorw("DB Commit", "Error", err)
		txn.Rollback()
		return
	}

	logger.ZapLogger.Infoln("Exit InsertPlayer DB")
	return
}

func (l *SqlLite) GetGameFromMatchId(matchId string) (gameId string, err error) {
	logger.ZapLogger.Infoln("Enter GetgameFromMatchId DB")

	query := `SELECT game_id FROM game WHERE match_id = ?`

	row := l.DB.QueryRow(query, matchId)

	var game_id sql.NullString

	err = row.Scan(&game_id)
	if err != nil {
		logger.ZapLogger.Errorw("DB Select", "Error", err)
		return
	}
	gameId = game_id.String

	logger.ZapLogger.Infoln("Exit GetgameFromMatchId DB")
	return
}

func (l *SqlLite) GetPlayerInfoById(playerId string) (player *models.Player, err error) {
	logger.ZapLogger.Infoln("Enter Get Player Info by Id DB")

	query := `SELECT player_id, player_name, position, game_id, cash, seq FROM player WHERE player_id = ?`

	row := l.DB.QueryRow(query, playerId)

	var (
		player_id   sql.NullString
		player_name sql.NullString
		position    sql.NullInt16
		game_id     sql.NullString
		cash        sql.NullInt64
		seq         sql.NullInt16
	)

	err = row.Scan(&player_id, &player_name, &position, &game_id, &cash, &seq)
	if err != nil {
		return
	}

	player = &models.Player{
		PlayerId: player_id.String,
		Name:     player_name.String,
		Pos:      int(position.Int16),
		GameId:   game_id.String,
		Cash:     int(cash.Int64),
		Seq:      int(seq.Int16),
	}

	logger.ZapLogger.Infoln("Exit Get Player Info by Id DB")
	return

}

func (l *SqlLite) UpdatePlayerPos(playerId string, position int) (err error) {
	logger.ZapLogger.Infoln("Enter Update Player Position")

	query := `UPDATE player SET position = ? WHERE player_id = ?`

	_, err = l.DB.Exec(query, position, playerId)
	if err != nil {
		return
	}

	logger.ZapLogger.Infoln("Exit Update Player Position")
	return
}

func (l *SqlLite) GetBlockState(position int, gameId string) (block *models.Block, er error) {

	var state bool
	var err error
	query := `SELECT block_id, block_tye from game_board where position = ?`

	row := l.DB.QueryRow(query, position)

	var (
		block_id   sql.NullString
		block_type sql.NullString
	)

	err = row.Scan(&block_id, &block_type)
	if err != nil {
		return
	}

	blockId := block_id.String

	query2 := `SELECT player_id FROM game_player WHERE game_id = ? and card_id = ?`

	row2 := l.DB.QueryRow(query2, gameId, blockId)

	var player_id sql.NullString

	err = row2.Scan(&player_id)
	if err == sql.ErrNoRows {
		state = false
	} else if err != nil {
		er = err
		return nil, er
	}
	// playerId := player_id.String
	state = true

	block = &models.Block{
		BlockId: blockId,
		Type:    block_type.String,
		State:   state,
	}

	return
}

func (l *SqlLite) GetBlockInfoById(blockId string) (block *models.Block, err error) {
	logger.ZapLogger.Infoln("Enter Get Block Info by Id")

	query := `SELECT block_id, block_type, block_name, colour, position, price, base_rent FROM game_board WHERE block_id = ?`

	row := l.DB.QueryRow(query, blockId)

	var (
		block_id     sql.NullString
		block_type   sql.NullString
		block_name   sql.NullString
		block_colour sql.NullString
		position     sql.NullInt16
		block_price  sql.NullInt64
		base_rent    sql.NullInt64
	)

	err = row.Scan(&block_id, &block_type, &block_name, &block_colour, &position, &block_price, &base_rent)
	if err != nil {
		return
	}

	block = &models.Block{
		BlockId:   block_id.String,
		Type:      block_type.String,
		BlockName: block_name.String,
		Colour:    block_colour.String,
		Position:  int(position.Int16),
		Price:     int(block_price.Int64),
		BaseRent:  int(base_rent.Int64),
	}

	logger.ZapLogger.Infoln("Exit Get Block Info by Id")
	return
}

func (l *SqlLite) UpdatePlayerCard(playerId, gameId, blockId string) (err error) {
	logger.ZapLogger.Infoln("Enter Update Player Card")

	query := `INSERT INTO game_player (player_id, game_id, card_id) VALUES (?, ?, ?)`

	_, err = l.DB.Exec(query, playerId, gameId, blockId)
	if err != nil {
		return
	}

	return
}

func (l *SqlLite) GetCardAction(cardNo string) (action string, err error) {
	logger.ZapLogger.Infoln("Enter Get Card Action")

	query := `SELECT action FROM block_info WHERE info_id = ?`

	row := l.DB.QueryRow(query, cardNo)

	var card_action sql.NullString
	err = row.Scan(&card_action)
	if err != nil {
		return
	}

	action = card_action.String

	return
}

func (l *SqlLite) GetPlayerCash(playerId string) (cash int, err error) {
	logger.ZapLogger.Infoln("Enter Get Player Cash")

	query := `SELECT cash FROM player WHERE player_id = ?`

	row := l.DB.QueryRow(query, playerId)

	var player_cash sql.NullInt64
	err = row.Scan(&player_cash)
	if err != nil {
		return
	}

	cash = int(player_cash.Int64)

	return
}

func (l *SqlLite) UpdatePlayerCash(playerId string, cash int) (err error) {
	logger.ZapLogger.Infoln("Update Player Cash")

	query := `UPDATE player SET cash = ? WHERE player_id = ?`

	_, err = l.DB.Exec(query, playerId)
	if err != nil {
		return
	}

	return
}

func (l *SqlLite) GetPosByBlockName(blockName string) (pos int, err error) {
	logger.ZapLogger.Infoln("Enter Get Position By Block Name")

	query := `SELECT position FROM game_board WHERE block_name = ?`

	row := l.DB.QueryRow(query, blockName)

	var position sql.NullInt16
	err = row.Scan(&position)
	if err != nil {
		return
	}

	pos = int(position.Int16)

	return
}

func (l *SqlLite) GetBlockInfoByBlockType(blockType string) (jailInfo []models.Jail, err error) {
	logger.ZapLogger.Infow("Enter Get Block Info by Block Type")

	query := `SELECT Info_id, block_info FROM block_info WHERE block_type = ?`

	rows, err := l.DB.Query(query, blockType)
	if err != nil {
		return
	}

	for rows.Next() {

		var (
			info_id sql.NullString
			block_info sql.NullString
		)

		err = rows.Scan(&info_id, &block_info)
		if err != nil {
			return
		}

		jail := models.Jail{
			InfoId: info_id.String,
			Info: block_info.String,
		}

		jailInfo = append(jailInfo, jail)

	}

	return
}

func (l *SqlLite) GetCardOwnership(blockId, gameId string) (playerId string, err error) {
	logger.ZapLogger.Infoln("Enter Get Card Ownership")

	query := `SELECT playerId FROM game_player WHERE block_id = ? AND game_id = ?`

	row := l.DB.QueryRow(query, blockId, gameId)

	var player_id sql.NullString

	err = row.Scan(&player_id)

	playerId = player_id.String

	return
}

func (l *SqlLite) UpdateGetOutOfJailCard(playerId, gameId string) (err error) {
	logger.ZapLogger.Infoln ("Enter Update Get Out Of Jail Card")

	query := `INSERT INTO game_player (player_id, game_id, card_id) VALUES (?, ?, (SELECT info_id FROM block_info WHERE block_type = 'SPECIAL_CARD'))`

	_, err = l.DB.Exec(query, playerId, gameId)
	if err != nil {
		return
	}

	return
}

func (l *SqlLite) DeleteGetOutOfJailCard(playerId, GameId string) (err error) {

	return
}

func (l *SqlLite) UpdatePlayerStatus(playerId, count string) (err error) {
	logger.ZapLogger.Infoln("Enter Update Player Status")

	status := models.BLOCKED + "_" + count
	query := `UPDATE player SET status = ? WHERE player_id = ?`

	_, err = l.DB.Exec(query, status, playerId)
	if err != nil {
		return
	}

	return
}