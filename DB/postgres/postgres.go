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
		logger.ZapLogger.Fatalln("Database Open", "Err", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logger.ZapLogger.Fatalf("Connecting to the Database", "Err" , err)
		return nil, err
	} else {
		logger.ZapLogger.Infoln("Database Connection Successfully Done")
	}

	data := &Postgres{DB: db}

	return data, err
}