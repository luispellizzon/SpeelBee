package games

import (
	"fmt"

	"github.com/luispellizzon/pangram/internal/dict"
	"github.com/luispellizzon/pangram/internal/pangram"
	"github.com/luispellizzon/pangram/internal/score"
)

// Single-player wrapper over the pangramGame itself
type pangramSingle struct{ core Game }

// Return a new Game instance
func NewPangramSingle(board pangram.GameBoard, repo dict.Repository, scorer score.Scorer) Game {
	return &pangramSingle{core: NewPangramFromGameBoard(board, repo, scorer)}
}

// Implementing Game interface
func (game *pangramSingle) Name() string { return fmt.Sprintf("%s - %s", game.core.Name(), "SINGLE PLAYER") }
func (game *pangramSingle) Info() ([]rune, rune) { return game.core.Info() }

func (game *pangramSingle) Submit(word string) (bool, string, int, int, bool) {
	return game.core.Submit(word)
}
