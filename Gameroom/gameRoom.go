package gameroom

import (
	"sync"
)


type Room interface {
	AddPlayer(gameId string, client *NewClient)
	RemovePlayer(playerId, gameId string)
}

type GameRoom struct {
	GameMap map[string]bool // Game Map
	Games []string	// Game List
	PlayerMap map[string]string // PlayerId Map
	ClientMap map[string]*NewClient // Client Map
	GamePlayerMap map[string][]string
	mut sync.Mutex
}

func CreateGameRoom() *GameRoom {
	return &GameRoom{
		GameMap: make(map[string]bool),
		PlayerMap: make(map[string]string),
		ClientMap: make(map[string]*NewClient),
		GamePlayerMap: make(map[string][]string),
	}
}

// Create New Room
func (r *GameRoom) CreateNewGame(gameId string) {
	r.mut.Lock()
	defer r.mut.Unlock()
	if !r.GameExists(gameId) {
		r.GameMap[gameId] = true
		r.Games = append(r.Games, gameId)
		r.GamePlayerMap[gameId] = make([]string, 6) 

	}

}

// Check existing room
func (r *GameRoom) GameExists(gameId string) bool {
	_, ok := r.GameMap[gameId] 
	return ok
}

// Get Room

// Remove Room
func (r *GameRoom) RemoveRoom(gameId string) {
	r.mut.Lock()
	defer r.mut.Unlock()
	delete(r.GameMap, gameId)
	delete(r.GamePlayerMap, gameId)
	r.Games = removeFromList(r.Games, gameId)
	
}

// Add New Player
func (r *GameRoom) AddPlayer(gameId string, client *NewClient) {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.PlayerMap[client.PlayerId] = gameId
	r.ClientMap[client.PlayerId] = client
	if playerList, ok := r.GamePlayerMap[gameId]; ok {
		r.GamePlayerMap[gameId] = append(playerList, client.PlayerId)
	}
}

// Check Existing Player
func (r *GameRoom) PlayerExists(playerId string) bool {
	_, ok := r.PlayerMap[playerId]
	return ok
}

// Get Player List
func (r *GameRoom) PlayerList(gameId string) (playerList []string) {
	return r.GamePlayerMap[gameId]
}

// Remove Player
func (r *GameRoom) RemovePlayer(playerId, gameId string) {
	r.mut.Lock()
	defer r.mut.Unlock()
	delete(r.PlayerMap, playerId)
	delete(r.ClientMap, playerId)
	playerList := r.GamePlayerMap[gameId]
	r.GamePlayerMap[gameId] = removeFromList(playerList, playerId)

}


func removeFromList(list []string, value string) []string {
	for i, v := range list {
		if v == value {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}