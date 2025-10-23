package pangram

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type GameBoard struct {
	Letters []rune
	Center rune
	Word string
}

type Source interface { TodaysPangram() (GameBoard, error) }

func LoadPangramsJSON(path string) ([]string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil { return nil, err }
	var objJSON map[string]any
	if err := json.Unmarshal(bytes, &objJSON); err != nil {
		return nil, err
	}
	words := make([]string, 0, len(objJSON))
	for key := range objJSON {
		key = strings.TrimSpace(key)
		if key != "" { words = append(words, key) }
	}
	return words, nil
}

// parse game pangram letters to be unique
func LettersFromWord(p string) ([]rune, error) {
	set := map[rune]struct{}{}
	letters := make([]rune, 0, 7)
	for _, r := range strings.ToLower(p) {
		if r < 'a' || r > 'z' { continue }
		if _, ok := set[r]; !ok {
			set[r] = struct{}{}; letters = append(letters, r)
		}
	}
	return letters, nil
}

type CurrentTodaysPangram struct { Words []string }

func (s CurrentTodaysPangram) TodaysPangram() (GameBoard, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	word := s.Words[rng.Intn(len(s.Words))]
	letters, err := LettersFromWord(word); if err != nil { return GameBoard{}, err }
	return GameBoard{Letters: letters, Center: letters[rng.Intn(7)], Word: word}, nil
}

// Singleton Board for everyone to read from
var (
	once sync.Once
	global GameBoard
	err	error
	src	Source
	srcSetMu sync.Mutex
)

func InitSource(s Source) {
	srcSetMu.Lock()
	defer srcSetMu.Unlock()
	if src == nil {
		src = s
	}
}

// Board returns the singleton GameBoard
func Board() (GameBoard, error) {
	once.Do(func() {
		if src == nil {
			err = errors.New("GAME-BOARD SOURCE NOT INITIALIZED")
			return
		}
		global, err = src.TodaysPangram()
	})
	return global, err
}