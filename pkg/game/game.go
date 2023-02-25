package game

import (
	"bytes"
	"sync"
	"time"

	"github.com/AlexEkdahl/snakes/pkg/network/protobuf"
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
	food          []*Food
	mu            sync.Mutex
	GameStateChan chan *protobuf.Message
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

type Food struct {
	x, y int
}

func NewGame(width, height int) *Game {
	return &Game{
		Players:       []*Player{},
		InputChan:     make(chan InputMessage),
		food:          []*Food{},
		mu:            sync.Mutex{},
		GameStateChan: make(chan *protobuf.Message, 1),
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

func (g *Game) SerializeGameState() *protobuf.Message {
	state := g.renderLevel()
	gameState := &protobuf.Message{
		Type: &protobuf.Message_Game{
			Game: state.String(),
		},
	}
	return gameState
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

func (g *Game) renderLevel() bytes.Buffer {
	var buff bytes.Buffer

	for h := 0; h < g.level.height; h++ {
		for w := 0; w < g.level.width; w++ {
			for _, p := range g.Players {
				if p.Snake.Occupies(h, w) {
					buff.WriteString("S")
				}
			}
			switch g.level.data[h][w] {
			case wall:
				buff.WriteString("X")
			case fruit:
				buff.WriteString("F")
			case void:
				buff.WriteString(" ")
			}
		}
		buff.WriteString("\n")
	}

	for _, p := range g.Players {
		buff.WriteString("\n")
		if p.Snake != nil {
			buff.WriteString("Direction: ")
			switch p.Snake.Dir {
			case Up:
				buff.WriteString("Up")
			case Right:
				buff.WriteString("Right")
			case Down:
				buff.WriteString("Down")
			case Left:
				buff.WriteString("Left")
			}
			buff.WriteString("\n")
		}
	}
	return buff
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
