package game

import (
	"net"
)

type Player struct {
	conn     net.Conn
    id int
	isWinner bool
}

func NewPlayer(conn net.Conn, name string, game *Game) *Player {
	return &Player{
		conn: conn,
		name: name,
		game: game,
	}
}

func (p *Player) Listen() {
	// code to listen for and handle player input
}

func (p *Player) Write(msg []byte) error {
	// code to write a message to the player's connection
}
