package game

import (
	"bytes"
	"sync"
	"time"

	"github.com/google/uuid"
)

type (
	Direction int
	cell      int
)

const (
	Up Direction = iota
	Down
	Left
	Right
)

const (
	void cell = iota
	wall
	fruit
	snake
)

type Game struct {
	level         *level
	Players       []*Player
	InputChan     chan InputMessage
	mu            sync.Mutex
	GameStateChan chan *[]byte
	running       bool
}

type InputMessage struct {
	PlayerID uuid.UUID
	Input    Direction
}

type level struct {
	height int
	width  int
	data   [][]cell
}

func NewGame(width, height int) *Game {
	return &Game{
		Players:       []*Player{},
		InputChan:     make(chan InputMessage),
		mu:            sync.Mutex{},
		GameStateChan: make(chan *[]byte, 1),
		level:         newLevel(width, height),
		running:       true,
	}
}

func newLevel(width, height int) *level {
	data := make([][]cell, height)

	// Initialize each element of the slice to empty
	for h := 0; h < height; h++ {
		data[h] = make([]cell, width)
		for w := 0; w < width; w++ {
			data[h][w] = void
		}
	}

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if h == 0 {
				data[h][w] = wall
			}
			if h == height-1 {
				data[h][w] = wall
			}
			if w == 0 {
				data[h][w] = wall
			}
			if w == width-1 {
				data[h][w] = wall
			}
		}
	}

	return &level{
		height: height,
		width:  width,
		data:   data,
	}
}

func (g *Game) gameLoop() {
	for {
		go g.updatePlayerSnakes()
		g.mu.Lock()
		gameState := g.SerializeGameState()
		g.GameStateChan <- gameState
		g.mu.Unlock()

		time.Sleep(time.Millisecond * 200)
	}
}

func (g *Game) SerializeGameState() *[]byte {
	state := g.renderLevel()
	// gameState := &protobuf.Message{
	// 	Type: &protobuf.Message_Game{
	// 		Game: state.String(),
	// 	},
	// }
	return &state
}

func (g *Game) Start() {
	if g.running {
		go g.gameLoop()
		go g.handleInput()
	}
}

func (g *Game) Stop() {
	g.running = false
}

var telnetControlChars = map[string]byte{
	"InterpretAsCommand": 255,
	"Will":               251,
	"Wont":               252,
	"Do":                 253,
	"Dont":               254,
	"SubnegotiationEnd":  240,
	"CarriageReturn":     13,
	"LineFeed":           10,
	"ClearScreen":        27,
}

func (g *Game) renderLevel() []byte {
	var buff bytes.Buffer
	for h := 0; h < g.level.height; h++ {
		for w := 0; w < g.level.width; w++ {
			occupied := false
			for _, p := range g.Players {
				if p.Snake.Occupies(h, w) {
					buff.WriteByte('S')
					occupied = true
					break
				}
			}
			if !occupied {
				switch g.level.data[h][w] {
				case wall:
					buff.WriteByte('X')
				case fruit:
					buff.WriteByte('X')
				case void:
					buff.WriteByte(' ')
				}
			}
		}
		buff.Write([]byte{telnetControlChars["CarriageReturn"], telnetControlChars["LineFeed"]})
	}

	return buff.Bytes()
}

func (g *Game) RemovePlayer(playerID uuid.UUID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i, player := range g.Players {
		if player.ID == playerID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			break
		}
	}
}

func (g *Game) handleInput() {
	for input := range g.InputChan {
		for _, player := range g.Players {
			if player.ID == input.PlayerID {
				player.Snake.SetDirection(input.Input)
			}
		}
	}
}

func (g *Game) updatePlayerSnakes() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, p := range g.Players {
		if p.Snake != nil {
			p.Snake.Move()
		}
	}
}
