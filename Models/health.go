package models


type Health struct {
	Msg 			string `json:"Msg"`
	DBStatus 		string `json:"DbStatus"`
	ServiceStatus 	string `json:"ServiceStatus"`
}