package score

// Scorer interface used to implement a strategy pattern
type Scorer interface { Score(length int, pangram bool) int }

// BasicScorer is the most basic way of scoring a word.
type BasicScorer struct{}
func (BasicScorer) Score(n int, pangram bool) int {
	if n < 4 { return 0 }
	if n == 4 { return 1 }
	return n
}

// Bonus scorer is a more advanced way of adding bonus to a word. We use strategy to leverage the basicScorer without the need of implementing the same functionality, but extending the strategy to also add bonus points and points according to the word, if the word is valid, length.
type BonusScorer struct { Inner Scorer; Bonus int }
func (b BonusScorer) Score(n int, pangram bool) int {
	base := b.Inner.Score(n, pangram)
	if base == 0 { return 0 }
	if pangram { return base + b.Bonus }
	return base
}
