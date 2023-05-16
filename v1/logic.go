package v1

import (
	"time"

	"github.com/samyfodil/tb_library_snake_001/types"
)

func info() types.BattlesnakeInfoResponse {
	return types.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",        // TODO: Your types.Battlesnake username
		Color:      "#099a40", // TODO: Choose color
		Head:       "fang",    // TODO: Choose head
		Tail:       "bolt",    // TODO: Choose tail
	}
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func Domove(state types.GameState) types.BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	myHead := state.You.Body[0] // types.Coordinates of your head

	// Prevent your types.Battlesnake from moving out of bounds
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

	// Prevent your types.Battlesnake from colliding with itself
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

	// Prevent your types.Battlesnake from colliding with other types.Battlesnakes
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
		nextMove = "down"
	}

	return types.BattlesnakeMoveResponse{Move: nextMove}
}

func Domove2(state types.GameState) types.BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	myHead := state.You.Body[0]
	myNeck := state.You.Body[1]

	if myNeck.X < myHead.X {
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X {
		isMoveSafe["right"] = false

	} else if myNeck.Y < myHead.Y {
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y {
		isMoveSafe["up"] = false
	}

	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	// Prevent moving out of bounds
	if myHead.X == 0 {
		isMoveSafe["left"] = false
	} else if myHead.X == boardWidth-1 {
		isMoveSafe["right"] = false
	}
	if myHead.Y == 0 {
		isMoveSafe["down"] = false
	} else if myHead.Y == boardHeight-1 {
		isMoveSafe["up"] = false
	}

	// Prevent colliding with itself
	for _, coord := range state.You.Body[1:] {
		if coord.X == myHead.X {
			if coord.Y < myHead.Y {
				isMoveSafe["down"] = false
			} else {
				isMoveSafe["up"] = false
			}
		} else if coord.Y == myHead.Y {
			if coord.X < myHead.X {
				isMoveSafe["left"] = false
			} else {
				isMoveSafe["right"] = false
			}
		}
	}

	// Prevent colliding with other snakes
	for _, snake := range state.Board.Snakes {
		for _, coord := range snake.Body {
			if coord.X == myHead.X {
				if coord.Y < myHead.Y {
					isMoveSafe["down"] = false
				} else {
					isMoveSafe["up"] = false
				}
			} else if coord.Y == myHead.Y {
				if coord.X < myHead.X {
					isMoveSafe["left"] = false
				} else {
					isMoveSafe["right"] = false
				}
			}
		}
	}

	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		return types.BattlesnakeMoveResponse{Move: "down"}
	}

	// Move towards the closest food
	closestFood := findClosestFood(state.You, state.Board)
	moveTowardsFood := getMoveTowardsFood(state.You.Head, closestFood)

	chosenMove := safeMoves[0]
	for _, move := range safeMoves {
		if move == moveTowardsFood {
			chosenMove = move
			break
		}
	}

	return types.BattlesnakeMoveResponse{Move: chosenMove}

}

func findClosestFood(snake types.Battlesnake, board types.Board) types.Coord {
	closestFood := board.Food[0]
	minDistance := manhattanDistance(snake.Head, closestFood)

	for _, food := range board.Food {
		distance := manhattanDistance(snake.Head, food)
		if distance < minDistance {
			minDistance = distance
			closestFood = food
		}
	}

	return closestFood

}

func manhattanDistance(a, b types.Coord) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getMoveTowardsFood(head, food types.Coord) string {
	if head.X < food.X {
		return "right"
	}
	if head.X > food.X {
		return "left"
	}
	if head.Y < food.Y {
		return "up"
	}
	if head.Y > food.Y {
		return "down"
	}
	return ""

}

func Domove3(state types.GameState) types.BattlesnakeMoveResponse {
	// Get safe moves
	safeMoves := getSafeMoves(state)

	if len(safeMoves) == 0 {
		return types.BattlesnakeMoveResponse{Move: "down"}
	}

	// Move towards the closest food
	closestFood := findClosestFood(state.You, state.Board)
	moveTowardsFood := getMoveTowardsFood(state.You.Head, closestFood)

	chosenMove := safeMoves[0]
	for _, move := range safeMoves {
		if move == moveTowardsFood {
			chosenMove = move
			break
		}
	}

	return types.BattlesnakeMoveResponse{Move: chosenMove}
}

