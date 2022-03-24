package utils

import "strings"

func ReplaceSpecialCharOfString(str string) string {
	result := ""

	strRune := []rune(str)

	for _, strRune := range strRune {
		currentChar := string(strRune)

		if strings.Contains("áàã", currentChar) {
			result += "a"
		} else if strings.Contains("éèẽ", currentChar) {
			result += "e"
		} else if strings.Contains("íìĩ", currentChar) {
			result += "i"
		} else if strings.Contains("óòõ", currentChar) {
			result += "o"
		} else if strings.Contains("úùũ", currentChar) {
			result += "u"
		} else if currentChar == "ç" {
			result += "c"
		} else if currentChar == "ñ" {
			result += "n"
		} else {
			result += currentChar
		}
	}

	return result
}
