package main

import (
    "strings"
)

// Обработка артиклей a/an
// Обработка артиклей a/an
func processArticles(words []string) []string {
	for i := 0; i < len(words)-1; i++ {
		if strings.ToLower(words[i]) == "a" {
			nextWord := cleanWord(words[i+1])
			if startsWithVowelSound(nextWord) {
				if words[i] == "A" {
					words[i] = "An"
				} else {
					words[i] = "an"
				}
			}
		}
	}
	return words
}

// Проверяет, начинается ли слово с гласного звука (учитывает исключения)
func startsWithVowelSound(word string) bool {
	wordLower := strings.ToLower(word)
	if wordLower == "" {
		return false
	}

	if startsWithConsonantSoundException(wordLower) {
		return false
	}

	firstChar := rune(wordLower[0])
	vowels := "aeiou"

	if strings.ContainsRune(vowels, firstChar) {
		return true
	}

	if firstChar == 'h' && isVowelSoundH(wordLower) {
		return true
	}

	if isAbbreviationWithVowelSound(word) {
		return true
	}

	return false
}

// Проверяет слова, начинающиеся с согласного звука, несмотря на гласную букву
func startsWithConsonantSoundException(word string) bool {
	// Например, "university" произносится с [j] звуком в начале (йуниверсити)
	exceptions := []string{"university", "unicorn", "user", "european", "eulogy", "unique", "unit", "utility"}

	for _, ex := range exceptions {
		if strings.HasPrefix(word, ex) {
			return true
		}
	}
	return false
}

// Проверяет, произносится ли h как гласная (немое h)
func isVowelSoundH(word string) bool {
	silentH := []string{"hour", "honest", "honor", "heir", "herb"}
	for _, silent := range silentH {
		if strings.HasPrefix(word, silent) {
			return true
		}
	}
	return false
}

// Проверяет аббревиатуры, которые начинаются с гласного звука (например, "FBI" произносится как "эф-би-ай")
func isAbbreviationWithVowelSound(word string) bool {
	if len(word) == 0 {
		return false
	}

	if word == strings.ToUpper(word) {
		vowelSoundLetters := "AEFHILMNORSX"
		firstChar := rune(word[0])
		for _, c := range vowelSoundLetters {
			if c == firstChar {
				return true
			}
		}
	}
	return false
}



// Вспомогательная функция для удаления пунктуации из слова
func cleanWord(word string) string {
	return strings.TrimFunc(word, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'))
	})
}
