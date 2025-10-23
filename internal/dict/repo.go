package dict

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/luispellizzon/pangram/internal/logger"
)

// Game repository interface
type Repository interface { Has(word string) (bool, error) }

// Save data in memory
type jsonMap struct{ data map[string]struct{} }

// Load dictionary (.json)
func loadJSON(path string) (*jsonMap, error) {
	bytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil { return nil, err }
	var objJSON map[string]any
	if err := json.Unmarshal(bytes, &objJSON); err != nil { return nil, err }
	mapper := make(map[string]struct{}, len(objJSON))
	for key := range objJSON {
		key = strings.TrimSpace(key)
		mapper[strings.ToLower(key)] = struct{}{}
	}
	return &jsonMap{data: mapper}, nil
}

// Adapter 
type JSONAdapter struct{ inner *jsonMap }

// Create a new repository using a json adapter, where it will convert json file, into a map in memory
func NewJSONAdapter(path string) (*JSONAdapter, error) {
	data, err := loadJSON(path); if err != nil { return nil, err }
	return &JSONAdapter{inner: data}, nil
}

// Implement repository
func (a *JSONAdapter) Has(word string) (bool, error) {
	_, ok := a.inner.data[strings.ToLower(word)]
	return ok, nil
}

// Cache proxy where it will take as a dependency the repository, so we can intercept requests before forward them to other layers of the server to check the dictionary
type CacheProxy struct {
	repo Repository
	cache map[string]bool
	order []string
	capacity   int
}

// Create proxy
func NewCacheProxy(repo Repository, capacity int) *CacheProxy {
	return &CacheProxy{repo: repo, cache: map[string]bool{}, capacity: capacity}
}

// Implement the same Repository interface, but the function will cache words that were already submitted so users that submits words that were checked before, will hit the cache without the need of reading the full repository where the full dictionary is loaded
func (p *CacheProxy) Has(word string) (bool, error) {
	pangram := strings.ToLower(word)
	isValid, ok := p.cache[pangram]
	if ok { 
		logger.Log().Infof("FROM CACHE: %v", word)
		return isValid, nil 
	}
	isValid, err := p.repo.Has(pangram)
	logger.Log().Infof("FROM DATABASE (REPOSITORY): %v", word)
	if err != nil { 
		logger.Log().Errorf("FROM CACHE: %v", word)
		return false, err 
	}

	// add to cache and clean older values. I set the capacity of the cache to be 5 words just for test purposes so I could see older values being cleaned. 
	p.cache[pangram] = isValid
	p.order = append(p.order, pangram)

	// Since we are using a queue way of caching, we remove the first value that from cache.
	if len(p.order) > p.capacity {
		old := p.order[0]
		p.order = p.order[1:]
		delete(p.cache, old)
	}
	return isValid, nil
}
