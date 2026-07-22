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
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
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
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
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
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Update Player Position")
	return
}

func (l *SqlLite) GetBlockState(position int, gameId string) (block *models.Block, err error) {

	query := `SELECT block_id, block_type, block_name, price from game_board where position = ?`

	row := l.DB.QueryRow(query, position)

	var (
		block_id   sql.NullString
		block_type sql.NullString
		block_name sql.NullString
		block_price sql.NullInt64
	)

	err = row.Scan(&block_id, &block_type, &block_name, &block_price)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	block = &models.Block{
		BlockId: block_id.String,
		Type:    block_type.String,
		BlockName: block_name.String,
		Price: int(block_price.Int64),
	}

	if block.Type != models.SPECIALCARD {
		playerId, _ := l.GetBlockOwner(block.BlockId, gameId)
		if playerId != "" {
			block.OwnerId = playerId
		}
	}

	return
}

func (l *SqlLite) GetBlockOwner(blockID, gameID string) (playerId string, err error) {

	query := `SELECT player_id FROM game_player WHERE game_id = ? and card_id = ?`

	row := l.DB.QueryRow(query, gameID, blockID)

	var player_id sql.NullString

	err = row.Scan(&player_id)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return "", nil
	}

	playerId = player_id.String

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
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
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

func (l *SqlLite) UpdateCardColour(playerId, colour string) (err error) {
	logger.ZapLogger.Infoln("Enter Update Card Colour")

	query := `UPDATE game_player SET = ? WHERE playerId = ? AND card_id IN (SELECT block_id FROM game_board WHERE colour = ?)`

	_, err = l.DB.Exec(query, models.COLOUR, playerId, colour)

	logger.ZapLogger.Infoln("Exit Uodate Card Count")
	return
}

func (l *SqlLite) UpdatePlayerCard(playerId, gameId, blockId string) (err error) {
	logger.ZapLogger.Infoln("Enter Update Player Card")

	query := `INSERT INTO game_player (player_id, game_id, card_id, status) VALUES (?, ?, ?, ?)`

	_, err = l.DB.Exec(query, playerId, gameId, blockId, models.BASE)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
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
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	action = card_action.String

	return
}

func (l *SqlLite) GetPlayerCashPos(playerId string) (cash, pos int, err error) {
	logger.ZapLogger.Infoln("Enter Get Player Cash")

	query := `SELECT cash, position FROM player WHERE player_id = ?`

	row := l.DB.QueryRow(query, playerId)

	var player_cash sql.NullInt64
	var player_pos sql.NullInt16
	err = row.Scan(&player_cash, &player_pos)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	cash = int(player_cash.Int64)
	pos = int(player_pos.Int16)

	return
}

func (l *SqlLite) UpdatePlayerCash(playerId string, cash, pos int) (err error) {
	logger.ZapLogger.Infoln("Update Player Cash")

	query := `UPDATE player SET cash = ?, position = ? WHERE player_id = ?`

	_, err = l.DB.Exec(query, cash, pos, playerId)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
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
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	pos = int(position.Int16)

	return
}

func (l *SqlLite) GetBlockInfoByBlockType(blockType string) (jailInfo []models.Jail, err error) {
	logger.ZapLogger.Infoln("Enter Get Block Info by Block Type")

	query := `SELECT Info_id, block_info FROM block_info WHERE block_type = ?`

	rows, err := l.DB.Query(query, blockType)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	for rows.Next() {

		var (
			info_id sql.NullString
			block_info sql.NullString
		)

		err = rows.Scan(&info_id, &block_info)
		if err != nil {
			logger.ZapLogger.Errorw("SQL Err", "Error", err)
			return
		}

		jail := models.Jail{
			InfoId: info_id.String,
			Info: block_info.String,
		}

		jailInfo = append(jailInfo, jail)

	}

	logger.ZapLogger.Infoln("Exit Get Block Info by Block Type")
	return
}

func (l *SqlLite) GetCardOwnership(blockId, gameId string) (playerId string, err error) {
	logger.ZapLogger.Infoln("Enter Get Card Ownership")

	query := `SELECT playerId FROM game_player WHERE block_id = ? AND game_id = ?`

	row := l.DB.QueryRow(query, blockId, gameId)

	var player_id sql.NullString

	err = row.Scan(&player_id)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	playerId = player_id.String

	return
}

func (l *SqlLite) UpdateGetOutOfJailCard(playerId, gameId string) (err error) {
	logger.ZapLogger.Infoln ("Enter Update Get Out Of Jail Card")

	query := `INSERT INTO game_player (player_id, game_id, card_id) VALUES (?, ?, (SELECT info_id FROM block_info WHERE block_type = 'SPECIAL_CARD'))`

	_, err = l.DB.Exec(query, playerId, gameId)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	return
}

func (l *SqlLite) DeleteGetOutOfJailCard(playerId, gameId, cardId string) (err error) {
	logger.ZapLogger.Infoln("Enter Delete Get Out of Jail")

	query := `DELETE FROM game_player WHERE player_id = ? AND game_id = ? AND card_id = ?`

	_, err = l.DB.Exec(query, playerId, gameId, cardId)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	logger.ZapLogger.Infoln("Exit Delete Get Out of Jail")
	return
}

func (l *SqlLite) UpdatePlayerStatus(playerId, status string) (err error) {
	logger.ZapLogger.Infoln("Enter Update Player Status")

	query := `UPDATE player SET status = ? WHERE player_id = ?`

	_, err = l.DB.Exec(query, status, playerId)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	return
}

