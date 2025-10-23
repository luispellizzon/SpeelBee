package games

import (
	"strings"

	"github.com/luispellizzon/pangram/internal/dict"
	"github.com/luispellizzon/pangram/internal/pangram"
	"github.com/luispellizzon/pangram/internal/score"
)

// Pangram Game itself to implement the Game interface. Contains attributes related to how the game can be checked against new word submissions and total points from each game. Later I can extend to fit multiple players and turn into multiplayer
type pangramGame struct {
	letters []rune
	center  rune
	seen    map[string]struct{}
	total   int
	dict    dict.Repository
	scorer  score.Scorer
}

// Create the actual user game according to what is the GameBoard singleton for every game
func NewPangramFromGameBoard(board pangram.GameBoard, repo dict.Repository, scoreStrategy score.Scorer) Game {
	return &pangramGame{
		letters: board.Letters,
		center:  board.Center,
		seen:    map[string]struct{}{},
		dict:    repo,
		scorer:  scoreStrategy,
	}
}

func (game *pangramGame) Name() string { return "PANGRAM GAME" }
func (game *pangramGame) Info() ([]rune, rune) { return game.letters, game.center }

// Main functionality of the word submission. This is called inside the singleplayer pangram wrapper. I did this way so I can extend functionality for different types of games, for example, when building the multiplayer game I can call pangramMultiplayer wrapper Submit function, and inside I can call the main pangramGame Submit to validate the word and get points, but inside the Submit wrapper, I can extend functionality after calling the main game.
func (game *pangramGame) Submit(value string) (bool, string, int, int, bool) {
	value = strings.ToLower(strings.TrimSpace(value))

	// Check word size rules
	if len([]rune(value)) < 4 { return false, "TOO_SHORT", 0, game.total, false }

	// Check if word already exists from previous submission
	if _, isDuplicated := game.seen[value]; isDuplicated { return false, "DUPLICATE", 0, game.total, false }

	allowed := map[rune]struct{}{}
	for _, chars := range game.letters { allowed[chars] = struct{}{} }

	// Check if contains the center letter and also if the letters are valid
	hasCenter := false
	for _, r := range value {
		_, ok := allowed[r]
		if !ok { return false, "INVALID_LETTER", 0, game.total, false }
		if r == game.center { hasCenter = true }
	}
	if !hasCenter { return false, "MISSING_CENTER", 0, game.total, false }

	// Check if word is in the dictionary. this will first hit the cache, and inside the cache will check the repository if not presented in the cache proxy
	ok, _ := game.dict.Has(value)
	if !ok { return false, "NOT_IN_DICT", 0, game.total, false }

	// check if length of the word is the same from the real pangram
	isSeenMap := map[rune]struct{}{}
	for _, r := range value { isSeenMap[r] = struct{}{} }
	pangram := len(isSeenMap) == len(allowed)

	// Score point according to its size. Remember the Score here is a Strategy pattern so whatever strategy we pass as a dependency injection, will represent the Score function. In this case, we inject a BonusStrategy that will take the BasicScorer and either return 1, 0 or the length of the word, and will sum up with whatever value is the bonus (+7)
	pts := game.scorer.Score(len(value), pangram)

	// Save game total points
	game.total += pts

	// Save word as seen
	game.seen[value] = struct{}{}

	// Return word response.
	return true, "OK", pts, game.total, pangram
}
