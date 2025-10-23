package dict

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/luispellizzon/pangram/internal/logger"
)

type Repository interface { Has(word string) (bool, error) }

type jsonMap struct{ data map[string]struct{} }

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

type JSONAdapter struct{ inner *jsonMap }

func NewJSONAdapter(path string) (*JSONAdapter, error) {
	data, err := loadJSON(path); if err != nil { return nil, err }
	return &JSONAdapter{inner: data}, nil
}

func (a *JSONAdapter) Has(word string) (bool, error) {
	_, ok := a.inner.data[strings.ToLower(word)]
	return ok, nil
}

type CacheProxy struct {
	repo Repository
	cache map[string]bool
	order []string
	capacity   int
}

func NewCacheProxy(repo Repository, capacity int) *CacheProxy {
	return &CacheProxy{repo: repo, cache: map[string]bool{}, capacity: capacity}
}
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

	// add to cache and clean older values
	p.cache[pangram] = isValid
	p.order = append(p.order, pangram)
	if len(p.order) > p.capacity {
		old := p.order[0]
		p.order = p.order[1:]
		delete(p.cache, old)
	}
	return isValid, nil
}
