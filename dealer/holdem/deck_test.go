package holdemdealer

import (
	"testing"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/frand"
	"github.com/stretchr/testify/require"
)

func TestDeck(t *testing.T) {
	deck := newDeck()
	require.Equal(t, 52, len(deck.cards))

	rng := frand.NewUnsafeInt(42)

	m := make([]card.Card, 2)

	ok := fillrnd(rng, deck, m)
	require.True(t, ok)

	require.Equal(t, []card.Card{card.CardTS, card.Card8S}, m)
	require.Equal(t, 50, len(deck.cards))

	require.False(t, card.IsAnyMatch(m, deck.cards))

	ok = fillrnd(rng, deck, m)
	require.True(t, ok)
	require.Equal(t, 48, len(deck.cards))

	require.False(t, card.IsAnyMatch(m, deck.cards))

	m = make([]card.Card, 5)
	ok = fillrnd(rng, deck, m)
	require.True(t, ok)
	require.Equal(t, 43, len(deck.cards))

	require.False(t, card.IsAnyMatch(m, deck.cards))
}

func TestDeckReset(t *testing.T) {
	deck := newDeck()
	rng := frand.NewUnsafeInt(42)

	// Draw some cards
	m := make([]card.Card, 10)
	ok := fillrnd(rng, deck, m)
	require.True(t, ok)
	require.Equal(t, 42, len(deck.cards))

	// Reset the deck
	deck.reset()
	require.Equal(t, 52, len(deck.cards))

	// Verify all cards are present
	allCards := card.AllCopy()
	for _, c := range allCards {
		require.Contains(t, deck.cards, c)
	}
}

func TestCloneDeck(t *testing.T) {
	srcDeck := newDeck()
	rng := frand.NewUnsafeInt(42)

	// Draw some cards from source
	m := make([]card.Card, 5)
	ok := fillrnd(rng, srcDeck, m)
	require.True(t, ok)
	require.Equal(t, 47, len(srcDeck.cards))

	// Clone the deck
	destDeck := &deck{cards: make(card.Cards, 52)}
	cloneDeck(destDeck, srcDeck)

	// Verify the clone has the same cards
	require.Equal(t, len(srcDeck.cards), len(destDeck.cards))
	require.Equal(t, srcDeck.cards, destDeck.cards)

	// Verify modifying one doesn't affect the other
	n := make([]card.Card, 2)
	ok = fillrnd(rng, destDeck, n)
	require.True(t, ok)
	require.Equal(t, 45, len(destDeck.cards))
	require.Equal(t, 47, len(srcDeck.cards))
}

func TestPopcards(t *testing.T) {
	deck := newDeck()

	// Remove specific cards
	cardsToRemove := card.NewCardsFromString("As Ks Qs")
	ok := popcards(deck, cardsToRemove)
	require.True(t, ok)
	require.Equal(t, 49, len(deck.cards))

	// Verify cards were removed
	for _, c := range cardsToRemove {
		require.NotContains(t, deck.cards, c)
	}

	// Try to remove cards that don't exist in deck
	ok = popcards(deck, cardsToRemove)
	require.True(t, ok)                   // Should still return true even if cards not found
	require.Equal(t, 49, len(deck.cards)) // Length shouldn't change
}

func TestFillrndEmptyDeck(t *testing.T) {
	deck := &deck{cards: card.Cards{}}
	rng := frand.NewUnsafeInt(42)

	m := make([]card.Card, 2)
	ok := fillrnd(rng, deck, m)
	require.False(t, ok) // Should fail with empty deck
}

func TestPoprnd(t *testing.T) {
	rng := frand.NewUnsafeInt(42)

	// Test with empty slice
	var emptySlice []int
	val, newSlice, ok := poprnd(rng, emptySlice)
	require.False(t, ok)
	require.Equal(t, 0, val)
	require.Empty(t, newSlice)

	// Test with non-empty slice
	slice := []int{1, 2, 3, 4, 5}
	val, newSlice, ok = poprnd(rng, slice)
	require.True(t, ok)
	require.Equal(t, 4, len(newSlice))
	require.NotContains(t, newSlice, val) // Popped value should not be in new slice
}
