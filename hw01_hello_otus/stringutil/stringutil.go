package stringutil

func Reverse(str string) string {
	arrChars := []rune(str)

	for i, j := 0, len(arrChars)-1; i < j; i, j = i+1, j-1 {
		arrChars[i], arrChars[j] = arrChars[j], arrChars[i]
	}

	return string(arrChars)
}
