package main

import (
	"testing"
)

func TestNewSnake(t *testing.T) {
	snake := NewSnake(10, 20)
	if snake.head != snake.tail {
		t.Error("snake head and tail are not equal")
	}
	if snake.head.x != 10 || snake.head.y != 20 {
		t.Error("snake head position is incorrect")
	}
	if snake.length != 1 {
		t.Error("snake length is incorrect")
	}
	if snake.dir != Up {
		t.Error("snake direction is incorrect")
	}
}

func TestNewSnakeWithLength(t *testing.T) {
	snake := NewSnakeWithLength(10, 20)
	if snake.tail != snake.head.next.next {
		t.Error("snake tail is not correct")
	}
	if snake.head.x != 10 || snake.head.y != 20 {
		t.Error("snake head position is incorrect")
	}
	if snake.length != 3 {
		t.Error("snake length is incorrect")
	}
	if snake.dir != Up {
		t.Error("snake direction is incorrect")
	}
}

func TestSnakeMove(t *testing.T) {
	snake := NewSnake(10, 20)
	snake.Move()
	if snake.head.y != 19 {
		t.Error("snake head y position is incorrect")
	}
	if snake.length != 0 {
		t.Error("snake length should be 0 after moving")
	}
	snake.Turn(Right)
	snake.Move()
	if snake.head.x != 11 {
		t.Error("snake head x position is incorrect")
	}
	if snake.dir != Right {
		t.Error("snake direction is incorrect")
	}
	snake.Move()
	if snake.head.x != 12 {
		t.Error("snake head x position is incorrect")
	}
	snake.Turn(Down)
	snake.Move()
	if snake.head.y != 20 {
		t.Error("snake head y position is incorrect")
	}
}

func TestSnakeTurn(t *testing.T) {
	snake := NewSnake(0, 0)

	// turn right from default direction (up)
	snake.Turn(Right)
	if snake.dir != Right {
		t.Errorf("Expected snake direction to be Right, got %d", snake.dir)
	}

	// turn down from right
	snake.Turn(Down)
	if snake.dir != Down {
		t.Errorf("Expected snake direction to be Down, got %d", snake.dir)
	}

	// turn left from down
	snake.Turn(Left)
	if snake.dir != Left {
		t.Errorf("Expected snake direction to be Left, got %d", snake.dir)
	}

	// turn up from left
	snake.Turn(Up)
	if snake.dir != Up {
		t.Errorf("Expected snake direction to be Up, got %d", snake.dir)
	}

	// turn up from up (invalid turn)
	snake.Turn(Up)
	if snake.dir != Up {
		t.Errorf("Expected snake direction to be unchanged (Up), got %d", snake.dir)
	}
}
