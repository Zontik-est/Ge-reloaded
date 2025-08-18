package main

import (
    "regexp"
    "strconv"
    "strings"
)
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
