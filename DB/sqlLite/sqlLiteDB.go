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