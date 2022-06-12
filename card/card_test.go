package card

import (
	"fmt"
	"testing"
)

const STANDARD_DECK_SIZE = 52

// Filter out a lot of things so that you can check manually
// with faces -> very small num of faces
// without -> filter almost all etc.
// pos = bypass + (each_suit * s_i) + f_i
// (each_suit * s_i) this is happening too many times, get it out of the inner loop

// Find a way to objectify
func TestStandardDeck(t *testing.T) {

	test_name := "Standard Deck Test"

	t.Run(test_name, func(t *testing.T) {
		cg := StandardCardGenerator()
		d := NewDeck(cg)

		faces := all_faces()
		suits := standard_order_suit()
		each_suit := len(faces)

		var (
			i int
			// Checking if it is of the exact same order as it should be for a standard brand new deck
			card       Card
			group_step int = 0
		)
		for _, suit := range suits {
			for f_ind, face := range faces {
				i = group_step + f_ind
				card = d.cards[i]
				if face != card.face {
					t.Errorf("Deck card face at %d : %v, expected %v", i, card.face, face)
				} else if suit != card.suit {
					t.Errorf("Card suit at %d : %v, expected %v", i, card.suit, suit)
				}
			}
			group_step += each_suit

		}

	})

}

func TestWithSuitsStandardDeck(t *testing.T) {
	test_name := "Only Given Suit(s) Standard Deck Test"

	exp_size := 13

	t.Run(test_name, func(t *testing.T) {
		cg := CardGeneratorFactory(WithSuits(Hearts))
		d := NewDeck(cg)
		if d.Len() != exp_size {
			t.Errorf("There are %d cards in the deck, expected %d", len(d.cards), exp_size)
		}

		if ok, pos := check_if_deck_has_only_one_suit(d, Hearts); !ok {
			t.Errorf("There should have been only %ss in the deck, but found %s", Hearts.String(), d.cards[pos].suit.String())
		}

	})
}

func TestWithRemovedCardsStandardDeck(t *testing.T) {
	test_name := "Unwanted Cards Removal From Standard Deck Test"

	exp := 13
	cg := StandardCardGenerator()
	unwanted := make([]Card, STANDARD_DECK_SIZE-exp)

	for s_i, suit := range []CardSuit{Clubs, Diamonds, Spades} {
		for i := Ace; i < Joker; i++ {
			// (exp*s_i)+int(i)-1 <--- -1 because we start by iota+1 when we set the face constants
			unwanted[(exp*s_i)+int(i)-1] = *CardFactory(suit, i)
		}
	}

	t.Run(test_name, func(t *testing.T) {
		d := NewDeck(cg, WithRemovedCards(unwanted...))
		if d.Len() != exp {
			t.Errorf("There are %d cards in the deck, expected %d", len(d.cards), d.Len()-len(unwanted))
		}
		if ok, pos := check_if_deck_has_only_one_suit(d, Hearts); !ok {
			t.Errorf("There should have been only %ss in the deck, but found %s", Hearts.String(), d.cards[pos].suit.String())
		}
	})
}

func TestWithAddedCardsStandardDeck(t *testing.T) {
	test_name := "Added Cards To Standard Deck Test"

	cg := StandardCardGenerator()
	add_these := []Card{
		*CardFactory(Hearts, Ace),
		*CardFactory(Clubs, Nine),
		*CardFactory(Hearts, Jack),
		*CardFactory(Clubs, Queen),
		*CardFactory(Spades, King),
	}

	t.Run(test_name, func(t *testing.T) {
		exptected_size := STANDARD_DECK_SIZE + len(add_these)
		d := NewDeck(cg, WithAddedCards(add_these...))
		// print_deck(d)
		if d.Len() != exptected_size {
			t.Errorf("There are %d cards in the deck, expected %d", len(d.cards), exptected_size)
		}
		for ind, card := range add_these {
			if d.cards[STANDARD_DECK_SIZE+ind] != card {
				t.Errorf("The card is %s, expected %s", d.cards[STANDARD_DECK_SIZE+ind].String(), card.String())
			}
		}
	})
}

