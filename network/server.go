package network

import (
	"fmt"
	"net"

	"github.com/AlexEkdahl/snakes/game"
	"github.com/AlexEkdahl/snakes/network/protobuf"
	"github.com/google/uuid"
)

const maxConcurrentConnections = 2

type Server struct {
	game                  *game.Game
	conn                  net.Listener
	player                map[uuid.UUID]*net.Conn
	Messenger             Messenger
	concurrentConnections chan struct{}
}

func NewServer(g *game.Game) (*Server, error) {
	c, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}
	fmt.Println("Server started and listening on :8080")

	return &Server{
		game:                  g,
		conn:                  c,
		player:                make(map[uuid.UUID]*net.Conn),
		Messenger:             &protobufHandler{},
		concurrentConnections: make(chan struct{}, maxConcurrentConnections),
	}, nil
}

func (s *Server) Start() {
	fmt.Println("shjkdsafhjkasdhfjkasfdhjk")
	// s.game.Start()
	for {
		s.concurrentConnections <- struct{}{} // Acquire a slot in the channel
		conn, err := s.conn.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			<-s.concurrentConnections // Release the slot in the channel
			continue
		}

		p := game.NewPlayer(conn)
		ale, err := s.Messenger.EncodeMessage()
		conn.Write(ale)
		fmt.Println("p", p)
		s.game.Players = append(s.game.Players, p)
		s.player[p.ID] = &p.Conn
		fmt.Printf("New player connected, ID: %v\n", p.ID)
		go s.handlePlayer(p)
	}
}

func (s *Server) handlePlayer(p *game.Player) {
	// Start listening for player input
	// go p.Listen()

	// Send the initial state to the player
	// initialState := s.game.SerializeGameState()
	// err := p.Write(initialState)
	// if err != nil {
	// 	fmt.Printf("Error writing to player connection: %v", err)
	// 	return
	// }
	defer func() {
		p.Conn.Close()
		// delete(s.player, p.ID)
		// s.game.Players = []
		<-s.concurrentConnections // Release the slot in the channel
	}()
	// Keep listening for updates from the player
	for {
		buf := make([]byte, 1024)
		n, err := p.Conn.Read(buf)
		msg, err := s.Messenger.DecodeMessage(buf[:n])
		if err != nil {
			fmt.Printf("Error reading message from player: %v\n", err)
			return
		}
		fmt.Println("server decode msg", msg)

		switch m := msg.Type.(type) {
		case *protobuf.Message_Move:
			s.game.InputChan <- game.InputMessage{
				PlayerID: p.ID,
				Input:    game.Direction(*m.Move.GetDirection().Enum()),
			}
		case *protobuf.Message_Disconnect:
			fmt.Printf("Player %v disconnected\n", p.ID)
			return
		}
	}
}

func (s *Server) handleInputMessage(msg *game.InputMessage) {
	s.game.InputChan <- *msg
}

func (s *Server) Stop() {
	s.conn.Close()
}
