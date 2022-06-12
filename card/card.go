package card

import (
	"fmt"
	"strconv"
)

// https://unicode-table.com/en/blocks/playing-cards/

type CardSuit uint8

const (
	Spades CardSuit = iota
	Diamonds
	Clubs
	Hearts
	Star
)

var suit_strings = []string{"♠︎", "♦︎", "♣︎", "♥︎", "☆"}

func (c CardSuit) String() string {
	return suit_strings[c]
}

type CardFace uint8

func (c CardFace) String() string {
	var v string
	var ok bool
	if v, ok = conversion_lut_from_uint[uint8(c)]; !ok {
		v = strconv.Itoa(int(c))
	}
	return v
}

const (
	Ace CardFace = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Joker
)

const SUIT_CARD_AMOUNT uint8 = uint8(King)

var conversion_lut_from_uint = map[uint8]string{
	1:  "A",
	11: "J",
	12: "Q",
	13: "K",
	14: "?",
}

func uint_to_str(cf uint8) string {
	return strconv.Itoa(int(cf))
}

type Card struct {
	suit CardSuit
	face CardFace
}

func CardFactory(s CardSuit, f CardFace) *Card {
	return &Card{s, f}
}

func (c *Card) Suit() CardSuit {
	return c.suit
}

func (f *Card) Face() CardFace {
	return f.face
}

func (c *Card) String() string {
	return fmt.Sprintf(" %s %v ", c.suit.String(), c.face.String())
}
