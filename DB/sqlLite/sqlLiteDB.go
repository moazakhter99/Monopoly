package sqlLite

import (
	"Monopoly/logger"
	"database/sql"
	"os"

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
