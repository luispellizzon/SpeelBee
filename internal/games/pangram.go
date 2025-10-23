package games

import (
	"strings"

	"github.com/luispellizzon/pangram/internal/dict"
	"github.com/luispellizzon/pangram/internal/pangram"
	"github.com/luispellizzon/pangram/internal/score"
)

type pangramGame struct {
	letters []rune
	center  rune
	seen    map[string]struct{}
	total   int
	dict    dict.Repository
	scorer  score.Scorer
}

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

func (game *pangramGame) Submit(value string) (bool, string, int, int, bool) {
	value = strings.ToLower(strings.TrimSpace(value))
	if len([]rune(value)) < 4 { return false, "TOO_SHORT", 0, game.total, false }
	if _, isDuplicated := game.seen[value]; isDuplicated { return false, "DUPLICATE", 0, game.total, false }

	allowed := map[rune]struct{}{}
	for _, chars := range game.letters { allowed[chars] = struct{}{} }

	hasCenter := false
	for _, r := range value {
		_, ok := allowed[r]
		if !ok { return false, "INVALID_LETTER", 0, game.total, false }
		if r == game.center { hasCenter = true }
	}
	if !hasCenter { return false, "MISSING_CENTER", 0, game.total, false }

	ok, _ := game.dict.Has(value)
	if !ok { return false, "NOT_IN_DICT", 0, game.total, false }

	isSeenMap := map[rune]struct{}{}
	for _, r := range value { isSeenMap[r] = struct{}{} }
	pangram := len(isSeenMap) == len(allowed)

	pts := game.scorer.Score(len(value), pangram)
	game.total += pts
	game.seen[value] = struct{}{}
	return true, "OK", pts, game.total, pangram
}