func TestWithRemovedFacesStandardDeck(t *testing.T) {
	test_name := "Unwanted Face Removal From Standard Deck Test"

	only_kings_size := 4
	faces_except_king := 12

	unwanted := make([]CardFace, faces_except_king)
	// only kings are to be in the deck
	for i := Ace; i < King; i++ {
		unwanted[int(i)-1] = i
	}

	t.Run(test_name, func(t *testing.T) {
		cg := CardGeneratorFactory(WithOutFaces(unwanted...))
		d := NewDeck(cg)
		// print_deck(d)
		if d.Len() != only_kings_size {
			t.Errorf("There are %d cards in the deck, expected less.", len(d.cards))
		}
		if ok, pos := check_if_deck_has_only_one_face(d, King); !ok {
			t.Errorf("The card %s is of different type, expected %s", d.cards[pos].String(), King.String())
		}
	})
}

func TestWithAddedFacesStandardDeck(t *testing.T) {
	test_name := "Added Faces To Standard Deck Test"

	add_these := []CardFace{Ace, CardFace(Five), Jack}

	t.Run(test_name, func(t *testing.T) {
		cg := CardGeneratorFactory(WithAdditionalFaces(add_these...))
		d := NewDeck(cg)
		// print_deck(d)
		if d.Len() != (STANDARD_DECK_SIZE + len(add_these)*len(cg.suits)) {
			t.Errorf("There are %d cards in the deck, expected more.", len(d.cards))
		}
		bypass_and_start_here := STANDARD_DECK_SIZE / 4
		for i, face := range add_these {
			if d.cards[bypass_and_start_here+i].face != face {
				t.Errorf("The card is %s, expected %s", d.cards[bypass_and_start_here+i].String(), face.String())
			}
		}
	})
}

func TestWithAddedSuitsStandardDeck(t *testing.T) {
	test_name := "Added Faces To Standard Deck Test"

	add_these := []CardSuit{Hearts}

	t.Run(test_name, func(t *testing.T) {
		cg := CardGeneratorFactory(WithExtendedSuits(add_these...))
		d := NewDeck(cg)
		// print_deck(d)
		if d.Len() != STANDARD_DECK_SIZE+(len(cg.faces)) {
			t.Errorf("There are %d cards in the deck, expected more.", len(d.cards))
		}
		for i := STANDARD_DECK_SIZE; i < d.Len(); i++ {
			if d.cards[i].suit != add_these[0] {
				t.Errorf("The card is %s, expected %s", d.cards[i].String(), add_these[0].String())
			}
		}
	})
}

func TestWithCustomFacesStandardDeck(t *testing.T) {
	test_name := "Only Wanted Face Cards Standard Deck Test"

	wanted := []CardFace{Ace}

	t.Run(test_name, func(t *testing.T) {
		cg := CardGeneratorFactory(WithFaces(wanted...))
		d := NewDeck(cg)
		if d.Len() != len(wanted)*4 {
			t.Errorf("There are %d cards in the deck, expected %d.", d.Len(), len(wanted)*4)
		}
		if ok, pos := check_if_deck_has_only_one_face(d, wanted[0]); !ok {
			t.Errorf("The card %s is of different face, expected %s", d.cards[pos].String(), wanted[0].String())
		}
	})
}

func TestMultipliedStandardDeckCards(t *testing.T) {
	test_name := "Multiple Standard Decks Tests - Bigger Deck"

	tests := []int{2, 3, 4}

	cg := StandardCardGenerator()

	for _, multiply := range tests {
		t.Run(test_name+fmt.Sprintf("-%d", multiply), func(t *testing.T) {
			d := NewDeck(cg, WithDeckTimes(multiply))
			// print_deck(d)
			if d.Len() != STANDARD_DECK_SIZE*multiply {
				t.Errorf("There are %d cards in the deck, expected %d", d.Len(), STANDARD_DECK_SIZE*multiply)
			}

			bypass := 0
			var (
				pos, i, suit_num int
				c                Card
			)
			for i < multiply {
				for _, suit := range cg.suits {
					for f_i, face := range cg.faces {
						pos = bypass + suit_num + f_i
						c = *CardFactory(suit, face)
						if d.cards[pos] != c {
							t.Errorf("Card at %d is  %s, expected %s", pos, d.cards[pos].String(), c.String())

						}
					}
					suit_num += len(cg.faces)
					suit_num %= STANDARD_DECK_SIZE
				}
				i++
				bypass += (STANDARD_DECK_SIZE)
			}
		})
	}

}

