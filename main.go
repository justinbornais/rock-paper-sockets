package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type GameRoom struct {
	Code    string
	Players [2]*websocket.Conn // Two players max
	Moves   map[int]string     // PlayerID to move
	Scores  [2]int             // Scores for both players
	Mutex   sync.Mutex         // Protect shared resources
}

var (
	games = make(map[string]*GameRoom)
	mutex sync.Mutex
)

const WinningPoints = 3

func main() {
	app := fiber.New()

	// Serve static files
	app.Static("/static", "./static")

	// Serve index.html
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./static/index.html")
	})

	// WebSocket endpoint
	app.Get("/ws/:game_code", websocket.New(func(c *websocket.Conn) {
		gameCode := c.Params("game_code")
		mutex.Lock()
		if _, exists := games[gameCode]; !exists {
			games[gameCode] = &GameRoom{
				Code:  gameCode,
				Moves: make(map[int]string),
			}
		}
		game := games[gameCode]
		mutex.Unlock()

		// Handle players
		game.Mutex.Lock()
		var playerID int
		if game.Players[0] == nil {
			playerID = 0
			game.Players[0] = c
		} else if game.Players[1] == nil {
			playerID = 1
			game.Players[1] = c
		} else {
			game.Mutex.Unlock()
			c.WriteMessage(websocket.TextMessage, []byte("Game room is full."))
			c.Close()
			return
		}
		game.Mutex.Unlock()

		// Notify players when the game starts
		if game.Players[0] != nil && game.Players[1] != nil {
			for i, player := range game.Players {
				if player != nil {
					player.WriteJSON(map[string]string{
						"action":    "game_started",
						"message":   "Game has started! Please make your move.",
						"player_id": fmt.Sprintf("%d", i),
					})
				}
			}
		}

		defer func() {
			// Handle player disconnection
			game.Mutex.Lock()
			game.Players[playerID] = nil
			if game.Players[0] == nil && game.Players[1] == nil {
				delete(games, gameCode)
			}
			game.Mutex.Unlock()
		}()

		// Main game loop
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			// Parse the received JSON
			var input map[string]string
			if err := json.Unmarshal(msg, &input); err != nil {
				log.Println("Error parsing JSON:", err)
				continue
			}

			// Check for "move" action
			action, ok := input["action"]
			if ok && action == "move" {
				move, moveExists := input["move"]
				if !moveExists {
					continue
				}
				game.Mutex.Lock()
				game.Moves[playerID] = move

				// If both players made a move
				if len(game.Moves) == 2 {
					winner := calculateWinner(game.Moves[0], game.Moves[1])
					if winner != -1 {
						game.Scores[winner]++
					}

					// Check if the game is over
					if game.Scores[0] == WinningPoints || game.Scores[1] == WinningPoints {
						for _, player := range game.Players {
							if player != nil {
								player.WriteJSON(map[string]interface{}{
									"action": "game_over",
									"winner": winner,
								})
							}
						}
						delete(games, gameCode)
						break
					}

					// Send results to players
					for i, player := range game.Players {
						if player != nil {
							player.WriteJSON(map[string]interface{}{
								"action":        "result",
								"winner":        winner,
								"scores":        game.Scores,
								"opponent_move": game.Moves[1-i],
							})
						}
					}

					// Reset moves
					game.Moves = make(map[int]string)
				}
			}
			game.Mutex.Unlock()
		}
	}))

	log.Fatal(app.Listen(":8080"))
}

// calculateWinner determines the winner based on moves.
func calculateWinner(move1, move2 string) int {
	if move1 == move2 {
		return -1 // Tie
	}
	if (move1 == "r" && move2 == "s") || (move1 == "s" && move2 == "p") || (move1 == "p" && move2 == "r") {
		return 0 // Player 1 wins
	}
	return 1 // Player 2 wins
}
