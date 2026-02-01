package sqlite

import "database/sql"


type Sqlite struct {
	DB *sql.DB

}


func OpenDatabase() (*Sqlite, error) {


	return nil, nil
}