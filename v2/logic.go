package v2

import (
	"math/rand"
	"time"

	"github.com/samyfodil/tb_library_snake_001/types"
)

var LookStepsAhead = 2

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

func predictSnakesNextPositions(state types.GameState) types.Board {
	board := state.Board
	for i, snake := range board.Snakes {
		// Skip dead snakes
		if snake.Health <= 1 || (isCoordInList(snake.Body[0], state.Board.Hazards) && snake.Health <= 16) {
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

func getSafeMoves(state types.GameState, head types.Coord, body []types.Coord) []string {
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

		// Check if the new head position is in a hazard
		if isCoordInList(newHead, board.Hazards) {
			if !isHeadInHazard {
				moveScores[move] = -500
			} else {
				moveScores[move] = 0
			}
			continue
		}

		// Default score for a safe move
		moveScores[move] = 100
	}

	// Find the best move based on scores
	bestMove := ""
	bestScore := -1001
	for move, score := range moveScores {
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	// If there's no best move, return an empty slice
	if bestScore <= -1000 {
		return []string{}
	}

	return []string{bestMove}
}

func isCoordInSnakeLists(state types.GameState, coord types.Coord) bool {
	for _, snake := range state.Board.Snakes {
		// Skip the dead snakes
		if snake.Health <= 1 || (isCoordInList(snake.Body[0], state.Board.Hazards) && snake.Health <= 16) {
			continue
		}

		if isCoordInList(coord, snake.Body) {
			return true
		}
	}
	return false
}

func distance(a, b types.Coord) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func countSegmentsInHazard(snake types.Battlesnake, board types.Board) int {
	segmentsInHazard := 0
	for _, segment := range snake.Body {
		for _, hazard := range board.Hazards {
			if segment.X == hazard.X && segment.Y == hazard.Y {
				segmentsInHazard++
			}
		}
	}
	return segmentsInHazard
}

func freeSpaceRatio(state types.GameState) float64 {
	totalSpaces := state.Board.Width * state.Board.Height
	occupiedSpaces := 0

	for _, snake := range state.Board.Snakes {
		occupiedSpaces += len(snake.Body)
	}

	freeSpaces := totalSpaces - occupiedSpaces
	return float64(freeSpaces) / float64(totalSpaces)
}

func shuffleMoves(moves []string) {
	for i := len(moves) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}
}

func chooseBestMove(state types.GameState, safeMoves []string) string {
	myHead := state.You.Head
	minDist := state.Board.Width*state.Board.Height + 1
	maxDist := -1
	bestMove := safeMoves[0]

	// Shuffle safeMoves to add randomness
	shuffleMoves(safeMoves)

	// Calculate the number of snake body segments in the hazard area
	segmentsInHazard := countSegmentsInHazard(state.You, state.Board)
	hazardWeight := float64(segmentsInHazard) * 1.5 // Adjust the multiplier as needed to fine-tune the prioritization

	// Calculate the free space ratio
	spaceRatio := freeSpaceRatio(state)

	// Calculate the dynamic health threshold
	healthThreshold := 30 + int((1-spaceRatio)*40) // This threshold ranges between 30 and 70 based on spaceRatio

	// Prioritize getting food when health is below the dynamic threshold
	shouldGetFood := state.You.Health < healthThreshold

	for _, move := range safeMoves {
		newHead := applyMove(myHead, move)

		// Check if the new head position is in a hazard
		inHazard := isCoordInList(newHead, state.Board.Hazards)

		// Move out of the hazard area when more of the snake's body is in it
		if inHazard {
			if hazardWeight > 0 {
				hazardWeight -= 1
				continue
			}
		}

		// If the snake should get food, prioritize the moves that minimize the distance to food
		// If the snake should not get food, prioritize the moves that maximize the distance to food
		for _, food := range state.Board.Food {
			dist := distance(newHead, food)
			if shouldGetFood {
				if dist < minDist {
					minDist = dist
					bestMove = move
				}
			} else {
				if dist > maxDist {
					maxDist = dist
					bestMove = move
				}
			}
		}
	}

	return bestMove
}

func isMoveSafeAfterNSteps(state types.GameState, move string, steps int) bool {
	if steps == 0 {
		return true
	}

	// Apply the move to the current head position
	newHead := applyMove(state.You.Head, move)

	// Check if the new head position is inside the board
	if newHead.X < 0 || newHead.X >= state.Board.Width || newHead.Y < 0 || newHead.Y >= state.Board.Height {
		return false
	}

	// Create a new state where our snake has made the move
	newBody := append([]types.Coord{newHead}, state.You.Body[:len(state.You.Body)-1]...)
	newState := state
	newState.You.Body = newBody
	newState.You.Head = newHead

	// Predict the next positions of all snakes, including our own
	newState.Board = predictSnakesNextPositions(newState)

	// Get the safe moves for the new state
	safeMoves := getSafeMoves(newState, newState.You.Head, newState.You.Body)

	// If there are no safe moves left in the new state, the initial move is not safe
	if len(safeMoves) == 0 {
		return false
	}

	// Check if the moves are safe after N-1 steps
	for _, nextMove := range safeMoves {
		if !isMoveSafeAfterNSteps(newState, nextMove, steps-1) {
			return false
		}
	}

	return true
}

func Move(state types.GameState) types.BattlesnakeMoveResponse {
	// Get safe moves for our snake based on the current state
	safeMoves := getSafeMoves(state, state.You.Head, state.You.Body)

	// Filter out moves that would not be safe after N steps
	safeMovesAfterNSteps := make([]string, 0, len(safeMoves))
	for _, move := range safeMoves {
		if isMoveSafeAfterNSteps(state, move, LookStepsAhead) {
			safeMovesAfterNSteps = append(safeMovesAfterNSteps, move)
		}
	}

	// If there are no safe moves left after filtering, fall back to the initial safe moves
	if len(safeMovesAfterNSteps) == 0 {
		safeMovesAfterNSteps = safeMoves
	}

	// Choose the best move based on your criteria (e.g., move towards food)
	nextMove := chooseBestMove(state, safeMovesAfterNSteps)

	return types.BattlesnakeMoveResponse{Move: nextMove}
}
