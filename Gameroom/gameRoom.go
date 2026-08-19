package gameroom

import (
	db "Monopoly/DB"
	"Monopoly/logger"
	"sync"
)

type Room interface {
	AddPlayer(gameId string, client NewClient)
	RemovePlayer(playerId, gameId string)
	UpdateGameState(gameId, playerId, state string)
	Flush()
	flush()
	GetClientListByGame(gameId string) (clientMap ClientMap)
}

type ClientMap map[string]*NewClient

type GameRoom struct {
	// Remove whatever not needed later
	GameMap       map[string]bool // Game Map
	GameState     map[string][2]string
	Games         []string        // Game List
	PlayerMap     map[string]bool // PlayerId Map
	GamePlayerMap map[string]ClientMap
	db            db.DbOperations
	mut           sync.Mutex
}

func CreateGameRoom(db db.DbOperations) *GameRoom {
	return &GameRoom{
		GameMap:       make(map[string]bool),
		GameState:     make(map[string][2]string),
		PlayerMap:     make(map[string]bool),
		GamePlayerMap: make(map[string]ClientMap, 6),
		db: db,
	}
}

// Add New Player
func (r *GameRoom) AddPlayer(gameId string, client NewClient) {
	r.mut.Lock()
	defer r.mut.Unlock()

	if !r.GameMap[gameId] {
		r.GameMap[gameId] = true
		r.Games = append(r.Games, gameId)
	}

	if !r.PlayerMap[client.PlayerId] {
		r.PlayerMap[client.PlayerId] = true
		if r.GamePlayerMap[gameId] == nil {
			r.GamePlayerMap[gameId] = make(ClientMap)
		}
		r.GamePlayerMap[gameId][client.PlayerId] = &client
		r.GameState[gameId] = [2]string{
			client.PlayerId, "Added",
		}
		logger.ZapLogger.Infow("State Update", "GameId", gameId, "PlayerId", client.PlayerId, "Updated State", "Added")
	}
}

// Remove Player
func (r *GameRoom) RemovePlayer(playerId, gameId string) {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.flush()
	if r.PlayerMap[playerId] {
		delete(r.GamePlayerMap[gameId], playerId)
		delete(r.PlayerMap, playerId)
		logger.ZapLogger.Infow("Remove Player", "PlayerId", playerId, "Player left", len(r.GamePlayerMap[gameId]))
	}

	if len(r.GamePlayerMap[gameId]) == 0 {
		delete(r.GameMap, gameId)
		delete(r.GameState, gameId)
		delete(r.GamePlayerMap, gameId)
		logger.ZapLogger.Infow("Drop Game", "GameId", gameId)
	}

}

func (r *GameRoom) UpdateGameState(gameId, playerId, state string) {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.GameState[gameId] = [2]string{
		playerId, state,
	}
	logger.ZapLogger.Infow("State Update", "GameId", gameId, "PlayerId", playerId, "Updated State", state)
}

func (r *GameRoom) Flush() {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.flush()
}

func (r *GameRoom) flush() {
	// Write all the info to db
}

func (r *GameRoom) GetClientListByGame(gameId string) (clientMap ClientMap) {
	return r.GamePlayerMap[gameId]
}
