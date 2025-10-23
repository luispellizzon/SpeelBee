package games

// This is how the Game itself will work, then later pangramGame will implement this interface and we can branch to pangramSinglePlayer and pangramMultiPlayer concrete classes that will implement this interface for each type of game
type Game interface {
	Name() string
	Info() (letters []rune, center rune)
	Submit(word string) (valid bool, reason string, points int, total int, pangram bool)
}
