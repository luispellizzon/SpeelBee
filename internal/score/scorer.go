package score

type Scorer interface { Score(length int, pangram bool) int }

type BasicScorer struct{}
func (BasicScorer) Score(n int, pangram bool) int {
	if n < 4 { return 0 }
	if n == 4 { return 1 }
	return n
}

type BonusScorer struct { Inner Scorer; Bonus int }
func (b BonusScorer) Score(n int, pangram bool) int {
	base := b.Inner.Score(n, pangram)
	if base == 0 { return 0 }
	if pangram { return base + b.Bonus }
	return base
}
