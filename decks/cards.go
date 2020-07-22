package decks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// CardSet is the set of int64's that uniquely identify a single card
type CardSet struct {
	DeckId int64 `json:"deckId"`
	CardId int64 `json:"cardId"`
}

// Deck contains information about a collection of cards
type Deck struct {
	Uuid        int64  `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cards       []Card `json:"cards"`
}

type Card struct {
	Uuid     int64  `json:"uuid"`
	Category string `json:"category"`
	Value    string `json:"value"`
}

func Equals(a, b Deck) bool {
	if a.Uuid != b.Uuid {
		return false
	}

	if len(a.Cards) != len(b.Cards) {
		return false
	}

	for i, v := range a.Cards {
		if v != b.Cards[i] {
			return false
		}
	}
	return true
}

// FromFile creates a new Deck from the provided json file path
func FromFile(filepath string) (Deck, error) {
	jsonFile, err := os.Open(filepath)
	defer jsonFile.Close()

	if err != nil {
		return Deck{}, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Deck{}, err
	}

	var deck Deck
	err = json.Unmarshal(byteValue, &deck)
	if err != nil {
		return Deck{}, err
	}

	return deck, nil
}

// DealNN deals up to count cards each to group number of people. Returning groups+1 number
// of newly dealt decks where the last deck is the remaining cards
func (d Deck) DealNN(groups, count int) []Deck {

	results := make([]Deck, groups+1)

	for j := 0; j < count; j++ {
		for i := 0; i < groups; i++ {
			if j*groups+i >= len(d.Cards) {
				break
			}

			results[i].Cards = append(results[i].Cards, d.Cards[j*groups+i])
		}
	}

	// Put remaining cards in last 'undealt' deck
	for i := count * groups; i < len(d.Cards); i++ {
		results[groups].Cards = append(results[groups].Cards, d.Cards[i])
	}

	return results
}

// Shuffle returns a new Deck with cards in a randomized order
func (d Deck) Shuffle() Deck {
	c := d
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c.Cards), func(i, j int) { c.Cards[i], c.Cards[j] = c.Cards[j], c.Cards[i] })
	return c
}

//OutOfBoundsError is an error thrown when requesting an index out of bounds
var InvalidIndexError error = errors.New("invalid index")

// Remove removes card at index returning card and deck with card removed
func (d Deck) Remove(idx int64) (Card, Deck, error) {

	i := -1
	for j, v := range d.Cards {
		if v.Uuid == idx {
			i = int(j)
		}
	}
	if i == -1 {
		return Card{}, Deck{}, fmt.Errorf("Index '%d' is %w", idx, InvalidIndexError)
	}

	card := d.Cards[i]
	deck := d
	deck.Cards = append(deck.Cards[:i], deck.Cards[i+1:]...)

	return card, deck, nil
}
