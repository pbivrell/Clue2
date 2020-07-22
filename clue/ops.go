package clue

import (
	"errors"
	"fmt"

	"github.com/pbivrell/clue/decks"
)

func ReleaseHands(userid string, cards []decks.CardSet) func(map[string]User, map[int64]decks.Deck) error {
	return func(users map[string]User, globalDecks map[int64]decks.Deck) error {

		user, has := users[userid]
		if !has {
			return fmt.Errorf("%w: user '%s' does not exist", InvalidUserError, userid)
		}

		for _, card := range cards {
			aCard, deck, err := user.decks[card.DeckId].Remove(card.CardId)
			if err != nil {
				return err
			}
			user.decks[card.DeckId] = deck
			x := globalDecks[card.DeckId]
			x.Cards = append(globalDecks[card.DeckId].Cards, aCard)
			globalDecks[card.DeckId] = x
		}
		return nil
	}
}

var InvalidUserError = errors.New("Invalid user")

func InputHands(userid string, cards []decks.CardSet) func(map[string]User, map[int64]decks.Deck) error {
	return func(users map[string]User, globalDecks map[int64]decks.Deck) error {

		user, has := users[userid]
		if !has {
			return fmt.Errorf("%w: user '%s' does not exist", InvalidUserError, userid)
		}

		for _, card := range cards {
			aCard, deck, err := globalDecks[card.DeckId].Remove(card.CardId)
			if err != nil {
				return err
			}
			globalDecks[card.DeckId] = deck
			x := user.decks[card.DeckId]
			x.Cards = append(user.decks[card.DeckId].Cards, aCard)
			user.decks[card.DeckId] = x
		}
		return nil
	}
}

var NoSolutionError = errors.New("no valid solution")

func DrawSolution(userid string, selector []string) func(map[string]User, map[int64]decks.Deck) error {
	return func(users map[string]User, globalDecks map[int64]decks.Deck) error {

		solutions := make([]decks.Card, len(selector))
		for _, card := range globalDecks[0].Cards {
			for i, selected := range selector {
				if card.Category == selected {
					solutions[i] = card
				}
			}
		}

		var blank decks.Card

		for _, card := range solutions {
			if card == blank {
				return NoSolutionError
			} else {
				_, deck, err := globalDecks[0].Remove(card.Uuid)
				if err != nil {
					return err
				}
				globalDecks[0] = deck
			}
		}

		user, ok := users[userid]
		if !ok {
			return fmt.Errorf("%w: user '%s' does not exist", InvalidUserError, userid)
		}
		x := user.decks[0]
		x.Cards = solutions
		user.decks[0] = x
		return nil
	}
}

func ShuffleDeck() func(map[string]User, map[int64]decks.Deck) error {
	return func(users map[string]User, globalDecks map[int64]decks.Deck) error {
		for k, d := range globalDecks {
			newDeck := d.Shuffle()
			globalDecks[k] = newDeck
		}
		return nil
	}
}

func DealHands(userids []string) func(map[string]User, map[int64]decks.Deck) error {
	return func(users map[string]User, globalDecks map[int64]decks.Deck) error {
		deck := globalDecks[0]
		newDecks := deck.DealNN(len(userids), len(deck.Cards))
		for i, userid := range userids {
			user, ok := users[userid]
			if !ok {
				return fmt.Errorf("%w: user '%s' does not exist", InvalidUserError, userid)
			}
			user.decks[0] = newDecks[i]
		}
		globalDecks[0] = newDecks[len(userids)]
		return nil
	}

}
