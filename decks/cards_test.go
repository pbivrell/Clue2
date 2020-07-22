package decks

import (
	"fmt"
	"testing"
)

func TestDealNN(t *testing.T) {
	var tests = []struct {
		deck   Deck
		groups int
		count  int
		result []Deck
	}{
		{
			groups: 0,
			count:  0,
			deck: Deck{
				Cards: []Card{},
			},
			result: make([]Deck, 1),
		},
		{
			groups: 5,
			count:  0,
			deck: Deck{
				Cards: []Card{},
			},
			result: make([]Deck, 6),
		},
		{
			groups: 5,
			count:  8,
			deck: Deck{
				Cards: []Card{},
			},
			result: make([]Deck, 6),
		},
		{
			groups: 3,
			count:  1,
			deck: Deck{
				Cards: []Card{
					Card{
						Uuid: 0,
					},
					Card{
						Uuid: 1,
					},
					Card{
						Uuid: 2,
					},
				},
			},
			result: []Deck{
				Deck{
					Cards: []Card{
						Card{
							Uuid: 0,
						},
					},
				},
				Deck{
					Cards: []Card{
						Card{
							Uuid: 1,
						},
					},
				},
				Deck{
					Cards: []Card{
						Card{
							Uuid: 2,
						},
					},
				},
				Deck{},
			},
		}, {
			groups: 3,
			count:  1,
			deck: Deck{
				Cards: []Card{
					Card{
						Uuid: 0,
					},
					Card{
						Uuid: 1,
					},
					Card{
						Uuid: 2,
					},
					Card{
						Uuid: 3,
					},
				},
			},
			result: []Deck{
				Deck{
					Cards: []Card{
						Card{
							Uuid: 0,
						},
					},
				},
				Deck{
					Cards: []Card{
						Card{
							Uuid: 1,
						},
					},
				},
				Deck{
					Cards: []Card{
						Card{
							Uuid: 2,
						},
					},
				},
				Deck{
					Cards: []Card{
						Card{
							Uuid: 3,
						},
					},
				},
			},
		}, {
			groups: 3,
			count:  2,
			deck: Deck{
				Cards: []Card{
					Card{
						Uuid: 0,
					},
					Card{
						Uuid: 1,
					},
					Card{
						Uuid: 2,
					},
					Card{
						Uuid: 3,
					},
				},
			},
			result: []Deck{
				Deck{
					Cards: []Card{
						Card{
							Uuid: 0,
						},
						Card{
							Uuid: 3,
						},
					},
				},
				Deck{
					Cards: []Card{
						Card{
							Uuid: 1,
						},
					},
				},
				Deck{
					Cards: []Card{
						Card{
							Uuid: 2,
						},
					},
				},
				Deck{
					Cards: []Card{},
				},
			},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("DealNN(%d,%d)", tt.groups, tt.count)
		t.Run(testname, func(t *testing.T) {
			result := tt.deck.DealNN(tt.groups, tt.count)
			if !equal(result, tt.result) {
				t.Errorf("got %v, want %v", result, tt.result)
			}
		})
	}
}

func equal(a, b []Deck) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !Equals(v, b[i]) {
			return false
		}
	}
	return true
}
