package manager

import (
	"fmt"
	"sync"

	"github.com/luispellizzon/pangram/internal/games"
)

type Manager interface {
	Create(kind string) (string, games.Game, error)
	Get(id string) (games.Game, bool)
}

type mgr struct {
	mu     sync.RWMutex
	nextID int
	inGames  map[string]games.Game
	factory *games.Factory
}

func New(factory *games.Factory) Manager {
	return &mgr{inGames: map[string]games.Game{}, factory: factory}
}

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

func (m *mgr) Get(id string) (games.Game, bool) {
	m.mu.RLock(); defer m.mu.RUnlock()
	g, ok := m.inGames[id]
	return g, ok
}