func (l *SqlLite) GetNextPlayer(gameId string, seq int) (playerId string, err error) {
	logger.ZapLogger.Infoln("Enter Get Next Player")

	query := `SELECT player_id FROM player WHERE game_id = ? AND seq = ?`

	row := l.DB.QueryRow(query, gameId, seq)

	var player_id sql.NullString

	err = row.Scan(&player_id)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	playerId = player_id.String

	return
}

func (l *SqlLite) GetPlayerSeqAndCount(playerId string) (seq, count int, err error) {
	logger.ZapLogger.Infoln("Enter Get Player Seq and Count")

	query := `SELECT p.seq, g.player_count FROM game as g JOIN player as p ON g.game_id = p.game_id WHERE p.player_id = ?`

	row := l.DB.QueryRow(query, playerId)

	var (
		player_seq sql.NullInt16
		game_count sql.NullInt16
	)

	err = row.Scan(&player_seq, &game_count)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	seq = int(player_seq.Int16)
	count = int(game_count.Int16)

	return
}

func (l *SqlLite) GetCardOwnershipStatus(playerId, blockId string) (status string, err error) {
	logger.ZapLogger.Infoln("Enter Get Card Ownership Status")

	query := `SELECT status FROM game_player WHERE player_id = ? AND card_id = ?`

	row := l.DB.QueryRow(query, playerId, blockId)

	var player_status sql.NullString

	err = row.Scan(&player_status)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	status = player_status.String

	return
}


func (l *SqlLite) GetCardOwnerCount(playerId, gameId string) (count int, err error) {
	logger.ZapLogger.Infoln("Enter Get Card Owner Count")

	query := `SELECT count(*) 
				FROM game_player AS gp 
				JOIN game_board AS gb ON gp.card_id = gb.block_id 
				WHERE gb.block_type = ? AND gp.player_id = ? AND game_id = ?;`

	row := l.DB.QueryRow(query, models.UTILITY, playerId, gameId)

	var card_count sql.NullInt16
	err = row.Scan(&card_count)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}
	
	count = int(card_count.Int16)

	return
}


func (l *SqlLite) GetPlayerStatusList(playerId, gameId string) (statusList []string, err error) {
	logger.ZapLogger.Infoln("Enter Get Player Status List")

	query := `SELECT status FROM game_player WHERE player_id = ? AND game_id = ?`

	rows, err := l.DB.Query(query, playerId, gameId)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	for rows.Next() {
		var card_status sql.NullString

		err = rows.Scan(&card_status)
		if err != nil {
			logger.ZapLogger.Errorw("SQL Err", "Error", err)
			return
		}
		statusList = append(statusList, card_status.String)

	}

	logger.ZapLogger.Infoln("Exit Get Player Status List")
	return
}

func (l *SqlLite) GetPlayerStatus(playerId, gameId string) (status string, err error) {
	logger.ZapLogger.Infoln("Enter Get Player Status")

	query := `SELECT status FROM player WHERE player_id = ? and game_id = ?`

	row := l.DB.QueryRow(query, playerId, gameId)

	var player_status sql.NullString

	err = row.Scan(&player_status)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	status = player_status.String

	return
}

func (l *SqlLite) GetBlockInfo(blockType string, infoNo int) (cardInfo string, err error) {
	logger.ZapLogger.Infoln("Enter Get Block Info")

	query := `SELECT block_info FROM block_info WHERE block_type = ? AND info_no = ?`

	row := l.DB.QueryRow(query, blockType, infoNo)

	var card_info sql.NullString
	err = row.Scan(&cardInfo)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	cardInfo = card_info.String

	return
}

func (l *SqlLite) GetBlockPrice(blockId string) (price int, err error) {
	logger.ZapLogger.Infoln("Enter Get Block Price")

	query := `SELECT price FROM game_board WHERE block_id = ?`

	row := l.DB.QueryRow(query, blockId)

	var block_price sql.NullInt64
	err = row.Scan(&block_price)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	price = int(block_price.Int64)
	
	return
}

func (l *SqlLite) GetPlayerColourCount(playerId, blockId string) (cardCount int, cardColour string, err error) {
	logger.ZapLogger.Infoln("Enter Get Player Colour Count")

	query := `SELECT count(*), gb.colour FROM game_board AS gb 
				JOIN game_player AS gp ON gb.block_id = gp.card_id 
				WHERE gp.player_id = ? AND gb.colour = (
				SELECT colour FROM  game_board WHERE block_id = ?) GROUP BY gb.colour`
	
	row := l.DB.QueryRow(query, playerId, blockId)

	var card_count sql.NullInt16
	var card_colour sql.NullString
	err = row.Scan(&card_count, &card_colour)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	cardCount = int(card_count.Int16)
	cardColour = card_colour.String

	return
}

func (l *SqlLite) UpdateCardStatus(playerId, gameId, blockId, status string) (err error) {
	logger.ZapLogger.Infoln("Enter Update Card Status")

	query := `UPDATE game_player SET status = ? WHERE player_id = ? AND game_id = ? AND block_id = ?`

	_, err = l.DB.Exec(query, status, playerId, gameId, blockId)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}


	return
}

func (l *SqlLite) GetCardStatus(playerId, gameId, blockId string) (status string, err error) {
	logger.ZapLogger.Infoln("Enter Get Card Status")

	query := `SELECT status FROM game_player WHERE player_id = ? AND game_id = ? AND blockId = ?`

	row := l.DB.QueryRow(query, playerId, gameId, blockId)

	var card_status sql.NullString
	err = row.Scan(&card_status)
	if err != nil {
		logger.ZapLogger.Errorw("SQL Err", "Error", err)
		return
	}

	return
}