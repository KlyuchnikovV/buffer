package line

type Line []rune

func (l *Line) Insert(r rune, position int) {
	var (
		temp = *l
	)

	(*l) = append(temp[:position], append([]rune{r}, temp[position:]...)...)
}

func (l *Line) Remove(position int) {
	if position < 0 || position >= len(*l) {
		return
	}
	var (
		temp = *l
	)
	(*l) = append(temp[:position], temp[position+1:]...)
}

func (l *Line) Data() []rune {
	return *l
}

func (l *Line) Line(position, length int) []rune {
	if position < 0 || position >= len(*l) || length < 0 || length >= position+len(*l) {
		return nil
	}
	return (*l)[position : position+length]
}

func (l *Line) Length() int {
	return len(*l)
}

func (l *Line) FindAll(runes []rune) []int {
	var (
		index, ok = l.FindNext(-1, runes)
		result    = make([]int, 0)
	)
	for ok {
		result = append(result, index)
		index, ok = l.FindNext(index, runes)
	}
	return result
}

func (l *Line) FindNext(previousPosition int, runes []rune) (int, bool) {
	if len(runes) == 0 {
		return -1, false
	}

	var temp = *l
	if previousPosition != -1 {
		temp = temp[previousPosition:]
	}

loop:
	for i, ch := range temp {
		if ch != runes[0] {
			continue
		}
		for j, r := range runes {
			if temp[i+j] != r {
				continue loop
			}
		}
		return i, true
	}
	return -1, false
}
