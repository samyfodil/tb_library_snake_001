package lib

import (
	"crypto/rand"
	"math/big"

	wrand "github.com/taubyte/go-sdk/crypto/rand"
)

func info() BattlesnakeInfoResponse {
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",           // TODO: Your Battlesnake username
		Color:      "#099a40",    // TODO: Choose color
		Head:       "all-seeing", // TODO: Choose head
		Tail:       "mlh-gene",   // TODO: Choose tail
	}
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func domove(state GameState) BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	myHead := state.You.Body[0] // Coordinates of your head

	// Prevent your Battlesnake from moving out of bounds
	boardWidth := state.Board.Width

	if myHead.X == boardWidth-1 {
		isMoveSafe["right"] = false
	} else if myHead.X == 0 {
		isMoveSafe["left"] = false
	}

	boardHeight := state.Board.Height

	if myHead.Y == 0 {
		isMoveSafe["down"] = false
	} else if myHead.Y == boardHeight-1 {
		isMoveSafe["up"] = false
	}

	// Prevent your Battlesnake from colliding with itself
	for _, cor := range state.You.Body[1:] /*skip head*/ {
		x, y := cor.X, cor.Y

		if myHead.X+1 == x {
			isMoveSafe["right"] = false
		} else if myHead.X-1 == x {
			isMoveSafe["left"] = false
		}

		if myHead.Y-1 == y {
			isMoveSafe["down"] = false
		} else if myHead.Y+1 == y {
			isMoveSafe["up"] = false
		}
	}

	// Prevent your Battlesnake from colliding with other Battlesnakes
	for _, snk := range state.Board.Snakes {
		for _, cor := range snk.Body {
			x, y := cor.X, cor.Y

			if myHead.X+1 == x {
				isMoveSafe["right"] = false
			} else if myHead.X-1 == x {
				isMoveSafe["left"] = false
			}

			if myHead.Y-1 == y {
				isMoveSafe["down"] = false
			} else if myHead.Y+1 == y {
				isMoveSafe["up"] = false
			}
		}
	}

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	// if none, include all
	if len(safeMoves) == 0 {
		for move, _ := range isMoveSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	// Choose a random move from the safe ones
	var nextMove string

	// Step 4 - Move towards food instead of random, to regain health and survive longer
	scoredSafeMoves := make(map[string]int)
	for _, food := range state.Board.Food {
		for _, smv := range safeMoves {
			x, y := myHead.X, myHead.Y
			switch smv {
			case "up":
				y++
			case "down":
				y--
			case "right":
				x++
			case "left":
				x--
			}
			old_score := scoredSafeMoves[smv]
			score := (food.X-x)*(food.X-x) + (food.Y-y)*(food.Y-y)
			if score < old_score {
				scoredSafeMoves[smv] = score
			}
		}
	}

	if len(scoredSafeMoves) > 0 {
		best := state.Board.Width*state.Board.Width + state.Board.Height*state.Board.Height
		for smv, score := range scoredSafeMoves {
			if score < best {
				best = score
				nextMove = smv
			}
		}
	} else {
		i, _ := rand.Int(wrand.NewReader(), big.NewInt(3000))
		nextMove = safeMoves[i.Int64()%int64(len(safeMoves))]
	}

	return BattlesnakeMoveResponse{Move: nextMove}
}
