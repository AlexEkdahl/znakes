package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/AlexEkdahl/snakes/pkg/game"
	"github.com/AlexEkdahl/snakes/pkg/network"
)

func main() {
	port := flag.String("port", ":8080", "The port number to use")
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	g := game.NewGame(50, 19)
	server, err := network.NewServer(*port, g)
	if err != nil {
		log.Fatal(err)
	}
	go server.Start()

	<-interrupt
	fmt.Println("\nShutting down the server...")
	server.Stop()
}