func getSafeMoves(state types.GameState) []string {
	myHead := state.You.Body[0]
	myBody := state.You.Body[1:]

	possibleMoves := []string{"up", "down", "left", "right"}
	safeMoves := []string{}

	for _, move := range possibleMoves {
		newHead := getNewHead(myHead, move)
		isSafe := true

		// Check for wall collisions
		if newHead.X < 0 || newHead.Y < 0 || newHead.X >= state.Board.Width || newHead.Y >= state.Board.Height {
			isSafe = false
		}

		// Check for self-collisions
		if isSafe {
			for _, coord := range myBody {
				if newHead.X == coord.X && newHead.Y == coord.Y {
					isSafe = false
					break
				}
			}
		}

		// Check for opponent collisions
		if isSafe {
			for _, snake := range state.Board.Snakes {
				if snake.ID == state.You.ID {
					continue
				}
				for _, coord := range snake.Body {
					if newHead.X == coord.X && newHead.Y == coord.Y {
						isSafe = false
						break
					}
				}
				if !isSafe {
					break
				}
			}
		}

		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	return safeMoves
}

func isOutOfBoard(width, height int, coord types.Coord) bool {
	return coord.X < 0 || coord.X >= width || coord.Y < 0 || coord.Y >= height
}

func isSnakeCollision(board types.Board, coord types.Coord) bool {
	for _, snake := range board.Snakes {
		for _, snakeBody := range snake.Body {
			if coord.X == snakeBody.X && coord.Y == snakeBody.Y {
				return true
			}
		}
	}
	return false
}

func getNewHead(head types.Coord, move string) types.Coord {
	newHead := types.Coord{X: head.X, Y: head.Y}
	switch move {
	case "up":
		newHead.Y++
	case "down":
		newHead.Y--
	case "left":
		newHead.X--
	case "right":
		newHead.X++
	}
	return newHead
}

/******************************/

func getSafeMovesFromOpponents(myHead types.Coord, safeMoves []string, opponentMoves map[string][]types.Coord) []string {
	safestMoves := make([]string, 0)

	for _, move := range safeMoves {
		newHead := getNewHead(myHead, move)
		isSafe := true

		for _, opponentMoveSet := range opponentMoves {
			for _, opponentMove := range opponentMoveSet {
				if newHead.X == opponentMove.X && newHead.Y == opponentMove.Y {
					isSafe = false
					break
				}
			}
			if !isSafe {
				break
			}
		}

		if isSafe {
			safestMoves = append(safestMoves, move)
		}
	}

	return safestMoves
}

func getAllOpponentMoves(state types.GameState) map[string][]types.Coord {
	opponentMoves := make(map[string][]types.Coord)

	for _, snake := range state.Board.Snakes {
		if snake.ID != state.You.ID {
			possibleMoves := getPossibleMoves(snake.Head)
			opponentMoves[snake.ID] = possibleMoves
		}
	}

	return opponentMoves
}

func getPossibleMoves(head types.Coord) []types.Coord {
	return []types.Coord{
		{X: head.X, Y: head.Y + 1},
		{X: head.X, Y: head.Y - 1},
		{X: head.X + 1, Y: head.Y},
		{X: head.X - 1, Y: head.Y},
	}
}

var lookAheadMoves = 8

func Domove4(state types.GameState) types.BattlesnakeMoveResponse {
	myHead := state.You.Body[0]

	// Get safe moves
	safeMoves := getSafeMoves(state)

	if len(safeMoves) == 0 {
		return types.BattlesnakeMoveResponse{Move: "down"}
	}

	// Get the positions of all possible moves for each opponent snake
	opponentMoves := getAllOpponentMoves(state)

	// Remove moves that would collide with opponents' possible moves
	safestMoves := getSafeMovesFromOpponents(myHead, safeMoves, opponentMoves)

	chosenMove := safestMoves[0]

	// If health is below the threshold, look for food
	if state.You.Health < 50 {
		closestFood := findClosestFood(state.You, state.Board)
		moveTowardsFood := getMoveTowardsFood(state.You.Head, closestFood)

		for _, move := range safestMoves {
			if move == moveTowardsFood {
				chosenMove = move
				break
			}
		}
	}

	// Simulate all possible moves for the next two turns and evaluate their safety
	safestNextMoves := make([]string, 0)
	maxSafetyScore := -1

	for _, move := range safestMoves {
		simulatedState := simulateMove(state, move)
		safetyScore, nextMove := getNextMoveSafetyScore(simulatedState, opponentMoves)

		if safetyScore > maxSafetyScore {
			maxSafetyScore = safetyScore
			safestNextMoves = []string{nextMove}
		} else if safetyScore == maxSafetyScore {
			safestNextMoves = append(safestNextMoves, nextMove)
		}
	}

	// Choose a random move from the safest next moves
	if len(safestNextMoves) > 0 {
		chosenMove = safestNextMoves[0]
	}

	return types.BattlesnakeMoveResponse{Move: chosenMove}
}

func deepCopyGameState(original types.GameState) types.GameState {
	copied := types.GameState{
		Game: original.Game,
		Turn: original.Turn,
		Board: types.Board{
			Height: original.Board.Height,
			Width:  original.Board.Width,
			Food:   make([]types.Coord, len(original.Board.Food)),
			Snakes: make([]types.Battlesnake, len(original.Board.Snakes)),
		},
		You: original.You,
	}

	copy(copied.Board.Food, original.Board.Food)

	for i, snake := range original.Board.Snakes {
		copiedSnake := types.Battlesnake{
			ID:     snake.ID,
			Name:   snake.Name,
			Health: snake.Health,
			Body:   make([]types.Coord, len(snake.Body)),
		}
		copy(copiedSnake.Body, snake.Body)
		copied.Board.Snakes[i] = copiedSnake
	}

	return copied
}

func simulateMove(state types.GameState, move string) types.GameState {
	simulatedState := deepCopyGameState(state)
	newHead := getNewHead(simulatedState.You.Body[0], move)
	simulatedState.You.Body = append([]types.Coord{newHead}, simulatedState.You.Body...)
	return simulatedState
}

func getNextMoveSafetyScore(state types.GameState, opponentMoves map[string][]types.Coord) (int, string) {
	myHead := state.You.Body[0]
	safeMoves := getSafeMoves(state)

	safetyScores := make(map[string]int)

	for _, move := range safeMoves {
		newHead := getNewHead(myHead, move)
		isSafe := true

		for _, opponentMoveSet := range opponentMoves {
			for _, opponentMove := range opponentMoveSet {
				if newHead.X == opponentMove.X && newHead.Y == opponentMove.Y {
					isSafe = false
					break
				}
			}
			if !isSafe {
				break
			}
		}

		if isSafe {
			safetyScores[move] = 1
		} else {
			safetyScores[move] = 0
		}
	}

	maxSafetyScore := -1
	bestMove := ""

	for move, score := range safetyScores {
		if score > maxSafetyScore {
			maxSafetyScore = score
			bestMove = move
		} else if score == maxSafetyScore {
			// If the random number is 0, update the bestMove
			if time.Now().Unix()%2 == 0 {
				bestMove = move
			}
		}
	}

	return maxSafetyScore, bestMove
}

/**************************************************************/

// Move function
func Domove5(state types.GameState) types.BattlesnakeMoveResponse {
	opponentMoves := getAllOpponentMoves(state)
	safetyScore, chosenMove := getNextMoveSafetyScoreV2(state, opponentMoves)

	// Reduce lookAheadMoves until a move is found
	for safetyScore == -1 && lookAheadMoves > 0 {
		lookAheadMoves--
		safetyScore, chosenMove = getNextMoveSafetyScoreV2(state, opponentMoves)
	}

	// If no safe move is found, choose a random move from all possible moves
	if safetyScore == -1 {
		possibleMoves := []string{"up", "down", "left", "right"}
		chosenMove = possibleMoves[int(time.Now().Unix())%len(possibleMoves)]
	}

	return types.BattlesnakeMoveResponse{Move: chosenMove}
}

func simulateMoveV2(state types.GameState, move string, lookAheadMoves int) types.GameState {
	if lookAheadMoves <= 0 {
		return state
	}

	simulatedState := deepCopyGameState(state)
	newHead := getNewHead(simulatedState.You.Body[0], move)
	simulatedState.You.Body = append([]types.Coord{newHead}, simulatedState.You.Body[:len(simulatedState.You.Body)-1]...)

	// Recursively simulate moves
	opponentMoves := getAllOpponentMoves(simulatedState)
	for _, moves := range opponentMoves {
		for _, coord := range moves {
			m := getDirection(simulatedState.You.Body[0], coord)
			simulateMoveV2(simulatedState, m, lookAheadMoves-1)
		}
	}

	return simulatedState
}

func getDirection(currentHead, newHead types.Coord) string {
	if newHead.X < currentHead.X {
		return "left"
	} else if newHead.X > currentHead.X {
		return "right"
	} else if newHead.Y < currentHead.Y {
		return "down"
	} else {
		return "up"
	}
}

// getNextMoveSafetyScoreV2 function
func getNextMoveSafetyScoreV2(state types.GameState, opponentMoves map[string][]types.Coord) (int, string) {
	myHead := state.You.Body[0]
	safeMoves := getSafeMoves(state)

	safetyScores := make(map[string]int)

	for _, move := range safeMoves {
		newHead := getNewHead(myHead, move)
		isSafe := true

		for _, opponentMoveSet := range opponentMoves {
			for _, opponentMove := range opponentMoveSet {
				if newHead.X == opponentMove.X && newHead.Y == opponentMove.Y {
					isSafe = false
					break
				}
			}
			if !isSafe {
				break
			}
		}

		if isSafe {
			simulatedState := simulateMoveV2(state, move, lookAheadMoves-1)
			opponentSimulatedMoves := getAllOpponentMoves(simulatedState)
			safetyScore, _ := getNextMoveSafetyScore(simulatedState, opponentSimulatedMoves)
			safetyScores[move] = safetyScore + 1
		} else {
			safetyScores[move] = 0
		}
	}

	maxSafetyScore := -1
	bestMove := ""

	for move, score := range safetyScores {
		if score > maxSafetyScore {
			maxSafetyScore = score
			bestMove = move
		} else if score == maxSafetyScore {
			// If the random number is 0, update the bestMove
			if time.Now().Unix()%2 == 0 {
				bestMove = move
			}
		}
	}

	return maxSafetyScore, bestMove
}

/************************************************/

const maxCalculationTime = 30 * time.Millisecond

func Domove6(state types.GameState) types.BattlesnakeMoveResponse {
	myHead := state.You.Body[0]
	opponentMoves := getAllOpponentMoves(state)
	var chosenMove string
	var collisionDetected bool
	var bestSafetyScore int = -1

	startTime := time.Now()

	for {
		currentSafetyScore, currentMove := getNextMoveSafetyScoreV2(state, opponentMoves)
		if currentSafetyScore > bestSafetyScore {
			bestSafetyScore = currentSafetyScore
			chosenMove = currentMove
		}

		newHead := getNewHead(myHead, chosenMove)
		collisionDetected = false

		// Check for wall collisions
		if newHead.X < 0 || newHead.Y < 0 || newHead.X >= state.Board.Width || newHead.Y >= state.Board.Height {
			collisionDetected = true
		}

		// Check for self-collisions
		if !collisionDetected {
			for _, coord := range state.You.Body {
				if newHead.X == coord.X && newHead.Y == coord.Y {
					collisionDetected = true
					break
				}
			}
		}

		// Check for opponent collisions
		if !collisionDetected {
			for _, snake := range state.Board.Snakes {
				if snake.ID == state.You.ID {
					continue
				}
				for _, coord := range snake.Body {
					if newHead.X == coord.X && newHead.Y == coord.Y {
						collisionDetected = true
						break
					}
				}
				if collisionDetected {
					break
				}
			}
		}

		if !collisionDetected || time.Since(startTime) > maxCalculationTime {
			break
		}
	}

	return types.BattlesnakeMoveResponse{Move: chosenMove}
}
