package main

import (
    "strconv"
)
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