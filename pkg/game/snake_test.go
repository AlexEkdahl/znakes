package game

import (
	"testing"
)

func TestNewSnake(t *testing.T) {
	s := NewSnake(1, 1)

	if s.Head.x != 1 || s.Head.y != 1 {
		t.Errorf("Expected head position to be (1,1), got (%d,%d)", s.Head.x, s.Head.y)
	}

	if s.Tail.x != 1 || s.Tail.y != 2 {
		t.Errorf("Expected tail position to be (1,2), got (%d,%d)", s.Tail.x, s.Tail.y)
	}

	if s.Length != 2 {
		t.Errorf("Expected snake length to be 2, got %d", s.Length)
	}

	if s.Dir != Right {
		t.Errorf("Expected snake direction to be Right, got %v", s.Dir)
	}
}

func TestNewSnakeWithLength(t *testing.T) {
	s := NewSnakeWithLength(1, 1)

	if s.Head.x != 1 || s.Head.y != 1 {
		t.Errorf("Expected head position to be (1,1), got (%d,%d)", s.Head.x, s.Head.y)
	}

	if s.Tail.x != 1 || s.Tail.y != 3 {
		t.Errorf("Expected tail position to be (1,3), got (%d,%d)", s.Tail.x, s.Tail.y)
	}

	if s.Length != 3 {
		t.Errorf("Expected snake length to be 3, got %d", s.Length)
	}

	if s.Dir != Right {
		t.Errorf("Expected snake direction to be Right, got %v", s.Dir)
	}
}

func TestSnake_RemoveTail(t *testing.T) {
	s := NewSnake(1, 1)
	s.RemoveTail()

	if s.Head.x != 1 || s.Head.y != 1 {
		t.Errorf("Expected head position to be (1,1), got (%d,%d)", s.Head.x, s.Head.y)
	}

	if s.Tail.x != 1 || s.Tail.y != 1 {
		t.Errorf("Expected tail position to be (1,1), got (%d,%d)", s.Tail.x, s.Tail.y)
	}

	if s.Length != 1 {
		t.Errorf("Expected snake length to be 1, got %d", s.Length)
	}

	if s.Dir != Right {
		t.Errorf("Expected snake direction to be Right, got %v", s.Dir)
	}
}
