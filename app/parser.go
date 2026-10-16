package main

import (
	"strings"
)

func parseCommand(command string) []string {
	result := []string{}
	var currentWord strings.Builder
	isBackSlash := false
	isSingleQuotes := false
	isDoubleQuotes := false
	for _, ch := range command {
		// если у нас был символ обратного слеша до этого
		// то мы текущий символ записываем в слово
		// и идем к следующему символу
		if isBackSlash {
			currentWord.WriteRune(ch)
			isBackSlash = false
			continue
		}

		if ch == '\\' && !isSingleQuotes {
			isBackSlash = true
			continue
		}

		if ch == '\'' && !isDoubleQuotes {
			isSingleQuotes = !isSingleQuotes
			continue
		} else if ch == '"' && !isSingleQuotes {
			isDoubleQuotes = !isDoubleQuotes
			continue
		} else if !isSingleQuotes && !isDoubleQuotes && ch == ' ' {
			if currentWord.Len() > 0 {
				result = append(result, currentWord.String())
				currentWord.Reset()
			}
		} else {
			currentWord.WriteRune(ch)
		}
	}
	if currentWord.Len() > 0 {
		result = append(result, currentWord.String())
	}
	return result
}
