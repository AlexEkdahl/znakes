package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/AlexEkdahl/snakes/game"
	"github.com/AlexEkdahl/snakes/network"
	"github.com/eiannone/keyboard"
)

func main() {
	isServer, port, err := parseArgs()
	if err != nil {
		panic(err)
	}

	moveChan := make(chan game.Direction, 10)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go readInput(moveChan, interrupt)

	if isServer {

		fmt.Println("isServer, port", isServer, port)
		g := game.NewGame(50, 19)
		server, err := network.NewServer("localhost:8080", g)
		if err != nil {
			log.Fatal(err)
		}
		go server.Start()
		<-interrupt
		fmt.Println("\nShutting down the server...")
		server.Stop()
	}
	if !isServer {

		err = keyboard.Open()
		if err != nil {
			panic(err)
		}

		client, err := network.NewClient("localhost:8080", moveChan)
		if err != nil {
			log.Fatal(err)
		}
		go client.Start()

		<-interrupt
		fmt.Println("\nShutting down the server...")
		keyboard.Close()
	}
	// Start the server

	// Start the client
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

func parseArgs() (isServer bool, port int, err error) {
	// Define command-line flags
	serverFlag := flag.Bool("server", false, "Starts the program as a server")
	clientFlag := flag.Bool("client", false, "Starts the program as a client")
	portFlag := flag.Int("port", 8080, "The port number to use")

	// Parse command-line flags
	flag.Parse()

	// Check that either the server or client flag is set
	if !*serverFlag && !*clientFlag {
		err = fmt.Errorf("must specify either --server or --client flag")
		return
	}

	// Check that only one of the server or client flag is set
	if *serverFlag && *clientFlag {
		err = fmt.Errorf("cannot specify both --server and --client flags")
		return
	}

	// Set the isServer flag based on the server flag
	isServer = *serverFlag

	// Set the port based on the port flag
	port = *portFlag

	return
}
