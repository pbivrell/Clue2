package clue

import (
	"errors"
	"fmt"

	"github.com/pbivrell/clue/decks"
	"github.com/pbivrell/clue/sessions"
)

type User struct {
	uuid  string
	decks map[int64]decks.Deck
}

func (u User) Decks() int {
	return len(u.decks)
}

func (u User) Deck(i int64) decks.Deck {
	if i < 0 || int(i) >= len(u.decks) {
		return decks.Deck{
			Cards: make([]decks.Card, 0),
		}
	}

	return u.decks[i]
}

type State struct {
	users map[string]User
	decks map[int64]decks.Deck
}

type ActiveState struct {
	releaseFunc func()
	*State
}

// UnexpectedStateError is an error that indicates the data stored at the session
// was not the type that was expected
var UnexpectedStateError = errors.New("Unexpected session state type")

// GetSessionState wraps session.Manager.Retrieve returning the state as a concret
// type instead of an interface type
func GetSessionState(sm *sessions.Manager, sessionid string) (*ActiveState, error) {
	data, releaseFunc, err := sm.Retrieve(sessionid)
	if err != nil {
		// We shouldn't need to do this but just to be safe
		releaseFunc()
		return &ActiveState{}, err
	}

	var state *State
	var ok bool
	if state, ok = data.(*State); !ok {
		// The data was not the type we expected release the lock
		releaseFunc()
		return &ActiveState{}, fmt.Errorf("%w: session state for '%s' was of type '%T' expected '%T'", UnexpectedStateError, sessionid, data, state)
	}

	return &ActiveState{releaseFunc, state}, nil
}

func NewState(ds ...decks.Deck) *State {

	deckMap := make(map[int64]decks.Deck)
	for _, deck := range ds {
		deckMap[deck.Uuid] = deck
	}

	return &State{
		users: map[string]User{},
		decks: deckMap,
	}
}

func (a *ActiveState) Join(uuid string) string {
	defer a.releaseFunc()
	a.users[uuid] = User{
		uuid:  uuid,
		decks: map[int64]decks.Deck{},
	}
	return uuid
}

func (a *ActiveState) Op(op func(map[string]User, map[int64]decks.Deck) error) error {
	defer a.releaseFunc()
	return op(a.users, a.decks)
}

func (a *ActiveState) Dump() (map[string]User, map[int64]decks.Deck) {
	defer a.releaseFunc()

	users := make(map[string]User)
	decks := make(map[int64]decks.Deck)

	for k, v := range a.users {
		users[k] = v
	}

	for k, v := range a.decks {
		decks[k] = v
	}

	return users, decks
}
