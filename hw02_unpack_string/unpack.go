package hw02unpackstring

import (
	"errors"
)

const (
	numberStart = 48
	numberLast  = 57
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	// var numberSymbols []string
	// for i := 0; i <= 9; i++ {
	//	//	numberSymbols = append(numberSymbols, string(i)) - error if more 57
	//	//	numberSymbols = append(numberSymbols, string(rune(i))) - error if more 57
	//	numberSymbols = append(numberSymbols, strconv.Itoa(i))
	// }

	var symbols = []rune(str)

	// var strBuilder strings.Builder
	var runesResult []rune
	for i := 0; i < len(symbols); i++ {
		if ok, err := isValidArrRunes(symbols, i); !ok {
			return "", err
		}

		if isNumber(symbols[i]) {
			// мы же данные провалидировали выше.
			// var count, errParse = strconv.Atoi(string(symbols[i]))
			// if errParse != nil {
			//	return "", errParse
			// }
			// по тз использовать string Repeats - хотя
			// strBuilder.WriteString(strings.Repeat(string(symbols[i-1]), count))
			var count = int(symbols[i]) - '0'

			// strBuilder.
			if count == 0 {
				runesResult = runesResult[0 : len(runesResult)-1]
			} else {
				for j := 1; j < count; j++ {
					// strBuilder.WriteRune(symbols[i-1])
					runesResult = append(runesResult, symbols[i-1])
				}
			}

		} else {
			runesResult = append(runesResult, symbols[i])
		}

	}
	return string(runesResult), nil
	// return strBuilder.String(), nil
}

func isValidArrRunes(runes []rune, currentIndex int) (bool, error) {
	if currentIndex == 0 {
		if isNumber(runes[0]) {
			return false, errors.New("first symbol is number")
		}

	} else {
		if isNumber(runes[currentIndex-1]) && isNumber(runes[currentIndex]) {
			return false, errors.New("previous symbol is number and current too")
		}
	}
	return true, nil
}

func isNumber(symbol rune) bool {
	return symbol >= numberStart && symbol <= numberLast
}
