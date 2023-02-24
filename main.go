package main

import (
	"fmt"

	"github.com/AlexEkdahl/snakes/game"
	"github.com/AlexEkdahl/snakes/network"
	"github.com/eiannone/keyboard"
)

func main() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	go readInput()

	g := game.NewGame(10, 10)
	s, err := network.NewServer(g)
	if err != nil {
		fmt.Println("err", err)
	}
	c, err := network.NewClient(":8080")
	if err != nil {
		fmt.Println("err", err)
	}

	go s.Start()
	c.Start()
}

func readInput() {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		fmt.Printf("key: %v, char: %c\n", key, char)

		if char == 'w' {
			// move snake up
		} else if char == 's' {
			// move snake down
		} else if char == 'a' {
			// move snake left
		} else if char == 'd' {
			// move snake right
		}
	}
}
