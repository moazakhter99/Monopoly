package gameroom

import (
	db "Monopoly/DB"
	"sync"
)

type Room interface {
	AddPlayer(gameId string, client *NewClient)
	RemovePlayer(playerId, gameId string)
	UpdateGameState(gameId, playerId, state string)
	Flush()
}

type ClientMap map[string]*NewClient

type GameRoom struct {
	// Remove whatever not needed later
	GameMap       map[string]bool // Game Map
	GameState     map[string][2]string
	Games         []string        // Game List
	PlayerMap     map[string]bool // PlayerId Map
	ClientMap     ClientMap       // Client Map
	GamePlayerMap map[string]ClientMap
	db            db.DbOperations
	mut           sync.Mutex
}

func CreateGameRoom(db db.DbOperations) *GameRoom {
	return &GameRoom{
		GameMap:       make(map[string]bool),
		GameState:     make(map[string][2]string),
		PlayerMap:     make(map[string]bool),
		ClientMap:     make(map[string]*NewClient),
		GamePlayerMap: make(map[string]ClientMap, 6),
		db: db,
	}
}

// Add New Player
func (r *GameRoom) AddPlayer(gameId string, client *NewClient) {
	r.mut.Lock()
	defer r.mut.Unlock()

	if !r.GameMap[gameId] {
		r.GameMap[gameId] = true
		r.Games = append(r.Games, gameId)
	}

	if !r.PlayerMap[client.PlayerId] {
		r.PlayerMap[client.PlayerId] = true
		r.ClientMap[client.PlayerId] = client
		r.GamePlayerMap[gameId][client.PlayerId] = client
		r.GameState[gameId] = [2]string{
			client.PlayerId, "Added",
		}
	}
}

// Remove Player
func (r *GameRoom) RemovePlayer(playerId, gameId string) {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.Flush()
	if r.PlayerMap[playerId] {
		delete(r.ClientMap, playerId)
		delete(r.GamePlayerMap[gameId], playerId)
		delete(r.PlayerMap, playerId)
	}

	if len(r.GamePlayerMap) == 0 {
		delete(r.GameMap, gameId)
		delete(r.GameState, gameId)
		delete(r.GamePlayerMap, gameId)
	}

}

func (r *GameRoom) UpdateGameState(gameId, playerId, state string) {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.GameState[gameId] = [2]string{
		playerId, state,
	}
}

func (r *GameRoom) Flush() {
	// Write all the info to db
}

// func removeFromList(list []string, value string) []string {
// 	for i, v := range list {
// 		if v == value {
// 			return append(list[:i], list[i+1:]...)
// 		}
// 	}
// 	return list
// }
// // Create New Room
// func (r *GameRoom) CreateNewGame(gameId string) {
// 	r.mut.Lock()
// 	defer r.mut.Unlock()
// 	if r.GameMap[""] {}
// 	if !r.GameExists(gameId) {
// 		r.GameMap[gameId] = true
// 		r.Games = append(r.Games, gameId)
// 		// r.GamePlayerMap[gameId] = make(map[string]*NewClient)

// 	}

// }

// // Check existing room
// func (r *GameRoom) GameExists(gameId string) bool {
// 	_, ok := r.GameMap[gameId]
// 	return ok
// }

// // Get Room

// // Remove Room
// func (r *GameRoom) RemoveRoom(gameId string) {
// 	r.mut.Lock()
// 	defer r.mut.Unlock()
// 	delete(r.GameMap, gameId)
// 	delete(r.GamePlayerMap, gameId)
// 	r.Games = removeFromList(r.Games, gameId)

// }

// // Check Existing Player
// func (r *GameRoom) PlayerExists(playerId string) bool {
// 	_, ok := r.PlayerMap[playerId]
// 	return ok
// }

// // Get Player List
// func (r *GameRoom) PlayerList(gameId string) (playerList []string) {
// 	return //r.GamePlayerMap[gameId]
// }
