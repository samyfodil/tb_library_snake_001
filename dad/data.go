package dad

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/samyfodil/tb_library_snake_001/types"
)

type MMArray [][]MinMaxData

func (ret MMArray) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("\n")
	for i := range ret {
		for j := range ret[i] {
			if ret[i][j].tie {
				// there are no artuculation points yet
				buffer.WriteString(" ")
				for _, sid := range ret[i][j].snakeIds {
					buffer.WriteString(fmt.Sprintf("%d", sid))
				}
				buffer.WriteString(" ")
			} else if len(ret[i][j].snakeIds) == 0 {
				buffer.WriteString(" XX ")
			} else {
				//buffer.WriteString(fmt.Sprintf("  %d ", ret[i][j].moves))
				buffer.WriteString(fmt.Sprintf("  %d ", ret[i][j].snakeIds[0]))
			}
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// MetaData
// contains any computed data about the move request.
// is used in composition with the move request so you cannot
// have name colisions with the MoveRequest struct
type MetaData struct {
	// denotes the number of moves until you reach the closest piece of food
	MyLength   int
	MyIndex    int
	Hazards    map[string]bool
	FoodMap    map[string]bool
	tightSpace bool
	DistToFood int
	// making this a pointer makes it able to be tested against
	// nil so we might as well keep it like this
	SnakeHash  map[string]*SnakeData
	KillZones  map[string]bool
	SnakeHeads map[string]bool
	minMaxArr  MMArray
	MinMaxMD   MinMaxMetaData
	KSD        KeySnakeData
}

// MetaDataDirec
// contains any computed data in a particular direction
// is used in composition with the move request so you cannot
// have name colisions with the MoveRequest struct
type MetaDataDirec struct {
	// denotes the number of moves until you reach the closest piece of food
	// indexed by their point
	// totals up your length and the ammount of food in a direction
	// if you would fill up the space make it unlikely to go that direction
	MovesVsSpace int
	// the total number of moves possible in this direction
	// contains a map to the last accessable piece of a snake
	// from your current location if you moved in this direction
	minMaxArr    MMArray
	ClosestFood  *Point
	Food         int
	Moves        int
	SeeTail      bool
	KeySnakeData KeySnakeData
	// indexed by their point
	FoodHash   map[string]*FoodData
	sortedFood []*FoodData
	MoveHash   map[string]*MinMaxData
	MinMaxMD   MinMaxMetaData
}

type KeySnakeData map[int]*SnakeData

// minKeySnakePart
// returns the snake data for the point you are waiting to open up
// it is the least number of moves that anyone around you can make before
// you are able to exit the area you are in
func (ksd KeySnakeData) minKeySnakePart() *SnakeData {
	var min *SnakeData
	for _, val := range ksd {
		if min == nil || min.lengthLeft > val.lengthLeft {
			min = val
		}
	}
	return min
}

func (m *MoveRequest) NoFood() bool {
	for _, val := range m.Direcs {
		if val.Food > 0 {
			return false
		}
	}
	return true
}

// String
// used to print the metadata for a particular direction
// it is necessary because the Static data is a pointer
// unfortunately this means that you have to manually manage this
// maybe I could make
func (m *MetaDataDirec) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("MetaDataDirec{")
	buffer.WriteString(fmt.Sprintf("movesVsSpace:%v, ", m.MovesVsSpace))
	buffer.WriteString(fmt.Sprintf("TotalMoves:%v, ", m.Moves))
	buffer.WriteString(fmt.Sprintf("KeySnakeData:\n"))
	for direc, val := range m.KeySnakeData {
		buffer.WriteString(fmt.Sprintf("\t%v:%v", direc, val))
	}
	buffer.WriteString("}\n")

	return buffer.String()
}

// used to find and set the length of your snake globally in the
// metatdata object
func (m *MetaData) SetMyLength(data *MoveRequest) {
	for i, snake := range data.Snakes {
		if snake.Id == data.You && len(data.You) > 0 {
			m.MyLength = len(snake.Coords)
			swap(data.Snakes, i, len(data.Snakes)-1)
			m.MyIndex = len(data.Snakes) - 1
			//m.MyIndex = i
		}
	}
}

// a little struct used to see the length left after this portion of a
// snakes body the tail of the snake has a value of 1
type MinMaxData struct {
	moves    int
	snakeIds []int
	snakes   map[int]bool
	tie      bool
	pnt      *Point
}

type MinMaxSnakeMD struct {
	moves int
	ties  int
}
type MinMaxMetaData struct {
	movesHash map[string]int
	tiesHash  map[string][]int
	snakes    map[int]MinMaxSnakeMD
}

// a little struct used to see the length left after this portion of a
// snakes body the tail of the snake has a value of 1
type SnakeData struct {
	id         int
	lengthLeft int
	pnt        *Point
}

func (s *SnakeData) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString(fmt.Sprintf("id:%v, ", s.id))
	buffer.WriteString(fmt.Sprintf("lengthLeft:%v, ", s.lengthLeft))
	buffer.WriteString(fmt.Sprintf("pnt:%v, ", s.pnt))
	buffer.WriteString("}\n")

	return buffer.String()
}

