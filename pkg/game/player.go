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
	s := NewSnake(5, 5) // get info someware

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
