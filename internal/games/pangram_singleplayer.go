package games

import (
	"fmt"

	"github.com/luispellizzon/pangram/internal/dict"
	"github.com/luispellizzon/pangram/internal/pangram"
	"github.com/luispellizzon/pangram/internal/score"
)

// Single-player wrapper over the game.
type pangramSingle struct{ core Game }

func NewPangramSingle(board pangram.GameBoard, repo dict.Repository, scorer score.Scorer) Game {
	return &pangramSingle{core: NewPangramFromGameBoard(board, repo, scorer)}
}

func (game *pangramSingle) Name() string { return fmt.Sprintf("%s - %s", game.core.Name(), "SINGLE PLAYER") }
func (game *pangramSingle) Info() ([]rune, rune) { return game.core.Info() }

func (game *pangramSingle) Submit(word string) (bool, string, int, int, bool) {
	return game.core.Submit(word)
}
