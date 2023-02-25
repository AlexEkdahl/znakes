package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

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

	moveChan := make(chan game.Direction, 10)
	interrupt := make(chan os.Signal, 1)

	go readInput(moveChan, interrupt)
	g := game.NewGame(50, 15)

	// Start the server
	server, err := network.NewServer("localhost:8080", g)
	if err != nil {
		log.Fatal(err)
	}
	go server.Start()

	// Start the client
	client, err := network.NewClient("localhost:8080", moveChan)
	if err != nil {
		log.Fatal(err)
	}
	go client.Start()

	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	fmt.Println("\nShutting down the server...")
	server.Stop()
	time.Sleep(1 * time.Second)
}

func readInput(mc chan game.Direction, stop chan os.Signal) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			continue
		}

		if key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
			keyboard.Close()
			stop <- os.Interrupt
			return
		}

		switch char {
		case 'k':
			mc <- game.Up
		case 'j':
			mc <- game.Down
		case 'h':
			mc <- game.Left
		case 'l':
			mc <- game.Right
		default:
			continue
		}
	}
}
