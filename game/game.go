package game

import (
	"bytes"
	"fmt"
	"math/rand"
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
	level     *level
	DrawBuff  bytes.Buffer
	Players   []*Player
	InputChan chan InputMessage
	food      []*Food
	mu        sync.Mutex
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
		DrawBuff:  bytes.Buffer{},
		Players:   []*Player{},
		InputChan: make(chan InputMessage),
		food:      []*Food{},
		mu:        sync.Mutex{},
	}
}

func NewFood(xBound, yBound int) *Food {
	rand.Seed(time.Now().UnixNano())
	return &Food{
		x: rand.Intn(xBound),
		y: rand.Intn(yBound),
	}
}

func newLevel(width, height int) *level {
	data := make([][]cell, height)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h] = make([]cell, width)
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
		time.Sleep(time.Millisecond * 1000)
		fmt.Println("linda")
	}
}

func (g *Game) Start() {
	go g.gameLoop()
}

func (g *Game) HandleCollision(playerNum int) {
	player := g.Players[playerNum]
	head := player.Snake.Head

	// Check for collision with walls
	if head.x < 0 || head.x >= g.level.width || head.y < 0 || head.y >= g.level.height {
		g.endGame(playerNum)
		return
	}

	// Check for collision with other snakes
	for i, otherPlayer := range g.Players {
		if i != playerNum {
			if g.snakeCollidesWithSnake(player.Snake, otherPlayer.Snake) {
				g.endGame(playerNum)
				return
			}
		}
	}
}

func (g *Game) endGame(playerNum int) {
	g.Players[playerNum].Conn.Close()
	g.Players = append(g.Players[:playerNum], g.Players[playerNum+1:]...)
}

func (g *Game) snakeCollidesWithSnake(snake1, snake2 *Snake) bool {
	for segment := snake2.Head; segment != nil; segment = segment.next {
		if snake1.collidesWithNode(segment) {
			return true
		}
	}
	return false
}

func (s *Snake) collidesWithNode(node *Node) bool {
	for segment := s.Head; segment != nil; segment = segment.next {
		if segment == node {
			return true
		}
	}
	return false
}

func (g *Game) renderLevel() {
	g.DrawBuff.Reset()

	for h := 0; h < g.level.height; h++ {
		for w := 0; w < g.level.width; w++ {
			switch g.level.data[h][w] {
			case wall:
				g.DrawBuff.WriteString("X")
			case snake:
				g.DrawBuff.WriteString("S")
			case fruit:
				g.DrawBuff.WriteString("F")
			case void:
				g.DrawBuff.WriteString(" ")
			}
		}
		g.DrawBuff.WriteString("\n")
	}
}
