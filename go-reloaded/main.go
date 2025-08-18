package main

import (
	"fmt"
	"io/ioutil"
	"os"
	// "regexp"
	// "strconv"
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
	processedText = strings.TrimSpace(processedText)


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







