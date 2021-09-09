package token

type Token struct {
	Value   []rune
	Classes Classes
}

func New(runes []rune, class ...Class) Token {
	return Token{
		Value:   runes,
		Classes: class,
	}
}

func (token Token) Length() int {
	return len(token.Value)
}

func (token *Token) Insert(position int, runes ...rune) {
	var temp = make([]rune, len(token.Value))
	copy(temp, token.Value)

	token.Value = append(temp[:position], append(runes, temp[position:]...)...)
}

func (token *Token) Remove(position int) {
	if position < 0 || position >= len(token.Value) {
		return
	}
	var temp = make([]rune, len(token.Value))
	copy(temp, token.Value)

	token.Value = append(temp[:position], temp[position+1:]...)
}

func (token Token) Split(position int) (Token, Token) {
	var temp = make([]rune, len(token.Value))
	copy(temp, token.Value)
	return New(temp[:position], token.Classes...), New(temp[position:], token.Classes...)
}
