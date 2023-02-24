package game

import (
	"fmt"
	"net"

	"github.com/google/uuid"
)

type Player struct {
	Conn  net.Conn
	ID    uuid.UUID
	Score int
	Snake *Snake
}

func NewPlayer(conn net.Conn) *Player {
	s := NewSnake(1, 1) // get info someware

	return &Player{
		Conn:  conn,
		ID:    uuid.New(),
		Score: 0,
		Snake: s,
	}
}

func (p *Player) String() string {
	return fmt.Sprintf("Player %s (score: %d, snake: %s)", p.ID.String(), p.Score, p.Snake.String())
}

// func (p *Player) Listen() {
// 	// code to listen for and handle player input
// }
//
// func (p *Player) Write(msg []byte) error {
// 	// code to write a message to the player's connection
// 	return nil
// }