func (m *MetaData) GenTailData(data *MoveRequest) {
	for i, snake := range data.Snakes {
		tail, _ := getTail(i, data)
		head := snake.Coords[0]

		m.Hazards[tail.String()] = false
		d := head.Down(data)
		if d != nil && m.FoodMap[d.String()] {
			m.Hazards[tail.String()] = true
		}
		d = head.Up(data)
		if d != nil && m.FoodMap[d.String()] {
			m.Hazards[tail.String()] = true
		}
		d = head.Right(data)
		if d != nil && m.FoodMap[d.String()] {
			m.Hazards[tail.String()] = true
		}
		d = head.Left(data)
		if d != nil && m.FoodMap[d.String()] {
			m.Hazards[tail.String()] = true
		}

	}
}

// GenenSnakeHash
//
//		generates a map of all the points in all the snakes
//		is used to determine how much of the snake must move
//	     in order for the area they are blocking to be open
func (m *MetaData) GenSnakeHash(data *MoveRequest) {
	m.SnakeHash = make(map[string]*SnakeData)
	for i, snake := range data.Snakes {
		for j, coord := range snake.Coords {
			m.SnakeHash[coord.String()] = &SnakeData{
				id:         i,
				lengthLeft: len(snake.Coords) - j - 1,
				pnt:        &Point{coord.X, coord.Y},
			}
		}
	}
}

func (m *MetaData) GenMinMax(data *MoveRequest) {
	defer data.GenHazards(data, true)

	tightSpace := true
	MinMax(data, "")
	for direc, direcData := range data.Direcs {
		MinMax(data, direc)
		//fmt.Printf(direc)
		ksd := direcData.KeySnakeData.minKeySnakePart()

		if ksd != nil {
			direcData.MovesVsSpace = direcData.Moves - direcData.Food - ksd.lengthLeft
		}
		if direcData.MovesVsSpace > 20 {
			tightSpace = false
		}
	}
	data.tightSpace = tightSpace
}

// Generates a map of hazards
// this is pretty janky and will need to get refactored
func (m *MetaData) GenHazards(data *MoveRequest, snakeMovesAsHazards bool) {
	m.Hazards = make(map[string]bool)
	m.KillZones = make(map[string]bool)
	m.GenTailData(data)
	for _, snake := range data.Snakes {
		snake.HeadPoint = &(snake.Coords[0])

		if len(snake.Coords) >= m.MyLength && data.You != snake.Id && snakeMovesAsHazards {
			head := snake.Head()
			d := head.Down(data)
			if d != nil {
				m.Hazards[d.String()] = true
			}
			d = head.Up(data)
			if d != nil {
				m.Hazards[d.String()] = true
			}
			d = head.Right(data)
			if d != nil {
				m.Hazards[d.String()] = true
			}
			d = head.Left(data)
			if d != nil {
				m.Hazards[d.String()] = true
			}

		} else if len(snake.Coords) < m.MyLength && data.You != snake.Id && snakeMovesAsHazards {
			head := snake.Head()
			d := head.Down(data)
			if d != nil {
				m.KillZones[d.String()] = true
			}
			d = head.Up(data)
			if d != nil {
				m.KillZones[d.String()] = true
			}
			d = head.Right(data)
			if d != nil {
				m.KillZones[d.String()] = true
			}
			d = head.Left(data)
			if d != nil {
				m.KillZones[d.String()] = true
			}

		}
		for i, coord := range snake.Coords {
			if i == len(snake.Coords)-1 {
				break
			}
			m.Hazards[coord.String()] = true
		}

	}
}

// generates a map of all the food points
func (m *MetaData) GenFoodMap(data *MoveRequest) {
	m.FoodMap = make(map[string]bool)
	for _, food := range data.Food {
		m.FoodMap[food.String()] = true
	}
}

// alias for the metadata map
type MoveMetaData map[string]*MetaDataDirec

type FoodData struct {
	moves   int
	pnt     *Point
	closest int
}

// RESPONSE AND REQUEST STRUCTS
type MoveResponse struct {
	Move  string `json:"move"`
	Taunt string `json:"taunt,omitempty"`
}

