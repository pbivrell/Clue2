package sessions

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

// CollisionError is an error type representing when a newly generated guid
// matches an already exists one. This should be impossible to observe
var CollisionError error = errors.New("guid collision")

// SessionError is an error type representing an attempt to opperate on an invalid sessions
var SessionError error = errors.New("invalid session")

// Manager mantains the running state for a card game by session guid
type Manager struct {
	sessions map[string]State
	m        sync.RWMutex
}

// State is the data stored in this session it is protected from concurrent
// access by the mutex
type State struct {
	m     sync.Mutex
	state interface{}
}

// NewManager maintains the state structure by generated session id
func NewManager() *Manager {
	return &Manager{
		sessions: map[string]State{},
		m:        sync.RWMutex{},
	}
}

func NewGuid() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b[0:2])
}

// NewSession creates a new session with the provided deck, returning
// the sessionid
func (m *Manager) NewSession(state interface{}) (string, error) {
	m.m.Lock()
	defer m.m.Unlock()
	uuid := NewGuid()
	if _, has := m.sessions[uuid]; has {
		return "", fmt.Errorf("%w: '%s' is an already existing session id", CollisionError, uuid)
	}
	m.sessions[uuid] = State{
		m:     sync.Mutex{},
		state: state,
	}
	return uuid, nil
}

// Retrieve gets exclusive access to the state of provided sessionid and func
// that releases exclusive access to the state
func (m *Manager) Retrieve(sessionid string) (interface{}, func(), error) {
	m.m.RLock()

	state, has := m.sessions[sessionid]
	if !has {
		defer m.m.RUnlock()
		return nil, func() {}, fmt.Errorf("%w: session '%s' does not exist", SessionError, sessionid)
	}

	// Second prevent all concurrent access to the retrieved state
	defer func() {
		state.m.Lock()
	}()
	// First unlock mutex protecting entire session data structure
	defer m.m.RUnlock()

	// This function is called when the caller is done maniulating the state
	releaseFunc := func() {
		// Ensure that the release can only be called once
		var once sync.Once
		once.Do(func() {
			state.m.Unlock()
		})
	}

	return state.state, releaseFunc, nil

}
