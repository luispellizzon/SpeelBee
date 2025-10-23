package main

import (
	"context"
	"fmt"
	"log"
	"net"

	gamepb "github.com/luispellizzon/pangram/api/pangram/v1"
	"github.com/luispellizzon/pangram/internal/dict"
	"github.com/luispellizzon/pangram/internal/games"
	"github.com/luispellizzon/pangram/internal/logger"
	"github.com/luispellizzon/pangram/internal/manager"
	"github.com/luispellizzon/pangram/internal/pangram"
	"github.com/luispellizzon/pangram/internal/score"
	"google.golang.org/grpc"
)

type server struct {
	gamepb.UnimplementedGameManagerServer
	mgr manager.Manager
}

func (s *server) CreateGame(ctx context.Context, req *gamepb.CreateGameRequest) (*gamepb.CreateGameResponse, error) {
	game_id, game, err := s.mgr.Create(req.GetKind())
	if err != nil { return nil, err }
	letters, center := game.Info()
	converted_letters := make([]string, 0, len(letters))
	for _, char := range letters { converted_letters = append(converted_letters, string(char)) }
	logger.Log().Infof("NEW GAME CREATED - ID: %v", game_id)
	return &gamepb.CreateGameResponse{Id: game_id, Name: game.Name(), Letters: converted_letters, Center: string(center)}, nil
}

func (s *server) SubmitWord(ctx context.Context, req *gamepb.SubmitWordRequest) (*gamepb.SubmitWordResponse, error) {
	game, ok := s.mgr.Get(req.GetId())
	if !ok {
		 logger.Log().Errorf("GAME NOT FOUND")
		return nil, fmt.Errorf("game not found") 
	}
	valid, reason, pts, total, pangram := game.Submit(req.GetWord())
	return &gamepb.SubmitWordResponse{
		Valid: valid, Reason: toEnum(reason), Points: int32(pts), Total: int32(total), Pangram: pangram,
	}, nil
}

func toEnum(response string) gamepb.WordResult {
	switch response {
	case "OK": return gamepb.WordResult_OK
	case "TOO_SHORT": return gamepb.WordResult_TOO_SHORT
	case "INVALID_LETTER": return gamepb.WordResult_INVALID_LETTER
	case "MISSING_CENTER": return gamepb.WordResult_MISSING_CENTER
	case "NOT_IN_DICT": return gamepb.WordResult_NOT_IN_DICT
	case "DUPLICATE": return gamepb.WordResult_DUPLICATE
	default: return gamepb.WordResult_ERROR
	}
}

func main() {
	dictPath := "assets/words_dictionary.json"
	data, err := dict.NewJSONAdapter(dictPath)
	if err != nil { logger.Log().Errorf("DICTIONARY: %v", err)}
	repo := dict.NewCacheProxy(data, 5)

	pangramPath := "assets/pangrams.json"
	words, err := pangram.LoadPangramsJSON(pangramPath)
	if err != nil || len(words) == 0 { logger.Log().Errorf("PANGRAMS: %v", err); panic(err) }
	src := pangram.CurrentTodaysPangram{Words: words}
	pangram.InitSource(src)

	scorer := score.BonusScorer{Inner: score.BasicScorer{}, Bonus: 7}
	factory := &games.Factory{Dict: repo, Scorer: scorer, Board: pangram.Provider{}}
	mgr := manager.New(factory)

	lis, err := net.Listen("tcp", ":50051"); if err != nil { log.Fatal(err) }
	s := grpc.NewServer()
	gamepb.RegisterGameManagerServer(s, &server{mgr: mgr})
	logger.Log().Infof("LISTENING ON :50051")
	if err := s.Serve(lis); err != nil { logger.Log().Errorf("SERVER %v", err) }
}