package main

import (
	"testing"
)

func TestSnake(t *testing.T) {
	// create a new snake with a single segment
	snake := NewSnake(0, 0)

	// test that the snake has length 1 and is located at the correct position
	if snake.length != 1 {
		t.Errorf("Expected snake length to be 1, got %d", snake.length)
	}
	if snake.head.x != 0 || snake.head.y != 0 {
		t.Errorf("Expected snake to be at position (0, 0), got (%d, %d)", snake.head.x, snake.head.y)
	}
}
