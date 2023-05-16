package lib

import (
	"math"
	"math/rand"
	"sort"
)

type MoveStrategy interface {
	GetMove(state GameState) BattlesnakeMoveResponse
}

type MoveStrategyX struct {
	MoveHelper
}

func (ms MoveStrategyX) GetMove(state GameState) BattlesnakeMoveResponse {
	movesWithScores := make([]moveScore, 0, len(ms.possibleMoves))

	for _, move := range ms.possibleMoves {
		score := ms.getNextMoveSafetyScore(state, move, ms.lookAheadMoves)
		movesWithScores = append(movesWithScores, moveScore{Move: move, Score: score})
	}

	ms.randomizeMoves(movesWithScores)
	ms.sortMovesByScore(movesWithScores)

	chosenMove := movesWithScores[0].Move

	return BattlesnakeMoveResponse{Move: chosenMove}
}

type moveScore struct {
	Move  string
	Score int
}

type MoveHelper struct {
	possibleMoves  []string
	lookAheadMoves int
}

func NewMoveHelper(lookAheadMoves int) MoveHelper {
	return MoveHelper{
		possibleMoves:  []string{"up", "down", "left", "right"},
		lookAheadMoves: lookAheadMoves,
	}
}

func (mh MoveHelper) getNextMoveSafetyScore(state GameState, move string, lookAheadMoves int) int {
	score := 0

	newHead := mh.getNewHead(state.You.Body[0], move)

	if mh.isOutOfBounds(newHead, state.Board.Width, state.Board.Height) {
		return math.MinInt32
	}

	if mh.isCollidingWithSnakes(newHead, state.Board.Snakes) {
		collidingSnake := mh.isHeadToHeadCollision(newHead, state.Board.Snakes)
		if collidingSnake == nil || len(state.You.Body) <= len(collidingSnake.Body) {
			return math.MinInt32
		}
		score += 1000 // Increase the score for colliding with a smaller snake's head
	}

	if lookAheadMoves > 0 {
		for _, m := range mh.possibleMoves {
			simulatedState := mh.simulateMove(state, move, lookAheadMoves-1)
			score += mh.getNextMoveSafetyScore(simulatedState, m, lookAheadMoves-1)
		}
	}

	score += mh.getAreaControlScore(state, move)

	return score
}

// Helper methods

func (mh MoveHelper) getNewHead(head Coord, move string) Coord {
	newHead := head
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

func (mh MoveHelper) isOutOfBounds(coord Coord, width, height int) bool {
	return coord.X < 0 || coord.Y < 0 || coord.X >= width || coord.Y >= height
}

func (mh MoveHelper) isCollidingWithSnakes(coord Coord, snakes []Battlesnake) bool {
	for _, snake := range snakes {
		for _, snakeCoord := range snake.Body {
			if coord.X == snakeCoord.X && coord.Y == snakeCoord.Y {
				return true
			}
		}
	}
	return false
}

func (mh MoveHelper) isHeadToHeadCollision(coord Coord, snakes []Battlesnake) *Battlesnake {
	for _, snake := range snakes {
		if coord.X == snake.Head.X && coord.Y == snake.Head.Y {
			return &snake
		}
	}
	return nil
}

func (mh MoveHelper) simulateMove(state GameState, move string, lookAheadMoves int) GameState {
	simulatedState := state.deepcopy()
	newHead := mh.getNewHead(simulatedState.You.Body[0], move)
	simulatedState.You.Body = append([]Coord{newHead}, simulatedState.You.Body[:len(simulatedState.You.Body)-1]...)

	return simulatedState
}

func (mh MoveHelper) randomizeMoves(movesWithScores []moveScore) {
	rand.Shuffle(len(movesWithScores), func(i, j int) {
		movesWithScores[i], movesWithScores[j] = movesWithScores[j], movesWithScores[i]
	})
}

func (mh MoveHelper) sortMovesByScore(movesWithScores []moveScore) {
	sort.Slice(movesWithScores, func(i, j int) bool {
		return movesWithScores[i].Score > movesWithScores[j].Score
	})
}

func (mh MoveHelper) getAreaControlScore(state GameState, move string) int {
	score := 0
	newHead := mh.getNewHead(state.You.Body[0], move)
	adjacentCells := mh.getAdjacentCells(state, newHead)

	for _, cell := range adjacentCells {
		if mh.isCellEmpty(state, cell) {
			score++
		}
	}

	return score
}

func (mh MoveHelper) getAdjacentCells(state GameState, coord Coord) []Coord {
	adjacentCells := []Coord{}

	if coord.X > 0 {
		adjacentCells = append(adjacentCells, Coord{X: coord.X - 1, Y: coord.Y})
	}
	if coord.X < state.Board.Width-1 {
		adjacentCells = append(adjacentCells, Coord{X: coord.X + 1, Y: coord.Y})
	}
	if coord.Y > 0 {
		adjacentCells = append(adjacentCells, Coord{X: coord.X, Y: coord.Y - 1})
	}
	if coord.Y < state.Board.Height-1 {
		adjacentCells = append(adjacentCells, Coord{X: coord.X, Y: coord.Y + 1})
	}

	return adjacentCells
}

func (mh MoveHelper) isCellEmpty(state GameState, coord Coord) bool {
	for _, snake := range state.Board.Snakes {
		for _, snakePart := range snake.Body {
			if snakePart.X == coord.X && snakePart.Y == coord.Y {
				return false
			}
		}
	}

	for _, hazard := range state.Board.Hazards {
		if hazard.X == coord.X && hazard.Y == coord.Y {
			return false
		}
	}

	return true
}

// deepcopy is a helper function for GameState deepcopy
func (state GameState) deepcopy() GameState {
	snakesCopy := make([]Battlesnake, len(state.Board.Snakes))
	copy(snakesCopy, state.Board.Snakes)
	foodCopy := make([]Coord, len(state.Board.Food))
	copy(foodCopy, state.Board.Food)

	hazardsCopy := make([]Coord, len(state.Board.Hazards))
	copy(hazardsCopy, state.Board.Hazards)

	youBodyCopy := make([]Coord, len(state.You.Body))
	copy(youBodyCopy, state.You.Body)

	return GameState{
		Game:  state.Game,
		Turn:  state.Turn,
		Board: Board{Height: state.Board.Height, Width: state.Board.Width, Snakes: snakesCopy, Food: foodCopy, Hazards: hazardsCopy},
		You:   Battlesnake{ID: state.You.ID, Name: state.You.Name, Health: state.You.Health, Body: youBodyCopy},
	}
}
