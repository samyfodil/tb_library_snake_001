package v3

import (
	"sort"

	"github.com/samyfodil/tb_library_snake_001/types"
)

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

func getAdjacentCoords(coord types.Coord) []types.Coord {
	return []types.Coord{
		{coord.X, coord.Y + 1},
		{coord.X, coord.Y - 1},
		{coord.X - 1, coord.Y},
		{coord.X + 1, coord.Y},
	}
}

func isCoordInList(coord types.Coord, list []types.Coord) bool {
	for _, c := range list {
		if c == coord {
			return true
		}
	}
	return false
}

func isCoordInBounds(coord types.Coord, width int, height int) bool {
	return coord.X >= 0 && coord.X < width && coord.Y >= 0 && coord.Y < height
}

func createBoard(state types.GameState) [][]float64 {
	board := make([][]float64, state.Board.Height)
	for i := range board {
		board[i] = make([]float64, state.Board.Width)
	}

	// Assign scores to cells
	for y := 0; y < state.Board.Height; y++ {
		for x := 0; x < state.Board.Width; x++ {
			coord := types.Coord{X: x, Y: y}

			// Score for snake bodies (including ours): 0
			if isCoordInSnakeLists(state, coord) {
				board[y][x] = 0
				continue
			}

			// Score for food: 1
			if isCoordInList(coord, state.Board.Food) {
				board[y][x] = 1
				continue
			}

			// Score for cells close to borders: 0.5
			if x == 0 || x == state.Board.Width-1 || y == 0 || y == state.Board.Height-1 {
				board[y][x] = 0.5
				continue
			}

			// Score for hazard cells: 0 to 0.5
			if isCoordInList(coord, state.Board.Hazards) {
				board[y][x] = 0.25 // You can adjust this value to control the hazard score
				continue
			}

			// Default score
			board[y][x] = 0
		}
	}

	return board
}

func isCoordInSnakeLists(state types.GameState, coord types.Coord) bool {
	for _, snake := range state.Board.Snakes {
		if isCoordInList(coord, snake.Body) {
			return true
		}
	}
	return false
}

func nextMove(head types.Coord, board [][]float64) string {
	adjacentCoords := getAdjacentCoords(head)
	sort.Slice(adjacentCoords, func(i, j int) bool {
		return board[adjacentCoords[i].Y][adjacentCoords[i].X] > board[adjacentCoords[j].Y][adjacentCoords[j].X]

	})

	// Filter out moves that would go out of bounds
	width := len(board[0])
	height := len(board)
	adjacentCoords = filter(adjacentCoords, func(coord types.Coord) bool {
		return isCoordInBounds(coord, width, height)
	})

	bestCoord := adjacentCoords[0]
	if head.X < bestCoord.X {
		return "right"
	} else if head.X > bestCoord.X {
		return "left"
	} else if head.Y < bestCoord.Y {
		return "up"
	} else {
		return "down"
	}
}

func filter(coords []types.Coord, condition func(types.Coord) bool) []types.Coord {
	filtered := make([]types.Coord, 0, len(coords))
	for _, coord := range coords {
		if condition(coord) {
			filtered = append(filtered, coord)
		}
	}
	return filtered
}

func calculateFutureBoards(state types.GameState, n int) [][][]float64 {
	futureBoards := make([][][]float64, n)

	for i := 0; i < n; i++ {
		futureState := state.Copy()
		for j, snake := range futureState.Board.Snakes {
			futureState.Board.Snakes[j].Body = snake.Body[:len(snake.Body)-1]
			for _, nhead := range getAdjacentCoords(snake.Head) {
				if isCoordInBounds(nhead, state.Board.Width, state.Board.Height) {
					futureState.Board.Snakes[j].Body = append(futureState.Board.Snakes[j].Body, nhead) // does not matter if it's added in the end
				}
			}
		}
		futureBoards[i] = createBoard(futureState)
	}

	return futureBoards
}

func averageBoards(boards [][][]float64) [][]float64 {
	if len(boards) == 0 {
		return nil
	}

	height := len(boards[0])
	width := len(boards[0][0])
	averagedBoard := make([][]float64, height)
	for i := range averagedBoard {
		averagedBoard[i] = make([]float64, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var sum float64
			for _, board := range boards {
				sum += board[y][x]
			}
			averagedBoard[y][x] = sum / float64(len(boards))
		}
	}

	return averagedBoard
}

func Move(state types.GameState) types.BattlesnakeMoveResponse {
	N := 4 // Number of possible future boards
	futureBoards := calculateFutureBoards(state, N)
	averagedBoard := averageBoards(futureBoards)

	me := state.You
	head := me.Head
	move := nextMove(head, averagedBoard)

	return types.BattlesnakeMoveResponse{
		Move: move,
	}
}
