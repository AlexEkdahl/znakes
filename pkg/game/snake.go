package game

import (
	"bytes"
	"fmt"
	"strconv"
)

type Node struct {
	x, y int   // position of the node in the terminal
	next *Node // pointer to the next node in the snake
}

type Snake struct {
	Head   *Node // pointer to the head node of the snake
	Tail   *Node // pointer to the tail node of the snake
	Length int   // length of the snake
	Dir    Direction
}

func NewSnake(x, y int) *Snake {
	head := &Node{x: x, y: y}

	return &Snake{
		Head:   head,
		Tail:   head,
		Length: 1,
		Dir:    Right,
	}
}

func NewSnakeWithLength(x, y int) *Snake {
	head := &Node{x: x, y: y}

	middle := &Node{x: x, y: y + 1}
	head.next = middle

	tail := &Node{x: x, y: y + 2}
	middle.next = tail

	return &Snake{
		Head:   head,
		Tail:   tail,
		Length: 3,
		Dir:    Right,
	}
}

func (s *Snake) Move() {
	// calculate the new position of the head of the snake
	var newHeadX, newHeadY int
	switch s.Dir {
	case Up:
		newHeadX, newHeadY = s.Head.x-1, s.Head.y
	case Right:
		newHeadX, newHeadY = s.Head.x, s.Head.y+1
	case Down:
		newHeadX, newHeadY = s.Head.x+1, s.Head.y
	case Left:
		newHeadX, newHeadY = s.Head.x, s.Head.y-1
	}

	// create a new head node and add it to the beginning of the snake
	newHead := &Node{newHeadX, newHeadY, s.Head}
	s.Head = newHead

	// if the snake hasn't grown, remove the last node
	if s.Length == 0 {
		// find the new tail node
		curr := s.Head
		for curr.next.next != nil {
			curr = curr.next
		}
		s.Tail = curr
		s.Tail.next = nil
	} else {
		s.Length--
	}
}

func (s *Snake) SetDirection(direction Direction) {
	if s.Dir == Up && direction == Down ||
		s.Dir == Down && direction == Up ||
		s.Dir == Left && direction == Right ||
		s.Dir == Right && direction == Left {
		return
	}

	s.Dir = direction
}

func (s *Snake) Occupies(x, y int) bool {
	curr := s.Head
	for curr != nil {
		if curr.x == x && curr.y == y {
			return true
		}
		curr = curr.next
	}
	return false
}

func (s *Snake) Check(matrix [][]int) bool {
	curr := s.Head
	for curr != nil {
		if matrix[curr.x][curr.y] == 1 {
			return false
		}
		curr = curr.next
	}
	return true
}

func (s *Snake) GetNodeLocations() [][]int {
	locations := make([][]int, s.Length)
	node := s.Head
	for i := 0; i < s.Length; i++ {
		if i == 0 {
			locations[i] = []int{node.x, node.y}
		} else {
			locations[i] = locations[i-1][:] // copy the previous slice
			locations[i][0], locations[i][1] = node.x, node.y
		}
		node = node.next
	}
	return locations
}

func (s *Snake) String() string {
	var buf bytes.Buffer

	buf.WriteString("Head: ")
	buf.WriteString(s.Head.String())
	buf.WriteString("\n")

	buf.WriteString("Tail: ")
	buf.WriteString(s.Tail.String())
	buf.WriteString("\n")

	buf.WriteString("Length: ")
	buf.WriteString(strconv.Itoa(s.Length))
	buf.WriteString("\n")

	buf.WriteString("Direction: ")
	switch s.Dir {
	case Up:
		buf.WriteString("Up")
	case Right:
		buf.WriteString("Right")
	case Down:
		buf.WriteString("Down")
	case Left:
		buf.WriteString("Left")
	}
	buf.WriteString("\n")

	buf.WriteString("Nodes:\n")
	for node := s.Head; node != nil; node = node.next {
		buf.WriteString("  ")
		buf.WriteString(node.String())
		buf.WriteString("\n")
	}

	return buf.String()
}

func (n *Node) String() string {
	return fmt.Sprintf("[x:%d, y:%d]", n.x, n.y)
}
