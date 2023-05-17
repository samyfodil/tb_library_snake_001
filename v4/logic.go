package v4

import (
	"math"
	"math/rand"
	"time"

	"github.com/samyfodil/tb_library_snake_001/types"
)

// Initialize the random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Helper functions

func applyMove(head types.Coord, move string) types.Coord {
	switch move {
	case "up":
		head.Y++
	case "down":
		head.Y--
	case "left":
		head.X--
	case "right":
		head.X++
	}
	return head
}

// Main logic

func isCoordInList(coord types.Coord, list []types.Coord) bool {
	for _, c := range list {
		if c == coord {
			return true
		}
	}
	return false
}

func predictSnakesNextPositions(state *types.GameState) types.Board {
	board := state.Board
	for i, snake := range board.Snakes {
		// Remove dead snakes
		if snake.Health <= 1 || (isCoordInList(snake.Body[0], state.Board.Hazards) && snake.Health <= 16) {
			board.Snakes = append(board.Snakes[:i], board.Snakes[i+1:]...)
			i-- // Adjust the loop index since we removed an element
			continue
		}

		safeMoves := getSafeMoves(state, snake.Head, snake.Body)
		if len(safeMoves) > 0 {
			move := safeMoves[rand.Intn(len(safeMoves))]
			newHead := applyMove(snake.Head, move)
			board.Snakes[i].Body = append([]types.Coord{newHead}, snake.Body[:len(snake.Body)-1]...)

			if isCoordInList(newHead, board.Food) {
				board.Snakes[i].Body = append([]types.Coord{newHead}, snake.Body...)
				board.Snakes[i].Health = 100
			} else {
				board.Snakes[i].Health -= 1
			}
		} else {
			board.Snakes[i].Health = 0
		}
	}

	return board
}

func getSafeMoves(state *types.GameState, head types.Coord, body []types.Coord) []string {
	board := state.Board
	possibleMoves := []string{"up", "down", "left", "right"}

	// Initialize move scores
	moveScores := make(map[string]int)

	isHeadInHazard := isCoordInList(head, board.Hazards)

	for _, move := range possibleMoves {
		newHead := applyMove(head, move)

		// Check if the new head position is out of the board
		if newHead.X < 0 || newHead.X >= board.Width || newHead.Y < 0 || newHead.Y >= board.Height {
			moveScores[move] = -1000
			continue
		}

		// Check if the new head position is in your own body
		if isCoordInList(newHead, body) {
			moveScores[move] = -1000
			continue
		}

		// Check if the new head position is in another snake's body
		if isCoordInSnakeLists(state, newHead) {
			moveScores[move] = -1000
			continue
		}

		// Check for possible head-to-head collisions with other snakes
		for _, otherSnake := range state.Board.Snakes {
			if otherSnake.ID == state.You.ID {
				continue
			}
			for _, otherMove := range possibleMoves {
				otherNewHead := applyMove(otherSnake.Head, otherMove)
				if newHead == otherNewHead {
					if len(body) > len(otherSnake.Body) {
						moveScores[move] = -500
					} else {
						moveScores[move] = -1000
					}
					break
				}
			}
		}

		// Add hazard penalty
		if isHeadInHazard && isCoordInList(newHead, board.Hazards) {
			moveScores[move] -= 200
		}
	}

	// Find the best moves
	bestScore := math.MinInt32
	bestMoves := []string{}
	for move, score := range moveScores {
		if score > bestScore {
			bestScore = score
			bestMoves = []string{move}
		} else if score == bestScore {
			bestMoves = append(bestMoves, move)
		}
	}

	return bestMoves
}

func isCoordInSnakeLists(state *types.GameState, coord types.Coord) bool {
	for _, snake := range state.Board.Snakes {
		if isCoordInList(coord, snake.Body) {
			return true
		}
	}
	return false
}

// Minimax

const maxDepth = 3

func calculateVoronoi(board *types.Board, snakes []types.Battlesnake) map[string]int {
	voronoi := make(map[string]int)

	for x := 0; x < board.Width; x++ {
		for y := 0; y < board.Height; y++ {
			minDistance := math.MaxInt32
			minSnakeID := ""

			for _, snake := range snakes {
				distance := abs(snake.Head.X-x) + abs(snake.Head.Y-y)
				if distance < minDistance {
					minDistance = distance
					minSnakeID = snake.ID
				}
			}

			voronoi[minSnakeID]++
		}
	}

	return voronoi
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func evaluateState(state *types.GameState) float64 {
	// Basic evaluation function that compares your snake's length to the other snakes' lengths and Voronoi territories
	mySnake := state.You
	myScore := float64(len(mySnake.Body))

	voronoi := calculateVoronoi(&state.Board, state.Board.Snakes)
	myScore += float64(voronoi[mySnake.ID])

	for _, otherSnake := range state.Board.Snakes {
		if otherSnake.ID == mySnake.ID {
			continue
		}
		myScore -= float64(len(otherSnake.Body))
		myScore -= float64(voronoi[otherSnake.ID])
	}

	return myScore
}

func minimax(state *types.GameState, depth int, isMaximizing bool) float64 {
	if depth == 0 {
		return evaluateState(state)
	}
	if isMaximizing {
		maxEval := math.Inf(-1)
		moves := getSafeMoves(state, state.You.Head, state.You.Body)

		for _, move := range moves {
			newState := state.Copy()
			newState.You.Head = applyMove(newState.You.Head, move)
			newState.Board = predictSnakesNextPositions(newState)
			eval := minimax(newState, depth-1, false)
			maxEval = math.Max(maxEval, eval)
		}

		return maxEval
	} else {
		minEval := math.Inf(1)
		moves := getSafeMoves(state, state.You.Head, state.You.Body)

		for _, move := range moves {
			newState := state.Copy()
			newState.You.Head = applyMove(newState.You.Head, move)
			newState.Board = predictSnakesNextPositions(newState)
			eval := minimax(newState, depth-1, true)
			minEval = math.Min(minEval, eval)
		}

		return minEval
	}
}

func Move(state *types.GameState) types.BattlesnakeMoveResponse {
	bestScore := math.Inf(-1)
	bestMove := ""

	moves := getSafeMoves(state, state.You.Head, state.You.Body)
	for _, move := range moves {
		newState := state.Copy()
		newState.You.Head = applyMove(newState.You.Head, move)
		newState.Board = predictSnakesNextPositions(newState)
		score := minimax(newState, maxDepth-1, false)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	return types.BattlesnakeMoveResponse{Move: bestMove}
}
