package types

// Position of Column and Row
type Position struct {
	Column byte `json:"column" yaml:"column"`
	Row    byte `json:"row" yaml:"row"`
}

// RobotColor defines an Enum of the colors in the following order: yellow, red, green, blue
type RobotColor byte

// RobotOrder represents the order of a board state by lining up all arranged robot colors into one byte
type RobotOrder byte

const (
	RobotColorYellow RobotColor = iota // RobotColorYellow Transforms RobotColor into an integer enum - Yellow = 0
	RobotColorRed    RobotColor = iota // RobotColorRed Transforms RobotColor into an integer enum - Red = 1
	RobotColorGreen  RobotColor = iota // RobotColorGreen Transforms RobotColor into an integer enum - Green = 2
	RobotColorBlue   RobotColor = iota // RobotColorBlue Transforms RobotColor into an integer enum - Blue = 3
)

// BoardState represents a state of the board in which only the robots have been moved.
type BoardState uint32

// Node represents the attributes of each board cell. It is divided into the minimal move count (bit 7-5) and the neighbors (bit 3-0).
// The minimal move count is the number of moves that a robot has to make, if it could stop everywhere, to get to the target. This value is saved as a 3bit uint and has an upper bound of 7.
// To save the configuration of the cell neighbors each direction has one designated bit in the following order: north, south, west, east. If the bit is set to 0 than the robot has no neighbor in that direction.
type Node byte

// Robot represents the position of the robot. It is divided into the column (bit 7-4) and the row (bit 3-0).
type Robot byte

// Target represents the color (bit 15-12), symbol (bit 11-8) and position (7-0) of the target.
// The color of the target is defined by the high bit in the color bit region. For that every color has one bit in the following order: yellow, red, green, blue.
// The symbol is defined in the same way with the order: circle, triangle, square, hexagon.
// The position is divided into the column (bit 7-4) and the row (bit 3-0).
type Target uint16

// Board represents the current board configuration to solve
type Board struct {
	Grid   [16][16]Node
	Target Target
}

// RobotStoppingPositions represents a matrix of robot stopping position. These positions are calculated by the assumption, that every robot could stop wherever it wants to.
// The four calculated positions are packed into one uint32 in the following order: up, down, left, right (bit 31 -> 0). Each position consist of a column (bit: 7-4) and a row (bit: 3-0).
type RobotStoppingPositions [16][16]uint32
