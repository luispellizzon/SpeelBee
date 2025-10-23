package games

import (
	"fmt"

	"github.com/luispellizzon/pangram/internal/dict"
	"github.com/luispellizzon/pangram/internal/pangram"
	"github.com/luispellizzon/pangram/internal/score"
)

type IBoardProvider interface{
	Board() (pangram.GameBoard, error)
}

type Factory struct {
	Dict   dict.Repository
	Scorer score.Scorer
	Board IBoardProvider
}

func (f *Factory) New(kind string) (Game, error) {
	board, err := f.Board.Board()
	if err != nil {return nil, err}
	switch kind {
	case "singleplayer":
		return NewPangramSingle(board, f.Dict, f.Scorer), nil
	case "multiplayer":
		return nil, fmt.Errorf("MULTIPLAYER NOT IMPLEMENTED")
	default:
		return nil, fmt.Errorf("GAME NOT IMPLEMENTED: %s", kind)
	}
}
