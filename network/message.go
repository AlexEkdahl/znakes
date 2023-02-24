package network

import "github.com/golang/protobuf/proto"

// MessageType represents the type of a message
type MessageType int32

const (
	// Welcome represents a welcome message
	Welcome MessageType = 0
	// Join represents a join message
	Join MessageType = 1
	// GameState represents a game state message
	GameState MessageType = 2
	// Input represents an input message
	Input MessageType = 3
	// GameOver represents a game over message
	GameOver MessageType = 4
)

// Message represents a message
type Message struct {
	Type             MessageType
	WelcomeMessage   *string
	PlayerNum        int32
	GameMode         GameMode
	PlayerMode       PlayerMode
	GameStateMessage *GameStateMessage
	InputMessage     *InputMessage
	GameOverMessage  *GameOverMessage
}

// InputMessage represents an input message
type InputMessage struct {
	Input      string
	PlayerNum  int32
	GameMode   GameMode
	PlayerMode PlayerMode
}

// GameStateMessage represents the game state in a message
type GameStateMessage struct {
	Game *Game
}

// GameOverMessage represents a game over message
type GameOverMessage struct {
	Message string
}

// Encode encodes a message as a protobuf message
func EncodeMessage(m *Message) ([]byte, error) {
	return proto.Marshal(m)
}

// Decode decodes a protobuf message into a message
func DecodeMessage(data []byte, m *Message) error {
	return proto.Unmarshal(data, m)
}
