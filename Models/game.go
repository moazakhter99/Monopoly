package models



type Game struct {
	GameId 		string 		`json:"gameId"`
	PlayerCount int 		`json:"playerCount"`
	PlayerList 	[]*Player 	`json:"playerList"`
	LastPlayer 	string 		`json:"lastPlayer"`
	NextPlayer 	string 		`json:"nextPlayer"`

}