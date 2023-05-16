package lib

import (
	"math"
	"sort"
)

const (
	DefaultLookAheadTurns = 3
	DefaultSimulateStuck  = 0.7 // 10% chance of getting stuck
)

type MoveStrategyY struct {
	MoveHelper
	LookAheadTurns int
	SimulateStuck  float64
}

func NewMoveStrategyY() MoveStrategyY {
	return MoveStrategyY{
		MoveHelper:     MoveHelper{},
		LookAheadTurns: DefaultLookAheadTurns,
		SimulateStuck:  DefaultSimulateStuck,
	}
}

func (ms MoveStrategyY) GetMove(state GameState) BattlesnakeMoveResponse {
	safeMoves := []string{}

	for _, move := range ms.getAllowedMoves(state) {
		if ms.isMoveSafe(state, move) {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		// No safe moves detected, loop over self
		safeMoves = ms.loopOverSelf(state)
	}

	// Calculate scores for each move based on the requirements
	moveScores := make(map[string]float64)
	for _, move := range safeMoves {
		moveScores[move] = ms.calculateMoveScoreY(state, move)
	}

	// Sort the safe moves based on their scores
	sort.SliceStable(safeMoves, func(i, j int) bool {
		return moveScores[safeMoves[i]] > moveScores[safeMoves[j]]
	})

	// Choose the move with the highest score
	chosenMove := safeMoves[0]

	return BattlesnakeMoveResponse{Move: chosenMove}
}

func (ms MoveStrategyY) calculateMoveScoreY(state GameState, move string) float64 {
	newState := ms.simulateMove(state, move, ms.LookAheadTurns, ms.SimulateStuck)
	score := 0.0

	// Check if the move is out of bounds
	if ms.isOutOfBounds(newState.Board, newState.You.Body[0]) {
		return math.Inf(-1) // return a negative infinity score to strongly penalize the move
	}

	// Encourage pushing other snakes to hit bounds
	score += ms.getBoundaryPushScore(state)

	// Discourage getting food unless necessary
	if state.You.Health < 60 {
		score += float64(ms.getFoodScore(state, newState))
	}

	// Avoid hitting other snakes
	score -= ms.getCollisionScore(newState)

	return score
}

func (ms MoveStrategyY) isCollidingWithSnake(head Coord, snakeBody []Coord) bool {
	for _, bodyPart := range snakeBody {
		if head.X == bodyPart.X && head.Y == bodyPart.Y {
			return true
		}
	}
	return false
}

func (ms MoveStrategyY) getCollisionScore(newState GameState) float64 {
	collisionScore := 0.0

	for _, snake := range newState.Board.Snakes {
		if snake.ID == newState.You.ID {
			continue
		}

		if ms.isCollidingWithSnake(newState.You.Body[0], snake.Body) {
			collisionScore += 100.0
		}
	}

	return collisionScore
}

// func (ms MoveStrategyY) getCollisionScore(newState GameState) float64 {
// 	collisionScore := 0.0

// 	for _, snake := range newState.Board.Snakes {
// 		if snake.ID == newState.You.ID {
// 			continue
// 		}

// 		if ms.isCollidingWithSnake(newState.You.Body[0], snake.Body) {
// 			collisionScore += 100.0
// 		}
// 	}

// 	return collisionScore
// }

func (ms MoveStrategyY) getFoodScore(oldState, newState GameState) int {
	oldDist := ms.distanceToClosestFood(oldState.You.Body[0], oldState.Board.Food)
	newDist := ms.distanceToClosestFood(newState.You.Body[0], newState.Board.Food)

	if newDist < oldDist {
		return -1
	}

	return 0
}

func (ms MoveStrategyY) loopOverSelf(state GameState) []string {
	loopMoves := []string{}

	for _, move := range ms.getAllowedMoves(state) {
		newHead := ms.getNewHead(state.You.Body[0], move)
		if ms.isInBounds(newHead, state.Board) && !ms.isCollidingWithSelf(newHead, state.You.Body) {
			loopMoves = append(loopMoves, move)
		}
	}

	return loopMoves
}

func (ms MoveStrategyY) simulateMove(state GameState, move string, lookAheadTurns int, stuckChance float64) GameState {
	newState := ms.DeepCopy(state)
	newHead := ms.getNewHead(state.You.Body[0], move)

	// Simulate stuck scenario (no response in time)
	if ms.randFloat() < stuckChance {
		move = ms.getPreviousMove(state)
		newHead = ms.getNewHead(state.You.Body[0], move)
	}

	newState.Board, newState.You = ms.updateGameStateForMove(newState.Board, newState.You, newHead)

	if lookAheadTurns > 1 {
		nextMove := ms.GetMove(newState).Move
		newState = ms.simulateMove(newState, nextMove, lookAheadTurns-1, stuckChance)
	}

	return newState
}

func (ms MoveStrategyY) getPreviousMove(state GameState) string {
	head := state.You.Body[0]
	neck := state.You.Body[1]
	if neck.X < head.X {
		return "left"
	} else if neck.X > head.X {
		return "right"
	} else if neck.Y < head.Y {
		return "down"
	} else if neck.Y > head.Y {
		return "up"
	}

	return "up"
}

func (ms MoveStrategyY) getBoundaryPushScore(state GameState) float64 {
	score := 0.0
	for _, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			continue
		}

		head := snake.Body[0]
		score += math.Min(float64(head.X), float64(state.Board.Width-head.X-1))
		score += math.Min(float64(head.Y), float64(state.Board.Height-head.Y-1))
	}

	return score
}
