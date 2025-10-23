package manager

import (
	"fmt"
	"sync"

	"github.com/luispellizzon/pangram/internal/games"
)

// GameManager interface
type Manager interface {
	Create(kind string) (string, games.Game, error)
	Get(id string) (games.Game, bool)
}

// mgr is the server manager that will save all games and will act as a singleton and also a proxy, since all requests will call this manager to grab a game by id from its inGames mapper or to create new games
type mgr struct {
	mu     sync.RWMutex
	nextID int
	inGames  map[string]games.Game
	factory *games.Factory
}

func New(factory *games.Factory) Manager {
	return &mgr{inGames: map[string]games.Game{}, factory: factory}
}

// Create new game using the server factory
func (m *mgr) Create(kind string) (string, games.Game, error) {
	game, err := m.factory.New(kind)
	if err != nil { 
		return "", nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	id := fmt.Sprintf("g-%d", m.nextID)
	m.inGames[id] = game
	return id, game, nil
}

// Get game by id, the game is saved inside the manager singleton inGames map
func (m *mgr) Get(id string) (games.Game, bool) {
	m.mu.RLock(); defer m.mu.RUnlock()
	g, ok := m.inGames[id]
	return g, ok
}
