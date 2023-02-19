package main

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	Up    = 0
	Down  = 1
	Left  = 2
	Right = 3
)

type Node struct {
	x, y int   // position of the node in the terminal
	next *Node // pointer to the next node in the snake
}

type Snake struct {
	head   *Node // pointer to the head node of the snake
	tail   *Node // pointer to the tail node of the snake
	length int   // length of the snake
	dir    int   // direction snake is facing (0 = Up, 1 = Down, 2 = Left, 3 = Right)
}

func NewSnake(x, y int) *Snake {
	// create the head node
	head := &Node{x: x, y: y}

	// create the snake with a length of 1 and no direction
	return &Snake{
		head:   head,
		tail:   head,
		length: 1,
		dir:    Up,
	}
}

func NewSnakeWithLength(x, y int) *Snake {
	// create the head node
	head := &Node{x: x, y: y}

	// create the middle node
	middle := &Node{x: x, y: y + 1}
	head.next = middle

	// create the tail node
	tail := &Node{x: x, y: y + 2}
	middle.next = tail

	// create the snake with length 3 and no direction
	return &Snake{
		head:   head,
		tail:   tail,
		length: 3,
		dir:    0,
	}
}

func (s *Snake) Move() {
	// calculate the new position of the head of the snake
	var newHeadX, newHeadY int
	switch s.dir {
	case Up:
		newHeadX, newHeadY = s.head.x, s.head.y-1
	case Right:
		newHeadX, newHeadY = s.head.x+1, s.head.y
	case Down:
		newHeadX, newHeadY = s.head.x, s.head.y+1
	case Left:
		newHeadX, newHeadY = s.head.x-1, s.head.y
	}

	// create a new head node and add it to the beginning of the snake
	newHead := &Node{newHeadX, newHeadY, s.head}
	s.head = newHead

	// if the snake hasn't grown, remove the last node
	if s.length == 0 {
		// find the new tail node
		curr := s.head
		for curr.next.next != nil {
			curr = curr.next
		}
		s.tail = curr
		s.tail.next = nil
	} else {
		s.length--
	}
}

func (s *Snake) Turn(dir int) {
	switch s.dir {
	case Up:
		if dir == Left || dir == Right {
			s.dir = dir
		}
	case Right:
		if dir == Up || dir == Down {
			s.dir = dir
		}
	case Down:
		if dir == Left || dir == Right {
			s.dir = dir
		}
	case Left:
		if dir == Up || dir == Down {
			s.dir = dir
		}
	}
}

func (s *Snake) CollidesWith(x, y int) bool {
	curr := s.head
	for curr != nil {
		if curr.x == x && curr.y == y {
			return true
		}
		curr = curr.next
	}
	return false
}

func (s *Snake) String() string {
	var buf bytes.Buffer

	buf.WriteString("Head: ")
	buf.WriteString(s.head.String())
	buf.WriteString("\n")

	buf.WriteString("Tail: ")
	buf.WriteString(s.tail.String())
	buf.WriteString("\n")

	buf.WriteString("Length: ")
	buf.WriteString(strconv.Itoa(s.length))
	buf.WriteString("\n")

	buf.WriteString("Direction: ")
	switch s.dir {
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
	for node := s.head; node != nil; node = node.next {
		buf.WriteString("  ")
		buf.WriteString(node.String())
		buf.WriteString("\n")
	}

	return buf.String()
}

func (n *Node) String() string {
	return fmt.Sprintf("[x:%d, y:%d]", n.x, n.y)
}
