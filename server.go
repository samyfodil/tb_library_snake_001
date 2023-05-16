package lib

import (
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

	state := GameState{}

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

	state := GameState{}

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

	var response BattlesnakeMoveResponse

	switch state.You.Name {
	case "tau001":
		response = domove(state)
	case "tau002":
		response = domove2(state)
	case "tau003":
		response = domove3(state)
	case "tau004":
		response = domove4(state)
	case "tau005":
		response = domove5(state)
	case "tau006":
		response = MoveStrategyX{NewMoveHelper(8)}.GetMove(state)
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

	state := GameState{}

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
