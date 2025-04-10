package entity

import (
	"sync"
)

// WoodcutterState represents the current state of a woodcutter
type WoodcutterState string

const (
	StateIdle      WoodcutterState = "IDLE"
	StateLearning  WoodcutterState = "LEARNING"
	StateSearching WoodcutterState = "SEARCHING"
	StateMoving    WoodcutterState = "MOVING"
	StateChopping  WoodcutterState = "CHOPPING"
)

// Woodcutter represents a woodcutter entity
type Woodcutter struct {
	mu sync.RWMutex

	ID        string
	Position  Position
	HP        float64
	Energy    float64
	State     WoodcutterState
	Skills    map[string]bool // Map of skill name to whether it's learned
	Inventory map[string]int  // Map of item name to count
}

// Position represents a 3D position
type Position struct {
	X, Y, Z float64
}

// NewWoodcutter creates a new woodcutter entity
func NewWoodcutter(id string, pos Position) *Woodcutter {
	return &Woodcutter{
		ID:        id,
		Position:  pos,
		HP:        100,
		Energy:    100,
		State:     StateIdle,
		Skills:    make(map[string]bool),
		Inventory: make(map[string]int),
	}
}

// LearnSkill adds a new skill to the woodcutter
func (w *Woodcutter) LearnSkill(skillName string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Skills[skillName] = true
}

// HasSkill checks if the woodcutter has a specific skill
func (w *Woodcutter) HasSkill(skillName string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Skills[skillName]
}

// UpdateState updates the woodcutter's state
func (w *Woodcutter) UpdateState(state WoodcutterState) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.State = state
}

// MoveTo updates the woodcutter's position
func (w *Woodcutter) MoveTo(pos Position) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Position = pos
}

// AddToInventory adds items to the woodcutter's inventory
func (w *Woodcutter) AddToInventory(itemName string, count int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Inventory[itemName] += count
}

// GetData returns the entity data for AI processing
func (w *Woodcutter) GetData() map[string]interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return map[string]interface{}{
		"id":       w.ID,
		"type":     "WOODCUTTER",
		"position": w.Position,
		"hp":       w.HP,
		"energy":   w.Energy,
		"state":    string(w.State),
		"skills":   w.Skills,
	}
}
