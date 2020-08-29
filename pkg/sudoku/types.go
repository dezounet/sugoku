package sudoku

// Difficulty specific type
type Difficulty int

// Level of difficulty
const (
	EASY      Difficulty = iota // 0
	MEDIUM                      // 1
	HARD                        // 2
	NIGHTMARE                   // 3
)

const empty = 0

// Coordinates as row / column
type Coordinates struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

// Cell coordinates and value
type Cell struct {
	Coordinates `mapstructure:",squash"`
	Value       int  `json:"value,omitempty"`
	Frozen      bool `json:"frozen,omitempty"`
}

// Grid made of several cells
type Grid struct {
	UUID      string   `json:"uuid"`
	BlockSize int      `json:"blocksize"`
	Cells     [][]Cell `json:"cells"`
}
