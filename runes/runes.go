package runes

func Split(runes []rune, delimeter rune) [][]rune {
	var (
		result    = make([][]rune, 0)
		prevIndex int
	)
	for index, r := range runes {
		if r == delimeter {
			result = append(result, runes[prevIndex:index])
			prevIndex = index + 1
		}
	}
	if runes[len(runes)-1] == delimeter {
		result = append(result, runes[len(runes)-1:len(runes)-1])
	}
	return result
}

func Join(slice [][]rune, delimeter ...rune) []rune {
	var (
		result   = make([]rune, 0)
		sliceLen = len(slice)
	)
	for i, item := range slice {
		result = append(result, item...)
		if i < sliceLen-1 {
			result = append(result, delimeter...)
		}
	}
	return result
}

func Count(runes []rune, symbol rune) int {
	var count = 0
	for _, r := range runes {
		if r == symbol {
			count++
		}
	}
	return count
}

func ReplaceAll(runes []rune, original, new []rune) []rune {
	var result = make([]rune, 0, len(runes))

	for i := 0; i < len(runes); i++ {
		var wasFound = true
		for j, s := range original {
			if i+j >= len(runes) || runes[i+j] != s {
				wasFound = false
				break
			}
		}
		if wasFound {
			result = append(result, new...)
			i += len(original) - 1
		} else {
			result = append(result, runes[i])
		}
	}
	return result
}
