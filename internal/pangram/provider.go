package pangram

type Provider struct{}


func (Provider) Board() (GameBoard, error){
	return Board()
}