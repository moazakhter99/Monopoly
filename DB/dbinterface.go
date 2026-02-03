package db


type DbOperations interface {
	Ping() (err error)	
}