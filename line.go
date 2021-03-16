package buffer

import (
	"github.com/KlyuchnikovV/buffer/runes"
)

func NewLine(data []rune) *Line {
	return &Line{data: data}
}

func NewLines(data []rune) []Line {
	var (
		lineRunes = runes.Split(data, '\n')
		result    = make([]Line, len(lineRunes))
	)
	for i, line := range lineRunes {
		result[i] = Line{data: line}
	}
	return result
}

func (l *Line) Insert(r rune, position int) {
	var temp = make([]rune, len(l.data)+1)
	copy(temp, l.data)

	l.data = append(temp[:position], append([]rune{r}, temp[position:]...)...)
}

func (l *Line) AppendData(runes []rune, position int) {
	var temp = make([]rune, len(l.data))
	copy(temp, l.data)

	l.data = append(temp[:position], append(runes, temp[position:]...)...)
}

func (l *Line) Remove(position int) {
	if position < 0 || position >= len(l.data) {
		return
	}
	var temp = make([]rune, len(l.data))
	copy(temp, l.data)

	l.data = append(temp[:position], temp[position+1:]...)
}

func (l Line) Data() []rune {
	return l.data
}

func (l Line) Line(position, length int) []rune {
	if position < 0 || position >= len(l.data) || length < 0 || length >= position+len(l.data) {
		return nil
	}
	return l.data[position : position+length]
}

func (l Line) Length() int {
	return len(l.data)
}

func (l Line) FindAll(runes []rune) []int {
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

func (l Line) FindNext(previousPosition int, runes []rune) (int, bool) {
	if len(runes) == 0 {
		return -1, false
	}

	var temp = l.data
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
