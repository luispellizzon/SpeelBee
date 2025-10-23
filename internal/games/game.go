package games

type Game interface {
	Name() string
	Info() (letters []rune, center rune)
	Submit(word string) (valid bool, reason string, points int, total int, pangram bool)
}
