package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/AlexEkdahl/snakes/pkg/game"
	"github.com/AlexEkdahl/snakes/pkg/network"
	"github.com/eiannone/keyboard"
)

// go:generate protoc --go_out=. --proto_path=pkg/network/ pkg/network/protobuf/message.proto

func main() {
	port := flag.String("port", ":8080", "The port number to use")
	flag.Parse()

	moveChan := make(chan game.Direction, 10)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go readInput(moveChan, interrupt)

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}

	client, err := network.NewClient(*port, moveChan)
	if err != nil {
		log.Fatal(err)
	}
	go client.Start()

	<-interrupt
	fmt.Println("\nShutting down the server...")
	keyboard.Close()
}

func readInput(mc chan game.Direction, stop chan os.Signal) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			continue
		}

		if key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
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