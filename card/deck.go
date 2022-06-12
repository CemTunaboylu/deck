package card

import (
	"math/rand"
	"sort"
	"time"
)

type Deck struct {
	cards      []Card
	deal_index int
}

func (d Deck) Len() int      { return len(d.cards) }
func (d Deck) Swap(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] }

// StandardDeckOrder implements sort.Interface by providing Less and using the Len and Swap methods with StandardCardOrder.
// Since we implement the sort.Interface, sort does not have to use reflection to determine elms and length thus will be faster.
type ByStandardOrder struct{ *Deck }

func (d ByStandardOrder) Less(i, j int) bool {
	return StandardCardOrder(d.cards[i]) < StandardCardOrder(d.cards[j])
}

func StandardSort(d *Deck) {
	sort.Stable(ByStandardOrder{d})
}

func StandardCardOrder(c Card) uint8 {
	return uint8(c.suit)*SUIT_CARD_AMOUNT + uint8(c.face)
	/*
		Nazar would be angry about this huehhuhe
		Use some constant named maybe as SUIT_CARD_AMOUNT instead of `King`
	*/
}

func Shuffle(d *Deck) {
	rand.Seed(time.Now().UnixNano())
	// rand.Shuffle(d.Len(), func(i, j int) { a[i], a[j] = a[j], a[i] })
	rand.Shuffle(d.Len(), d.Swap)
}

type CardGenerator struct {
	suits []CardSuit
	faces []CardFace
}

func StandardCardGenerator() *CardGenerator {
	return &CardGenerator{
		suits: standard_order_suit(),
		faces: all_faces(),
	}
}

func CardGeneratorFactory(func_opts ...func(cg *CardGenerator)) *CardGenerator {
	cg := &CardGenerator{
		suits: standard_order_suit(),
		faces: all_faces(),
	}
	for _, f_o := range func_opts {
		f_o(cg)
	}
	return cg
}

func (cg *CardGenerator) Generate() *Deck {

	each_suit := len(cg.faces)
	cards := make([]Card, len(cg.suits)*each_suit)

	suit_group_step := 0
	for _, suit := range cg.suits {
		for f_i, face := range cg.faces {
			cards[suit_group_step+f_i] = *CardFactory(suit, face)
		}
		suit_group_step += each_suit
	}
	d := &Deck{cards, 0}
	return d
}

func WithSuits(suits ...CardSuit) func(*CardGenerator) {
	return func(cg *CardGenerator) {
		cg.suits = suits
	}
}

func WithExtendedSuits(suits ...CardSuit) func(*CardGenerator) {
	return func(cg *CardGenerator) {
		cg.suits = append(cg.suits, suits...)
	}
	/*
		additional suits may be misleading, the user may think it will be a complete new suits
	*/
}

// WithFaces changes the CardGenerator faces with specified face cards, e.g. {A, 5, J} is all the faces that generator will generate regardless of their suit
func WithFaces(faces ...CardFace) func(*CardGenerator) {
	return func(cg *CardGenerator) {
		cg.faces = faces
	}
}

func WithAdditionalFaces(faces ...CardFace) func(*CardGenerator) {
	return func(cg *CardGenerator) {
		cg.faces = append(cg.faces, faces...)
	}
}

// WithoutFaces removes the specified face cards from the generator, e.g. {A, 5, J} removes all the Aces, 5s and Jacks regardless of their suit
func WithOutFaces(faces ...CardFace) func(*CardGenerator) {
	return func(cg *CardGenerator) {
		new_faces := make([]CardFace, len(cg.faces)-len(faces))
		unwanted_set := map[CardFace]struct{}{}
		for _, c := range faces {
			unwanted_set[c] = struct{}{}
		}
		i := 0
		for _, f := range cg.faces {
			if _, in := unwanted_set[f]; !in {
				new_faces[i] = f
				i++
			}
		}
		cg.faces = new_faces
	}
}

func NewDeck(cg *CardGenerator, func_params ...func(*Deck)) *Deck {
	d := cg.Generate()
	for _, f := range func_params {
		f(d)
	}
	return d
}

// Functional options called after another will disturb the order of the deck.
// Thus we cannot rely on that assumption, we should iterate over the deck for the deck editing options.

// WithAddedCard is for special cases where we do not have control over the deck with the CardGenerator as a proxy to add some cards.
func WithAddedCards(cards_to_add ...Card) func(*Deck) {
	return func(d *Deck) {

		// 2nd arg is length, 3rd arg is capacity. This is the fastest way to 'deep' copy a slice.
		new_deck := append(make([]Card, 0, len(d.cards)+len(cards_to_add)), d.cards...)

		for _, c := range cards_to_add {
			new_deck = append(new_deck, c)
		}
		d.cards = new_deck

	}
}

//  WithRemoved is for special cases where we do not have control over the deck with the CardGenerator as a proxy
func WithRemovedCards(cards_to_remove ...Card) func(*Deck) {
	return func(d *Deck) {
		new_card_amt := len(d.cards) - len(cards_to_remove)
		new_deck := make([]Card, new_card_amt)
		unwanted_set := map[Card]struct{}{}
		for _, c := range cards_to_remove {
			unwanted_set[c] = struct{}{}
		}

		index := 0
		for _, card := range d.cards {
			if new_card_amt == 0 {
				break
			}
			if _, in := unwanted_set[card]; !in {
				new_deck[index] = card
				new_card_amt--
				index++
			}
		}
		d.cards = new_deck

	}
}

// WithDeckTimes increases the deck by multiples, 2x for example
func WithDeckTimes(times int) func(*Deck) {
	return func(d *Deck) {
		odd := times%2 == 1
		new_deck := []Card{}
		if odd {
			new_deck = make([]Card, len(d.cards))
			for in, card := range d.cards {
				new_deck[in] = card
			}
		}
		for times > 1 {
			d.cards = append(d.cards, d.cards...) // double
			times /= 2
		}

		if odd {
			d.cards = append(d.cards, new_deck...)
		}
	}
}

func WithJokersOfNum(num_jokers int) func(*Deck) {
	return func(d *Deck) {
		for num_jokers > 0 {
			d.cards = append(d.cards, *CardFactory(Star, Joker))
			num_jokers--
		}
	}
}

func WithCustomFilter(filter func(c Card) bool) func(d *Deck) {
	return func(d *Deck) {
		new_deck := []Card{} // we have to settle with the amortized cost for this one
		for _, card := range d.cards {
			if !filter(card) {
				new_deck = append(new_deck, card)
			}
		}
		d.cards = new_deck
	}
}

func all_faces() []CardFace {
	return []CardFace{
		Ace,
		Two,
		Three,
		Four,
		Five,
		Six,
		Seven,
		Eight,
		Nine,
		Ten,
		Jack,
		Queen,
		King,
	}
}

// A brand new deck of cards is typically sorted by suit, spades, diamonds, clubs, and hearts.
func standard_order_suit() []CardSuit {
	return []CardSuit{Spades, Diamonds, Clubs, Hearts}
}
