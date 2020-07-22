package decks

import (
	"sync"
)

// Manager maintains the collection Decks
type Manager struct {
	decks map[int64]Deck
	m     sync.RWMutex
}

// NewManager creates instance of deck manager
func NewManager() *Manager {
	return &Manager{
		decks: map[int64]Deck{},
		m:     sync.RWMutex{},
	}
}

// Insert adds a deck to the manager
func (m *Manager) Insert(d Deck) {
	m.m.Lock()
	defer m.m.Unlock()
	m.decks[d.Uuid] = d
}

// Retrieve gets a deck from the manager
func (m *Manager) Retrieve(i int64) Deck {
	m.m.RLock()
	defer m.m.RUnlock()
	newDeck := Deck{
		Uuid:        m.decks[i].Uuid,
		Name:        m.decks[i].Name,
		Description: m.decks[i].Description,
		Cards:       make([]Card, len(m.decks[i].Cards)),
	}
	copy(newDeck.Cards, m.decks[i].Cards)
	return newDeck
}
