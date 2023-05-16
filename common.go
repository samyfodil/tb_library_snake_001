package lib

import (
	"function/types"
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
