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
	VOID        = 0
	WALL        = 1
	PLAYER      = 2
	GHOST       = 3
	VOID_ICON   = " "
	WALL_ICON   = "█"
	PLAYER_ICON = "☻"
	GHOST_ICON  = "⛇"
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
			case PLAYER:
				g.drawBuff.WriteString(PLAYER_ICON)
			case GHOST:
				g.drawBuff.WriteString(GHOST_ICON)
			case VOID:
				g.drawBuff.WriteString(VOID_ICON)
			}
		}
		g.drawBuff.WriteString("\n")
	}
}

func (g *game) gameLoop() {
	g.level.data[2][3] = PLAYER
	g.level.data[5][3] = GHOST
	for {
		clearScreen()
		g.renderLevel()
		fmt.Print(g.drawBuff.String())
		time.Sleep(time.Millisecond * 32)
	}
}

func (g *game) start() {
	g.gameLoop()
}

func main() {
	g := newGame(80, 18)
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