type GameStartRequest struct {
	GameId string `json:"game_id"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

func NewGameStartRequest(req *http.Request) (*GameStartRequest, error) {
	decoded := GameStartRequest{}
	err := json.NewDecoder(req.Body).Decode(&decoded)
	return &decoded, err
}

type GameStartResponse struct {
	Color    string  `json:"color"`
	HeadUrl  *string `json:"head_url,omitempty"`
	Name     string  `json:"name"`
	Taunt    *string `json:"taunt,omitempty"`
	HeadType string  `json:"head_type"`
	TailType string  `json:"tail_type"`
}

type MoveRequest struct {
	// static
	GameId string   `json:"game_id"`
	Height int      `json:"height"`
	Width  int      `json:"width"`
	Turn   int      `json:"turn"`
	Food   []Point  `json:"food"`
	Snakes []*Snake `json:"snakes"`
	You    string   `json:"you"`

	// added here for convenience
	MetaData
	Direcs MoveMetaData
}

func (m *MoveRequest) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("\nMoveRequest{\n")
	head, _ := getMyHead(m)
	buffer.WriteString(fmt.Sprintf("head: %v ", head))
	buffer.WriteString("Direcs{\n")
	for direc, val := range m.Direcs {
		buffer.WriteString(fmt.Sprintf("\t%v:%v", direc, val))
	}
	buffer.WriteString(fmt.Sprintf("tightSpace: %v ", m.tightSpace))
	buffer.WriteString(fmt.Sprintf("MyIndex: %v ", m.MyIndex))
	buffer.WriteString(fmt.Sprintf("MyIndex: %v ", m.KillZones))
	buffer.WriteString("}\n")
	buffer.WriteString("}\n")

	return buffer.String()
}

// initializes global meta data attributes
func (m *MoveRequest) init() {
	m.Direcs = make(MoveMetaData)
	m.Direcs[UP] = &MetaDataDirec{}
	m.Direcs[DOWN] = &MetaDataDirec{}
	m.Direcs[LEFT] = &MetaDataDirec{}
	m.Direcs[RIGHT] = &MetaDataDirec{}

	m.SetMyLength(m)
	m.GenFoodMap(m)
	m.GenSnakeHash(m)
	m.GenMinMax(m)
	m.GenHazards(m, true)
	m.MinMaxMD = GenMinMaxStats(m.minMaxArr)
	for _, val := range m.Direcs {
		val.MinMaxMD = GenMinMaxStats(val.minMaxArr)
	}
}

// de serializes the move request data into a string
func getMoveRequestString(req *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	return buf.String()
}

func battlesnakeToSnake(bs types.Battlesnake) *Snake {
	coords := make([]Point, len(bs.Body))
	for i, coord := range bs.Body {
		coords[i] = Point{X: coord.X, Y: coord.Y}
	}

	return &Snake{
		Coords:       coords,
		HealthPoints: bs.Health,
		Id:           bs.ID,
		Name:         bs.Name,
		Taunt:        bs.Shout, // You might need to adjust this field depending on the original struct
	}
}

func gameStateToMoveRequest(state *types.GameState) *MoveRequest {
	foodPoints := make([]Point, len(state.Board.Food))
	for i, food := range state.Board.Food {
		foodPoints[i] = Point{X: food.X, Y: food.Y}
	}

	snakes := make([]*Snake, len(state.Board.Snakes))
	for i, bs := range state.Board.Snakes {
		snakes[i] = battlesnakeToSnake(bs)
	}

	return &MoveRequest{
		GameId: state.Game.ID,
		Height: state.Board.Height,
		Width:  state.Board.Width,
		Turn:   state.Turn,
		Food:   foodPoints,
		Snakes: snakes,
		You:    state.You.ID,
	}
}

// // creates a new move request
// func NewMoveRequest(state *types.GameState) (*MoveRequest, error) {
// 	res := new(MoveRequest)

// 	err := json.Unmarshal([]byte(str), res)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res.init()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }

// Decode [number, number] JSON array into a Point

// Allows decoding a string or number identifier in JSON
// by removing any surrounding quotes and storing in a string
type Snake struct {
	Coords       []Point `json:"coords"`
	HealthPoints int     `json:"health_points"`
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Taunt        string  `json:"taunt"`
	HeadStack    *Stack
	TailStack    *Stack
	HeadPoint    *Point
	FoodHash     map[string]*FoodData
}

func (snake Snake) Head() *Point    { return snake.HeadPoint }
func (snake *Snake) String() string { return fmt.Sprintf("%#v", snake) }
