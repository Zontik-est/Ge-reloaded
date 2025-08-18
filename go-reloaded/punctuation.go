
package main

import (
    "regexp"
)

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
