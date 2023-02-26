package network

import (
	"fmt"
	"net"

	"github.com/AlexEkdahl/snakes/pkg/game"
	"github.com/AlexEkdahl/snakes/pkg/network/protobuf"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

const (
	maxConcurrentConnections = 2
	protocol                 = "tcp"
)

type Server struct {
	game                  *game.Game
	conn                  net.Listener
	player                map[uuid.UUID]*net.Conn
	Messenger             Messenger
	concurrentConnections chan struct{}
	running               bool
	hasConnections        chan struct{}
}

func NewServer(port string, g *game.Game) (*Server, error) {
	c, err := net.Listen(protocol, port)
	if err != nil {
		return nil, err
	}
	fmt.Println("Server started and listening on :8080")

	return &Server{
		game:                  g,
		conn:                  c,
		running:               true,
		Messenger:             &protobufHandler{},
		concurrentConnections: make(chan struct{}, maxConcurrentConnections),
		player:                make(map[uuid.UUID]*net.Conn),
		hasConnections:        make(chan struct{}),
	}, nil
}

func (s *Server) Start() {
	// Start a goroutine to listen to the GameStateChan
	go s.handleConnections()

	// Wait for at least one connection
	<-s.hasConnections

	go s.broadcastGameSate()
	go s.game.Start()
}

func (s *Server) broadcastGameSate() {
	for gameState := range s.game.GameStateChan {
		msg, err := proto.Marshal(gameState)
		if err != nil {
			fmt.Printf("Error encoding game state message: %v\n", err)
			continue
		}

		for _, conn := range s.player {
			_, err := (*conn).Write(msg)
			if err != nil {
				fmt.Printf("Error sending game state to player: %v\n", err)
				continue
			}
		}
	}
}

func (s *Server) handleConnections() {
	for s.running {
		s.concurrentConnections <- struct{}{} // Acquire a slot in the channel
		conn, err := s.conn.Accept()
		if err != nil {
			<-s.concurrentConnections // Release the slot in the channel
			continue
		}

		p := game.NewPlayer(conn)
		s.game.Players = append(s.game.Players, p)
		s.player[p.ID] = &p.Conn
		fmt.Printf("New player connected, ID: %v\n", p.ID)
		go s.handlePlayer(p)

		// signal that there is at least one connection
		select {
		case <-s.hasConnections:
		default:
			close(s.hasConnections)
		}
	}
}

func (s *Server) handlePlayer(p *game.Player) {
	for {
		buf := make([]byte, 1024)
		n, err := p.Conn.Read(buf)
		if err != nil {
			fmt.Printf("Error reading message from player: %v\n", err)
			return
		}
		msg, err := s.Messenger.DecodeMessage(buf[:n])
		if err != nil {
			fmt.Printf("Error decode message: %v\n", err)
			return
		}

		switch m := msg.Type.(type) {
		case *protobuf.Message_Move:
			s.handleInputMessage(m, p.ID)
		case *protobuf.Message_Disconnect:
			p.Conn.Close()
			delete(s.player, p.ID)
			s.game.RemovePlayer(p.ID)
			<-s.concurrentConnections // Release the slot in the channel

			fmt.Printf("Player %v disconnected\n", p.ID)
			return
		}
	}
}

func (s *Server) handleInputMessage(m *protobuf.Message_Move, id uuid.UUID) {
	s.game.InputChan <- game.InputMessage{
		PlayerID: id,
		Input:    game.Direction(*m.Move.GetDirection().Enum()),
	}
}

func (s *Server) Stop() {
	s.game.Stop()
	s.running = false
	s.conn.Close()
}
