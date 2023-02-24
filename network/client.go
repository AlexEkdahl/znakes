package network

import (
	"fmt"
	"net"
)

// Client is a client for the game server
type Client struct {
	conn      net.Conn
	ID        string
	Messenger Messenger
}

// NewClient creates a new client with a connection to the given address
func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:      conn,
		Messenger: &protobufHandler{},
	}, nil
}

// Start starts the client
func (c *Client) Start() {
	c.handlerMessages()
}

func (c *Client) handlerMessages() error {
	// Read and handle messages from the server
	msg, err := c.Messenger.EncodeMessage()
	if err != nil {
		return fmt.Errorf("error writing join message: %v", err)
	}
	c.conn.Write(msg)

	fmt.Println("client write msg", msg)
	for {
		msg, err := c.Messenger.EncodeMessage()
		if err != nil {
			return fmt.Errorf("error writing join message: %v", err)
		}
		c.conn.Write(msg)
		buf := make([]byte, 1024)
		n, err := c.conn.Read(buf)

		alex, err := c.Messenger.DecodeMessage(buf[:n])
		if err != nil {
			fmt.Printf("Error reading message from player: %v\n", err)
		}

		fmt.Println("client decoded", alex)
		switch m := alex.Type.(type) {
		default:
			fmt.Println("client decoded", m)
			fmt.Println("client decoded", alex.String())
		}
	}
}
