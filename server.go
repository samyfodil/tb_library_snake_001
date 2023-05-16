package lib

import (
	"function/types"
	v1 "function/v1"
	v2 "function/v2"
	"io"

	"github.com/taubyte/go-sdk/event"
)

const ServerID = "github.com/samyfodil/tb_library_snake_001"

// HTTP Handlers

//export index
func index(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	h.Headers().Set("Server", ServerID)
	h.Headers().Set("Content-Type", "application/json")

	response := info()

	data, err := response.MarshalJSON()
	if err != nil {
		h.Write([]byte(err.Error()))
		return 1
	}

	h.Write(data)

	return 0
}

//export start
func start(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	h.Headers().Set("Server", ServerID)

	state := types.GameState{}

	data, err := io.ReadAll(h.Body())
	if err != nil {
		h.Write([]byte(err.Error()))
		return 1
	}

	err = state.UnmarshalJSON(data)
	if err != nil {
		h.Write([]byte(err.Error()))
		return 1
	}

	return 0
}

//export move
func move(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	h.Headers().Set("Server", ServerID)

	state := types.GameState{}

	data, err := io.ReadAll(h.Body())
	if err != nil {
		h.Write([]byte(err.Error()))
		return 1
	}

	err = state.UnmarshalJSON(data)
	if err != nil {
		h.Write([]byte(err.Error()))
		return 1
	}

	response := types.BattlesnakeMoveResponse{
		Move: "down",
	}

	switch state.You.Name {
	case "tau001":
		response = v1.Domove(state)
	case "tau002":
		response = v1.Domove2(state)
	case "tau003":
		response = v1.Domove3(state)
	case "tau004":
		response = v1.Domove4(state)
	case "tau005":
		response = v1.Domove5(state)
	case "tau006":
		response = v2.Move(state)
	}

	h.Headers().Set("Content-Type", "application/json")

	data, err = response.MarshalJSON()
	if err != nil {
		h.Write([]byte(err.Error()))
		return 1
	}

	h.Write(data)

	h.Return(200)

	return 0
}

//export end
func end(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	h.Headers().Set("Server", ServerID)

	state := types.GameState{}

	data, err := io.ReadAll(h.Body())
	if err != nil {
		h.Write([]byte(err.Error()))
		return 1
	}

	err = state.UnmarshalJSON(data)
	if err != nil {
		h.Write([]byte(err.Error()))
		return 1
	}

	return 0
}
