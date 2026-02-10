package postgres

import (
	"Monopoly/logger"
	"database/sql"
	"fmt"

	"github.com/spf13/viper"
)

type Postgres struct {
	DB *sql.DB
}

func OpenDatabase() (*Postgres, error) {

	logger.ZapLogger.Infoln("Inside Open DB connection")

	host := viper.GetString("DB.DB_HOSTNAME")
	port := viper.GetString("DB.DB_PORT")
	user := viper.GetString("DB.DB_USER")
	password := viper.GetString("DB.DB_PASSWORD")
	dbname := viper.GetString("DB.DB_NAME")
	psqlSSLMode := viper.GetString("DB.PSQL_SSL_MODE")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, psqlSSLMode)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		logger.ZapLogger.Fatalw("Database Open", "Err", err)
		return nil, err
	}

	data := &Postgres{DB: db}

	return data, err
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
