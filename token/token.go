package token

type Token struct {
	value   []rune
	classes Classes
}

func New(runes []rune, class ...Class) Token {
	return Token{
		value:   runes,
		classes: class,
	}
}

func (token Token) Length() int {
	return len(token.value)
}

func (token *Token) Insert(position int, runes ...rune) {
	var temp = make([]rune, len(token.value)+1)
	copy(temp, token.value)

	token.value = append(temp[:position], append(runes, temp[position:]...)...)
}

func (token *Token) Remove(position int) {
	if position < 0 || position >= len(token.value) {
		return
	}
	var temp = make([]rune, len(token.value))
	copy(temp, token.value)

	token.value = append(temp[:position], temp[position+1:]...)
}
