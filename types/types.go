package types

type Position struct {
	Column byte `json:"column"`
	Row    byte `json:"row"`
}

type Neighbors struct {
	Up    bool `json:"top"`
	Down  bool `json:"bottom"`
	Left  bool `json:"left"`
	Right bool `json:"right"`
}

type Node struct {
	Neighbors Neighbors
	Position  Position
}

type Robot struct {
	Color    string `json:"color"`
	Position Position
}

type Colors uint8

const (
	Yellow Colors = iota
	Red           = iota
	Green         = iota
	Blue          = iota
)

type RawTarget struct {
	Color    string `json:"color"`
	Symbol   string `json:"symbol"`
	Position Position
}

type RawBoard struct {
	Nodes  [][]Node  `json:"nodes"`
	Robots []Robot   `json:"robots"`
	Target RawTarget `json:"target"`
}

type BoardState uint32

type Board struct {
	Board  [16][16]byte
	Robots [4]uint8
	Target uint16
}
