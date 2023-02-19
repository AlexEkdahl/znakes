package main

type SnakeNode struct {
	x, y int        // position of the node in the terminal
	next *SnakeNode // pointer to the next node in the snake
}

type Snake struct {
	head   *SnakeNode // pointer to the head node of the snake
	tail   *SnakeNode // pointer to the tail node of the snake
	length int        // length of the snake
}

func NewSnake(x, y int) *Snake {
	node := &SnakeNode{x: x, y: y}
	return &Snake{head: node, tail: node, length: 1}
}

// Move moves the snake one position in the specified direction
func (s *Snake) Move(row, col int) {
	// insert a new head node at the beginning of the linked list
	head := &SnakeNode{x: row, y: col, next: s.head}
	s.head = head
	s.length++

	// remove the tail node from the end of the linked list
	if s.length > 1 {
		node := s.tail
		for i := 1; i < s.length-1; i++ {
			node = node.next
		}
		node.next = nil
		s.tail = node
		s.length--
	}
}

// Grow adds a new segment to the end of the snake
func (s *Snake) Grow() {
	// insert a new tail node at the end of the linked list
	tail := &SnakeNode{x: s.tail.x, y: s.tail.y}
	s.tail.next = tail
	s.tail = tail
	s.length++
}
