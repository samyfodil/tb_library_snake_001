package main

import (
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
