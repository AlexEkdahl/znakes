package network

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/AlexEkdahl/snakes/pkg/game"
)

type TelnetServer struct {
	game    *game.Game
	conn    net.Listener
	clients map[net.Conn]struct{}
	mutex   sync.Mutex
}

func NewTelnetServer(port string, g *game.Game) (*TelnetServer, error) {
	c, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}
	fmt.Println("Server started and listening on", port)

	return &TelnetServer{
		game:    g,
		conn:    c,
		clients: make(map[net.Conn]struct{}),
	}, nil
}

func (s *TelnetServer) Start() {
	// Start a goroutine to listen to the GameStateChan
	go s.game.Start()
	go s.handleConnections()
	go s.broadcastGameState()
}

func (s *TelnetServer) broadcastGameState() {
	for gameState := range s.game.GameStateChan {
		s.mutex.Lock()
		for c := range s.clients {
			clearScreen := []byte{27, 91, 72, 27, 91, 74}
			message := append(clearScreen, *gameState...) // Send game state to client
			if _, err := c.Write(message); err != nil {
				fmt.Printf("Error sending game state to client: %v\n", err)
			}
		}
		s.mutex.Unlock()
	}
}

func (s *TelnetServer) handleConnections() {
	for {
		conn, err := s.conn.Accept()
		if err != nil {
			continue
		}

		fmt.Printf("New client connected: %s\n", conn.RemoteAddr())

		s.mutex.Lock()
		s.clients[conn] = struct{}{}
		p := s.game.AddPlayer(conn)
		s.mutex.Unlock()

		go s.handleClient(*p)
	}
}

func (s *TelnetServer) handleClient(p game.Player) {
	defer func() {
		s.mutex.Lock()
		delete(s.clients, p.Conn)
		s.mutex.Unlock()
		p.Conn.Close()
	}()

	// Negotiate Telnet options
	if err := negotiateTelnetOptions(p.Conn); err != nil {
		fmt.Printf("Error negotiating Telnet options: %v\n", err)
		return
	}

	reader := bufio.NewReader(p.Conn)
	for {
		input, err := reader.ReadByte()
		if err != nil {
			break
		}

		switch input {
		case 'h':
			s.game.InputChan <- game.InputMessage{Input: game.Left, PlayerID: p.ID}
		case 'j':
			s.game.InputChan <- game.InputMessage{Input: game.Down, PlayerID: p.ID}
		case 'k':
			s.game.InputChan <- game.InputMessage{Input: game.Up, PlayerID: p.ID}
		case 'l':
			s.game.InputChan <- game.InputMessage{Input: game.Right, PlayerID: p.ID}
		case 'q':
			p.Conn.Close()
			delete(s.clients, p.Conn)
			s.game.RemovePlayer(p.ID)

			fmt.Printf("Player %v disconnected\n", p.ID)
			return
		}
	}
}

func negotiateTelnetOptions(conn net.Conn) error {
	// Telnet option negotiation
	binaryOption := []byte{255, 253, 34, 255, 250, 34, 1, 0, 255, 240, 255, 251, 1}
	if _, err := conn.Write(binaryOption); err != nil {
		return err
	}

	return nil
}

func (s *TelnetServer) Stop() {
	s.game.Stop()
	s.conn.Close()
}
