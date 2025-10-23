package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	gamepb "github.com/luispellizzon/pangram/api/pangram/v1"
	"github.com/luispellizzon/pangram/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	gameID := flag.String("game_id", "", "--game_id flag to rejoin game, or create new game if not specified in the terminal")
	gameMode := flag.String("mode", "", "--mode flag for singleplayer OR multiplayer")
	flag.Parse()
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { logger.Log().Errorf("SERVER %v", err) }
	defer conn.Close()

	client := gamepb.NewGameManagerClient(conn)

	var id string
	// create new game
	if *gameID == "" {
		if *gameMode == "" {
			// prompt for mode singleplayer vs multiplayer
			*gameMode = getMode()
		}

		// check if mode chosen is valid
		if !isValidMode(*gameMode) {
			fmt.Printf("Invalid mode: %q. Use 'singleplayer' or 'multiplayer'.\n", *gameMode)
			os.Exit(2)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Create new game
		response, err := client.CreateGame(ctx, &gamepb.CreateGameRequest{Kind: *gameMode})
		if err != nil { panic(err) }

		// Get new game id and game info about the pangram
		id = response.GetId()
		fmt.Printf("Game ID: %s %s \nletters: %s \ncenter: %s\n",
			response.GetId(), response.GetName(), strings.Join(response.GetLetters(), " "), response.GetCenter())
	} else {
		// rejoin previous game using game id
		id = *gameID
		fmt.Printf("Joining existing game -> %s.\n", id)
	}

	// game loop
	cli := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter word (or /quit): ")
		if !cli.Scan() {
			break
		}
		w := strings.TrimSpace(cli.Text())
		if w == "" { continue }
		if w == "/quit" { break }

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		// Submit word to game by id to be checked on the server
		resp, err := client.SubmitWord(ctx, &gamepb.SubmitWordRequest{Id: id, Word: w})
		cancel()
		if err != nil {
			fmt.Printf("Response error: %v\n", err)
			continue
		}

		// Check server response
		// Check if word is valid
		if resp.GetValid() {
			// Check if the word is the pangram
			if resp.GetPangram(){
				fmt.Printf("IS PANGRAM: +%d \nTOTAL POINTS: %d\n",
				// Show points
				resp.GetPoints(), resp.GetTotal())
			} else {
				fmt.Printf("VALID: +%d \nTOTAL POINTS: %d\n",
				resp.GetPoints(), resp.GetTotal())
			}
		} else {
			fmt.Printf("INVALID: %s \nTOTAL POINTS: %d)\n", resp.GetReason().String(), resp.GetTotal())
		}
	}
}


func getMode() string {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Choose a game mode [singleplayer/multiplayer]: ")
		if !reader.Scan() {
			return "singleplayer" 
		}
		m := strings.ToLower(strings.TrimSpace(reader.Text()))
		if isValidMode(m) {
			return m
		}
		fmt.Println("Please, enter 'singleplayer' or 'multiplayer'.")
	}
}

func isValidMode(gameMode string) bool {
	return gameMode == "singleplayer" || gameMode == "multiplayer"
}