package postgres

import "database/sql"


type Posgres struct {
	DB *sql.DB
}


func OpenDatabase() (*Posgres, error) {


	return nil, nil
}