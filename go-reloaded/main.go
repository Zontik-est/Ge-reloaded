package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Читаем входной файл
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", inputFile, err)
		os.Exit(1)
	}

	// Обрабатываем текст
	processedText := processText(string(content))

	// Записываем результат в выходной файл
	err = ioutil.WriteFile(outputFile, []byte(processedText), 0644)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", outputFile, err)
		os.Exit(1)
	}
}

func processText(text string) string {
	// Разбиваем по строкам
	lines := strings.Split(text, "\n")
	var resultLines []string

	for _, line := range lines {
		// Если строка пустая — сохраняем как есть
		if strings.TrimSpace(line) == "" {
			resultLines = append(resultLines, "")
			continue
		}

		// Разбиваем на слова
		words := strings.Fields(line)

		// Обработка правил
		words = processHexBin(words)
		words = processCaseChanges(words)
		words = processArticles(words)

		// Собираем обратно строку
		lineResult := strings.Join(words, " ")
		lineResult = processPunctuation(lineResult)
		lineResult = processQuotes(lineResult)

		resultLines = append(resultLines, lineResult)
	}

	// Склеиваем строки обратно, сохраняя \n
	return strings.Join(resultLines, "\n")
}


// Обработка шестнадцатеричных и двоичных чисел
func processHexBin(words []string) []string {
	result := make([]string, 0)
	
	for i := 0; i < len(words); i++ {
		if words[i] == "(hex)" && i > 0 {
			// Конвертируем предыдущее слово из hex в decimal
			if hexVal, err := strconv.ParseInt(result[len(result)-1], 16, 64); err == nil {
				result[len(result)-1] = strconv.FormatInt(hexVal, 10)
			}
		} else if words[i] == "(bin)" && i > 0 {
			// Конвертируем предыдущее слово из bin в decimal
			if binVal, err := strconv.ParseInt(result[len(result)-1], 2, 64); err == nil {
				result[len(result)-1] = strconv.FormatInt(binVal, 10)
			}
		} else {
			result = append(result, words[i])
		}
	}
	
	return result
}

// Обработка изменений регистра
func processCaseChanges(words []string) []string {
	result := make([]string, 0)
	
	for i := 0; i < len(words); i++ {
		word := words[i]
		
		if word == "(up)" && len(result) > 0 {
			result[len(result)-1] = strings.ToUpper(result[len(result)-1])
		} else if word == "(low)" && len(result) > 0 {
			result[len(result)-1] = strings.ToLower(result[len(result)-1])
		} else if word == "(cap)" && len(result) > 0 {
			result[len(result)-1] = strings.Title(result[len(result)-1])
		} else if isNumberedCommand(word, "up") {
			count := extractNumber(word)
			applyCase(result, count, strings.ToUpper)
		} else if isNumberedCommand(word, "low") {
			count := extractNumber(word)
			applyCase(result, count, strings.ToLower)
		} else if isNumberedCommand(word, "cap") {
			count := extractNumber(word)
			applyCase(result, count, strings.Title)
		} else if isPartialCommand(word, i, words) {
			// Обрабатываем разбитые команды типа "(cap," "6)"
			cmdType, count := extractPartialCommand(word, i, words)
			if cmdType != "" {
				switch cmdType {
				case "up":
					applyCase(result, count, strings.ToUpper)
				case "low":
					applyCase(result, count, strings.ToLower)  
				case "cap":
					applyCase(result, count, strings.Title)
				}
				i++ // Пропускаем следующее слово (число с скобкой)
			} else {
				result = append(result, word)
			}
		} else if isNumberPart(word) {
			// Пропускаем числовую часть команды (она уже обработана)
			continue
		} else {
			result = append(result, word)
		}
	}
	
	return result
}

// Проверяет, является ли слово началом разбитой команды
func isPartialCommand(word string, index int, words []string) bool {
	if !strings.HasPrefix(word, "(") {
		return false
	}
	
	commands := []string{"up,", "low,", "cap,"}
	for _, cmd := range commands {
		if word == "("+cmd {
			// Проверяем, есть ли следующее слово с числом и скобкой
			if index+1 < len(words) && isNumberPart(words[index+1]) {
				return true
			}
		}
	}
	return false
}

// Проверяет, является ли слово числовой частью команды
func isNumberPart(word string) bool {
	if !strings.HasSuffix(word, ")") {
		return false
	}
	
	numPart := strings.TrimSuffix(word, ")")
	_, err := strconv.Atoi(numPart)
	return err == nil
}

// Извлекает команду и число из разбитой команды
func extractPartialCommand(word string, index int, words []string) (string, int) {
	if index+1 >= len(words) {
		return "", 1
	}
	
	nextWord := words[index+1]
	if !isNumberPart(nextWord) {
		return "", 1
	}
	
	// Извлекаем тип команды
	cmdType := ""
	if word == "(up," {
		cmdType = "up"
	} else if word == "(low," {
		cmdType = "low"
	} else if word == "(cap," {
		cmdType = "cap"
	}
	
	// Извлекаем число
	numStr := strings.TrimSuffix(nextWord, ")")
	if num, err := strconv.Atoi(numStr); err == nil {
		return cmdType, num
	}
	
	return cmdType, 1
}

// Проверяет, является ли слово командой с числом
func isNumberedCommand(word, command string) bool {
	pattern := `\(` + command + `,\s*\d+\)`
	matched, _ := regexp.MatchString(pattern, word)
	return matched
}

// Извлечение числа из команды типа (up, 2)
func extractNumber(command string) int {
	re := regexp.MustCompile(`\((?:up|low|cap),\s*(\d+)\)`)
	matches := re.FindStringSubmatch(command)
	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}
	return 1
}

// Применение изменения регистра к n предыдущим словам
func applyCase(words []string, count int, caseFunc func(string) string) {
	if len(words) == 0 {
		return
	}
	
	// Находим последние count слов (игнорируя пунктуацию)
	wordCount := 0
	for i := len(words) - 1; i >= 0 && wordCount < count; i-- {
		word := words[i]
		// Проверяем, является ли это словом (содержит буквы)
		if containsLetters(word) {
			words[i] = caseFunc(words[i])
			wordCount++
		}
	}
}

// Проверяет, содержит ли строка буквы
func containsLetters(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}

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
	exceptions := []string{"university", "unicorn", "user", "european", "eulogy"}
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

// Обработка пунктуации
func processPunctuation(text string) string {
	// Обрабатываем группы пунктуации типа ... или !! или !?
	// Убираем пробелы внутри групп
	re1 := regexp.MustCompile(`([.!?])\s+([.!?])`)
	for re1.MatchString(text) {
		text = re1.ReplaceAllString(text, "$1$2")
	}
	
	// Убираем пробелы перед пунктуацией
	re2 := regexp.MustCompile(`\s+([.,!?:;])`)
	text = re2.ReplaceAllString(text, "$1")
	
	// Добавляем пробел после пунктуации, если его нет и далее идет буква
	re3 := regexp.MustCompile(`([.,!?:;])([a-zA-Z])`)
	text = re3.ReplaceAllString(text, "$1 $2")
	
	return text
}

// Обработка кавычек
func processQuotes(text string) string {
	parts := strings.Split(text, "'")
	
	if len(parts) < 3 {
		return text
	}
	
	result := parts[0]
	
	for i := 1; i < len(parts); i += 2 {
		if i+1 < len(parts) {
			// Убираем пробелы вокруг текста в кавычках
			quoted := strings.TrimSpace(parts[i])
			result += "'" + quoted + "'" + parts[i+1]
		}
	}
	
	return result
}