package game

import (
	"bytes"
	"math/rand"
	"time"
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
	Level    *Level
	DrawBuff bytes.Buffer
	Snakes   []*Snake
	Players  []*Player
}

// type GameState struct {
// 	Board [][]bool
// 	Food  []Coord
// }

type Level struct {
	height int
	width  int
	data   [][]cell
}

func newLevel(width, height int) *Level {
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

	return &Level{
		height: height,
		width:  width,
		data:   data,
	}
}

func (g *Game) renderLevel() {
	g.DrawBuff.Reset()

	for h := 0; h < g.Level.height; h++ {
		for w := 0; w < g.Level.width; w++ {
			switch g.Level.data[h][w] {
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

type Food struct {
	x, y int // position of the fruit in the terminal
}

// NewFood creates a new food at a random position within the specified bounds
func NewFood(xBound, yBound int) *Food {
	rand.Seed(time.Now().UnixNano())
	return &Food{
		x: rand.Intn(xBound),
		y: rand.Intn(yBound),
	}
}
