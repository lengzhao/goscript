package vm

import (
	"testing"

	"github.com/lengzhao/goscript/context"
)

func TestBasicContext(t *testing.T) {
	parent := context.NewContext("parent", nil)
	if parent == nil {
		t.Fatal("Failed to create parent context")
	}

	if parent.GetPathKey() != "parent" {
		t.Errorf("Expected path key 'parent', got '%s'", parent.GetPathKey())
	}
}