func TestWithNJokersStandardDeck(t *testing.T) {
	test_name := "N Jokers in the Standard Decks Test"

	cg := StandardCardGenerator()
	num_jokers := 3

	t.Run(test_name, func(t *testing.T) {
		d := NewDeck(cg, WithJokersOfNum(num_jokers))
		// print_deck(d)
		if d.Len() != STANDARD_DECK_SIZE+num_jokers {
			t.Errorf("There are %d cards in the deck, expected %d", d.Len(), STANDARD_DECK_SIZE+num_jokers)
		}
	})
}

func TestShuffleStandardDeck(t *testing.T) {
	test_name := "Shuffle Standard Decks Test"

	cg := StandardCardGenerator()
	d := NewDeck(cg)

	t.Run(test_name, func(t *testing.T) {
		Shuffle(d)
		// print_deck(d)
		if d.cards[0] == *CardFactory(Spades, Ace) {
			t.Errorf("Deck is not shuffled.")
		}
	})
}

func TestSortStandardDeck(t *testing.T) {
	test_name := "Sort Standard Decks Test"

	cg := StandardCardGenerator()

	d := NewDeck(cg)
	d_upper := d.Len() - 1
	// Do not use shuffle, do not rely on anything
	for i := 0; i < d.Len()/2; i++ {
		d.Swap(i, d_upper-i)
	}
	// print_deck(d)

	t.Run(test_name, func(t *testing.T) {
		StandardSort(d)
		// print_deck(d)
		if d.cards[0] != *CardFactory(Spades, Ace) {
			t.Errorf("Deck is not sorted properly.")
		}
	})
}

func TestCardComparison(t *testing.T) {

	tests := []struct {
		a, b Card
		exp  bool
	}{
		{*CardFactory(Hearts, Ace), *CardFactory(Hearts, Ace), true},
		{*CardFactory(Spades, Nine), *CardFactory(Spades, Nine), true},
		{*CardFactory(Clubs, Jack), *CardFactory(Clubs, Jack), true},
		{*CardFactory(Star, Joker), *CardFactory(Star, Joker), true},

		{*CardFactory(Star, Joker), *CardFactory(Hearts, Jack), false},
		{*CardFactory(Star, Joker), *CardFactory(Clubs, Jack), false},
		{*CardFactory(Hearts, Jack), *CardFactory(Diamonds, Jack), false},
		{*CardFactory(Hearts, Jack), *CardFactory(Diamonds, Ten), false},
		{*CardFactory(Hearts, Jack), *CardFactory(Spades, Queen), false},
		{*CardFactory(Hearts, Jack), *CardFactory(Hearts, King), false},

		{*CardFactory(Spades, Nine), *CardFactory(Clubs, Nine), false},
		{*CardFactory(Spades, Nine), *CardFactory(Hearts, Eight), false},
	}

	test_name := "Card Comparison Tests"

	for i, test := range tests {
		t.Run(test_name+fmt.Sprintf("-%d", i), func(t *testing.T) {
			if test.a != test.b && test.exp {
				t.Errorf("	%+v and %+v should have been equal \n", test.a.String(), test.b.String())
			} else if test.a == test.b && !test.exp {
				t.Errorf("	%+v and %+v should have been different \n", test.a.String(), test.b.String())
			}
		})
	}

}

func print_deck(d *Deck) {
	var count = 0
	var row = d.Len() / 4
	for _, c := range d.cards {
		fmt.Print(c.String())
		count++
		if count%row == 0 {
			println()
			count = 0
		}
	}
}

func check_if_deck_has_only_one_suit(d *Deck, suit CardSuit) (ret bool, pos int) {
	// default bool is false
	ret = true
	for c_ind, card := range d.cards {
		if card.suit != suit {
			ret = false
			pos = c_ind
			break
		}
	}
	return
}
func check_if_deck_has_only_one_face(d *Deck, face CardFace) (ret bool, pos int) {
	// default bool is false
	ret = true
	for c_ind, card := range d.cards {
		if card.face != face {
			ret = false
			pos = c_ind
			break
		}
	}
	return
}
