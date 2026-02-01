package sqlLite

import (
	"Monopoly/logger"
	"database/sql"
)



type SqlLite struct {
	DB *sql.DB
}

func OpenDatabase() (db *SqlLite, err error) {

	DB, err := sql.Open("sqlite", "./DB/sqlLite/sqlLiteDB/monopoly.db") 
	if err != nil {
		logger.ZapLogger.Fatalln("Sqlite DB Open", "Error ",err) 
	}

	db = &SqlLite{
		DB: DB,
	}

	return
}