package entity

import (
	"sync"
)

// Tree represents a tree entity
type Tree struct {
	mu sync.RWMutex

	ID       string
	Position Position
	HP       float64
	Wood     int // Amount of wood this tree provides when chopped
}

// NewTree creates a new tree entity
func NewTree(id string, pos Position) *Tree {
	return &Tree{
		ID:       id,
		Position: pos,
		HP:       100,
		Wood:     10,
	}
}

// TakeDamage reduces the tree's HP by the specified amount
func (t *Tree) TakeDamage(damage float64) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.HP -= damage
	return t.HP <= 0
}

// GetData returns the entity data for AI processing
func (t *Tree) GetData() map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return map[string]interface{}{
		"id":       t.ID,
		"type":     "TREE",
		"position": t.Position,
		"hp":       t.HP,
	}
}
