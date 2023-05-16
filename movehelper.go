package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/samyfodil/tb_library_snake_001/types"
)

type MoveHelper struct{}

func (mh MoveHelper) isOutOfBounds(point types.Coord, width int, height int) bool {
	return point.X < 0 || point.Y < 0 || point.X >= width || point.Y >= height
}

func (mh MoveHelper) isBodyCollision(point types.Coord, body []types.Coord) bool {
	for _, bodyPart := range body {
		if point.X == bodyPart.X && point.Y == bodyPart.Y {
			return true
		}
	}
	return false
}

func (mh MoveHelper) getAllowedMoves(state types.GameState) []string {
	allowedMoves := []string{"up", "down", "left", "right"}

	head := state.You.Body[0]
	width := state.Board.Width
	height := state.Board.Height

	for _, move := range allowedMoves {
		newHead := mh.getNewHead(head, move)
		if mh.isOutOfBounds(newHead, width, height) || mh.isBodyCollision(newHead, state.You.Body) {
			allowedMoves = mh.remove(allowedMoves, move)
		}
	}

	return allowedMoves
}

func (mh MoveHelper) getNewHead(head types.Coord, move string) types.Coord {
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

func (mh MoveHelper) remove(slice []string, s string) []string {
	index := -1
	for i, v := range slice {
		if v == s {
			index = i
			break
		}
	}

	if index >= 0 {
		return append(slice[:index], slice[index+1:]...)
	}

	return slice
}

func (mh MoveHelper) euclideanDistance(a types.Coord, b types.Coord) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func (mh MoveHelper) randFloat() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

func (mh MoveHelper) isMoveSafe(state types.GameState, move string) bool {
	newHead := mh.getNewHead(state.You.Body[0], move)
	return mh.isInBounds(newHead, state.Board) && !mh.isCollidingWithSelf(newHead, state.You.Body)
}

func (mh MoveHelper) getCollisionScore(state types.GameState, move string) int {
	newHead := mh.getNewHead(state.You.Body[0], move)
	return mh.distanceToClosestCollision(newHead, state.Board)
}

func (mh MoveHelper) isInBounds(coord types.Coord, board types.Board) bool {
	return coord.X >= 0 && coord.X < board.Width && coord.Y >= 0 && coord.Y < board.Height
}

func (mh MoveHelper) isCollidingWithSelf(newHead types.Coord, body []types.Coord) bool {
	for _, bodyPart := range body {
		if newHead.X == bodyPart.X && newHead.Y == bodyPart.Y {
			return true
		}
	}
	return false
}

func (mh MoveHelper) updateGameStateForMove(board types.Board, you types.Battlesnake, newHead types.Coord) (types.Board, types.Battlesnake) {
	updatedBoard := board
	updatedSnake := you

	// Check for food consumption
	consumedFood := false
	for i, food := range board.Food {
		if newHead == food {
			consumedFood = true
			updatedBoard.Food = append(updatedBoard.Food[:i], updatedBoard.Food[i+1:]...)
			break
		}
	}

	// Update the snake's body
	if consumedFood {
		updatedSnake.Body = append([]types.Coord{newHead}, updatedSnake.Body...)
	} else {
		updatedSnake.Body = append([]types.Coord{newHead}, updatedSnake.Body[:len(updatedSnake.Body)-1]...)
	}

	// Update the snake's health
	if consumedFood {
		updatedSnake.Health = 100
	} else {
		updatedSnake.Health -= 1
	}

	return updatedBoard, updatedSnake
}

func (mh MoveHelper) distanceToClosestFood(coord types.Coord, food []types.Coord) int {
	minDist := math.MaxInt32
	for _, foodCoord := range food {
		dist := mh.manhattanDistance(coord, foodCoord)
		if dist < minDist {
			minDist = dist
		}
	}
	return minDist
}

func (mh MoveHelper) distanceToClosestCollision(coord types.Coord, board types.Board) int {
	minDist := math.MaxInt32
	for _, snake := range board.Snakes {
		for _, bodyCoord := range snake.Body {
			dist := mh.manhattanDistance(coord, bodyCoord)
			if dist < minDist {
				minDist = dist
			}
		}
	}
	return minDist
}

func (mh MoveHelper) DeepCopy(state types.GameState) types.GameState {
	copiedState := types.GameState{}
	copiedState.Turn = state.Turn
	copiedState.You = state.You
	copiedState.Board.Width = state.Board.Width
	copiedState.Board.Height = state.Board.Height
	copiedState.Board.Food = make([]types.Coord, len(state.Board.Food))
	copy(copiedState.Board.Food, state.Board.Food)
	copiedState.Board.Hazards = make([]types.Coord, len(state.Board.Hazards))
	copy(copiedState.Board.Hazards, state.Board.Hazards)
	copiedState.Board.Snakes = make([]types.Battlesnake, len(state.Board.Snakes))
	for i, snake := range state.Board.Snakes {
		copiedState.Board.Snakes[i] = types.Battlesnake{
			ID:     snake.ID,
			Name:   snake.Name,
			Health: snake.Health,
			Body:   make([]types.Coord, len(snake.Body)),
		}
		copy(copiedState.Board.Snakes[i].Body, snake.Body)
	}
	return copiedState
}

func (mh MoveHelper) manhattanDistance(a, b types.Coord) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}
