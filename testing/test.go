package testing

import (
	"bosch/stack" // Adjust the import path based on your project structure
	"testing"
)

func TestStack(t *testing.T) {
	// Initialize a new stack
	s := stack.InitStack()

	// Check if the stack is empty
	if !s.IsEmpty() {
		t.Errorf("Expected stack to be empty, but it is not")
	}

	// Push an item onto the stack
	s.Push(1)
	if s.IsEmpty() {
		t.Errorf("Expected stack not to be empty after pushing an item")
	}

	// Peek at the top item
	top := s.Peek()
	if top != 1 {
		t.Errorf("Expected top to be 1, but got %v", top)
	}

	// Pop the top item
	popped := s.Pop()
	if popped != 1 {
		t.Errorf("Expected popped value to be 1, but got %v", popped)
	}

	// Check if the stack is empty after popping
	if !s.IsEmpty() {
		t.Errorf("Expected stack to be empty after popping the last item")
	}

	// Test stack underflow
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on popping from empty stack, but no panic occurred")
		}
	}()
	s.Pop() // This should panic
}
