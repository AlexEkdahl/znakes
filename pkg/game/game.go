package game

import (
	"bytes"
	"math/rand"
	"net"
	"strconv"
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

type Fruit struct {
	X int
	Y int
}

type Game struct {
	level         *level
	Players       map[uuid.UUID]*Player
	Fruits        []*Fruit
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
		Players:       make(map[uuid.UUID]*Player),
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
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		go g.updatePlayerSnakes()

		// add new fruit to the game board
		g.spawnFruits()

		g.mu.Lock()
		gameState := g.renderLevel()
		g.GameStateChan <- &gameState
		g.mu.Unlock()
	}
}

func (g *Game) updatePlayerSnakes() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, player := range g.Players {
		if player.Snake != nil {
			player.Snake.Move(g.level.width, g.level.height)

			// check if snake has eaten any fruits
			for i := 0; i < len(g.Fruits); i++ {
				fruit := g.Fruits[i]
				if player.Snake.Occupies(fruit.Y, fruit.X, g.level.height, g.level.width) {
					player.Score++
					g.Fruits = append(g.Fruits[:i], g.Fruits[i+1:]...)
				}
			}
		}
	}
}

func (g *Game) spawnFruits() {
	if len(g.Fruits) < 2 {
		for i := len(g.Fruits); i < 2; i++ {
			g.Fruits = append(g.Fruits, &Fruit{rand.Intn(g.level.width - 1), rand.Intn(g.level.height - 1)})
		}
	}
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
			for _, player := range g.Players {
				if player.Snake.Occupies(h, w, g.level.height, g.level.width) {
					buff.WriteByte('S')
					occupied = true
					break
				}
			}
			if !occupied {
				// check for fruits at this position
				for _, fruit := range g.Fruits {
					if fruit.X == w && fruit.Y == h {
						buff.WriteByte('F')
						occupied = true
						break
					}
				}

				if !occupied {
					switch g.level.data[h][w] {
					case wall:
						buff.WriteByte('X')
					case void:
						buff.WriteByte(' ')
					}
				}
			}
		}
		buff.Write([]byte{telnetControlChars["CarriageReturn"], telnetControlChars["LineFeed"]})
	}

	for _, player := range g.Players {
		msg := "Player score: " + strconv.Itoa(player.Score)
		buff.WriteString(msg)
		buff.Write([]byte{telnetControlChars["CarriageReturn"], telnetControlChars["LineFeed"]})
	}
	return buff.Bytes()
}

func (g *Game) RemovePlayer(playerID uuid.UUID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.Players, playerID)
}

func (g *Game) AddPlayer(conn net.Conn) *Player {
	p := NewPlayer(conn)
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Players[p.ID] = p
	return p
}

func (g *Game) handleInput() {
	for input := range g.InputChan {
		if player, ok := g.Players[input.PlayerID]; ok {
			player.Snake.SetDirection(input.Input)
		}
	}
}
