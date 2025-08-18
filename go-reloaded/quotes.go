
package main

import (
    "strings"
)

// Обработка кавычек
func processQuotes(text string) string {
	parts := strings.Split(text, "'")
	    if len(parts)%2 == 0 {
        return text // нечётное количество кавычек — не обрабатываем
    }

	
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