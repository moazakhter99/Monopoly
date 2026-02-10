package db


type DbOperations interface {
	Ping() (err error)	
	InsertGame(gameId, matchId string) (err error)
	
}