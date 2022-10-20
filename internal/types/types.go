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

type Wall struct {
	Position1  Position `json:"position1"`
	Position2  Position `json:"position2"`
	Direction1 string   `json:"direction1"`
	Direction2 string   `json:"direction2"`
}

type Robot struct {
	Color    string `json:"color"`
	Position Position
}

type RobotColor byte
type RobotOrder byte

const (
	RobotColorYellow RobotColor = iota
	RobotColorRed    RobotColor = iota
	RobotColorGreen  RobotColor = iota
	RobotColorBlue   RobotColor = iota
)

type RawTarget struct {
	Color    string `json:"color"`
	Symbol   string `json:"symbol"`
	Position Position
}

type RawBoard struct {
	Walls  []Wall    `json:"walls"`
	Robots []Robot   `json:"robots"`
	Target RawTarget `json:"target"`
}

type BoardState uint32

type Board struct {
	Grid   [16][16]byte
	Target uint16
}

type RobotStoppingPositions [16][16]uint32