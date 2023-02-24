package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	VOID            = 0
	WALL            = 1
	SNAKE           = 2
	SNAKE_HEAD      = 4
	FRUIT           = 3
	VOID_ICON       = " "
	WALL_ICON       = "█"
	SNAKE_ICON      = "\033[31m█\033[0m"
	SNAKE_HEAD_ICON = "\033[37m█\033[0m"
	FRUIT_ICON      = "⛇"
)

type game struct {
	level    *level
	drawBuff bytes.Buffer
}

type level struct {
	height int
	width  int
	data   [][]int
}

func newGame(width, height int) *game {
	lvl := newLevel(width, height)
	return &game{
		level: lvl,
	}
}

func newLevel(width, height int) *level {
	data := make([][]int, height)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h] = make([]int, width)
		}
	}

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if h == 0 {
				data[h][w] = WALL
			}
			if h == height-1 {
				data[h][w] = WALL
			}
			if w == 0 {
				data[h][w] = WALL
			}
			if w == width-1 {
				data[h][w] = WALL
			}
		}
	}

	return &level{
		height: height,
		width:  width,
		data:   data,
	}
}

func (g *game) renderLevel() {
	g.drawBuff.Reset()

	for h := 0; h < g.level.height; h++ {
		for w := 0; w < g.level.width; w++ {
			switch g.level.data[h][w] {
			case WALL:
				g.drawBuff.WriteString(WALL_ICON)
			case SNAKE:
				g.drawBuff.WriteString(SNAKE_ICON)
			case SNAKE_HEAD:
				g.drawBuff.WriteString(SNAKE_HEAD_ICON)
			case FRUIT:
				g.drawBuff.WriteString(FRUIT_ICON)
			case VOID:
				g.drawBuff.WriteString(VOID_ICON)
			}
		}
		g.drawBuff.WriteString("\n")
	}
}

func (g *game) gameLoop() {
	s := NewSnakeWithLength(4, 4)
	s.Turn(Right)
	s.Move()
	for curr := s.head; curr != nil; curr = curr.next {
		if curr == s.head {
			g.level.data[curr.y][curr.x] = SNAKE_HEAD
			continue
		}
		g.level.data[curr.y][curr.x] = SNAKE
	}

	g.level.data[5][3] = FRUIT
	for {
		clearScreen()
		fmt.Println("s.Check(g.level.data)", s.Check(g.level.data))
		if s.Check(g.level.data) {
			s.Move()
		}

		for curr := s.head; curr != nil; curr = curr.next {
			if curr == s.head {
				g.level.data[curr.y][curr.x] = SNAKE_HEAD
				continue
			}
			g.level.data[curr.y][curr.x] = SNAKE
		}
		g.renderLevel()
		fmt.Print(g.drawBuff.String())
		time.Sleep(time.Millisecond * 32)
	}
}

func (g *game) start() {
	g.gameLoop()
}

func main() {
	g := newGame(40, 40)
	g.start()
}

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	default:
		fmt.Print("\033[2J") // clear entire screen
		fmt.Print("\033[H")  // move cursor to top-left corner
	}
}
