package network

import (
	"fmt"
	"net"

	"github.com/AlexEkdahl/snakes/game"
	"github.com/AlexEkdahl/snakes/network/protobuf"
)

type MessagePrinter interface {
	PrintMessage(msg interface{})
	Clear()
}

type Client struct {
	conn      net.Conn
	ID        string
	Messenger Messenger
	Printer   MessagePrinter
	moveChan  chan game.Direction
}

func NewClient(addr string, mc chan game.Direction) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:      conn,
		Messenger: &protobufHandler{},
		Printer:   NewPrinter(),
		moveChan:  mc,
	}, nil
}

func (c *Client) Start() {
	go c.sendMoves()
	c.handlerMessages()
}

func (c *Client) handlerMessages() error {
	for {
		buf := make([]byte, 1024)
		n, err := c.conn.Read(buf)

		msg, err := c.Messenger.DecodeMessage(buf[:n])
		if err != nil {
			fmt.Printf("Error reading message from player: %v\n", err)
		}

		switch msg.Type.(type) {
		case *protobuf.Message_Game:
			c.Printer.PrintMessage(msg.GetGame())
		default:
			continue

		}
	}
}

func (c *Client) sendMoves() {
	for move := range c.moveChan {
		// Send the move to the server
		msg := &protobuf.Message{
			Type: &protobuf.Message_Move{
				Move: &protobuf.MoveMessage{
					Direction: protobuf.Direction(move),
				},
			},
		}
		b, err := c.Messenger.EncodeMessage(msg)
		if err != nil {
			fmt.Printf("Error sending move to server: %v\n", err)
		}
		c.conn.Write(b)
	}
}

func (c *Client) Stop() {
	c.conn.Close()
}
