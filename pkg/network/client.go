package network

import (
	"net"

	"github.com/AlexEkdahl/snakes/pkg/game"
	"github.com/AlexEkdahl/snakes/pkg/network/protobuf"
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
	isRunning bool
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
		isRunning: true,
	}, nil
}

func (c *Client) Start() {
	go c.sendMoves()
	c.handlerMessages()
}

func (c *Client) handlerMessages() {
	for c.isRunning {
		buf := make([]byte, 1024)
		n, err := c.conn.Read(buf)
		if err != nil {
			continue
		}

		msg, err := c.Messenger.DecodeMessage(buf[:n])
		if err != nil {
			continue
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
			continue
		}

		if _, err := c.conn.Write(b); err != nil {
			continue
		}
	}
}

func (c *Client) disconect() {
	msg := &protobuf.Message{
		Type: &protobuf.Message_Disconnect{},
	}

	b, err := c.Messenger.EncodeMessage(msg)
	if err != nil {
		return
	}
	if _, err := c.conn.Write(b); err != nil {
		return
	}
}

func (c *Client) Stop() {
	c.isRunning = false
	c.disconect()
	c.conn.Close()
}
