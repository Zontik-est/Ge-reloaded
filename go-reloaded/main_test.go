package main

import (
	"fmt"
	"strings"
	"testing"
)

// Debug функция для тестирования отдельных частей
func TestDebugCommands(t *testing.T) {
	text := "it was the age of foolishness (cap, 6) , it was"
	words := strings.Fields(text)
	
	fmt.Println("=== DEBUG: Command parsing ===")
	fmt.Println("Original:", text)
	fmt.Println("Words:", words)
	
	for i, word := range words {
		fmt.Printf("Word %d: '%s'\n", i, word)
		if isNumberedCommand(word, "cap") {
			fmt.Printf("  ✓ Is whole cap command\n")
		} else if isPartialCommand(word, i, words) {
			fmt.Printf("  ✓ Is partial cap command\n")
			cmdType, count := extractPartialCommand(word, i, words)
			fmt.Printf("  ✓ Command: %s, Count: %d\n", cmdType, count)
		} else if isNumberPart(word) {
			fmt.Printf("  ✓ Is number part\n")
		}
	}
	
	result := processCaseChanges(words)
	fmt.Println("After processCaseChanges:", result)
}

func TestDebugPunctuation(t *testing.T) {
	tests := []string{
		"I was sitting over there ,and then BAMM !!",
		"I was thinking ... You were right",
		"Punctuation tests are ... kinda boring ,what do you think ?",
	}
	
	fmt.Println("=== DEBUG: Punctuation ===")
	for _, test := range tests {
		fmt.Printf("Before: '%s'\n", test)
		result := processPunctuation(test)
		fmt.Printf("After:  '%s'\n\n", result)
	}
}

func TestProcessHexBin(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"1E", "(hex)", "files", "were", "added"},
			expected: []string{"30", "files", "were", "added"},
		},
		{
			input:    []string{"It", "has", "been", "10", "(bin)", "years"},
			expected: []string{"It", "has", "been", "2", "years"},
		},
		{
			input:    []string{"Simply", "add", "42", "(hex)", "and", "10", "(bin)"},
			expected: []string{"Simply", "add", "66", "and", "2"},
		},
	}

	for i, test := range tests {
		result := processHexBin(test.input)
		if !sliceEqual(result, test.expected) {
			t.Errorf("Test %d failed. Expected %v, got %v", i+1, test.expected, result)
		}
	}
}

func TestProcessCaseChanges(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"Ready,", "set,", "go", "(up)", "!"},
			expected: []string{"Ready,", "set,", "GO", "!"},
		},
		{
			input:    []string{"I", "should", "stop", "SHOUTING", "(low)"},
			expected: []string{"I", "should", "stop", "shouting"},
		},
		{
			input:    []string{"Welcome", "to", "the", "Brooklyn", "bridge", "(cap)"},
			expected: []string{"Welcome", "to", "the", "Brooklyn", "Bridge"},
		},
		{
			input:    []string{"This", "is", "so", "exciting", "(up, 2)"},
			expected: []string{"This", "is", "SO", "EXCITING"},
		},
	}

	for i, test := range tests {
		result := processCaseChanges(test.input)
		if !sliceEqual(result, test.expected) {
			t.Errorf("Test %d failed. Expected %v, got %v", i+1, test.expected, result)
		}
	}
}

func TestProcessArticles(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"There", "it", "was.", "A", "amazing", "rock!"},
			expected: []string{"There", "it", "was.", "An", "amazing", "rock!"},
		},
		{
			input:    []string{"There", "is", "no", "greater", "agony", "than", "bearing", "a", "untold", "story"},
			expected: []string{"There", "is", "no", "greater", "agony", "than", "bearing", "an", "untold", "story"},
		},
		{
			input:    []string{"I", "saw", "a", "elephant", "and", "a", "house"},
			expected: []string{"I", "saw", "an", "elephant", "and", "a", "house"},
		},
	}

	for i, test := range tests {
		result := processArticles(test.input)
		if !sliceEqual(result, test.expected) {
			t.Errorf("Test %d failed. Expected %v, got %v", i+1, test.expected, result)
		}
	}
}

func TestProcessPunctuation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "I was sitting over there ,and then BAMM !!",
			expected: "I was sitting over there, and then BAMM!!",
		},
		{
			input:    "I was thinking ... You were right",
			expected: "I was thinking... You were right",
		},
		{
			input:    "Punctuation tests are ... kinda boring ,what do you think ?",
			expected: "Punctuation tests are... kinda boring, what do you think?",
		},
	}

	for i, test := range tests {
		result := processPunctuation(test.input)
		if result != test.expected {
			t.Errorf("Test %d failed. Expected '%s', got '%s'", i+1, test.expected, result)
		}
	}
}

func TestProcessQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "I am exactly how they describe me: ' awesome '",
			expected: "I am exactly how they describe me: 'awesome'",
		},
		{
			input:    "As Elton John said: ' I am the most well-known homosexual in the world '",
			expected: "As Elton John said: 'I am the most well-known homosexual in the world'",
		},
	}

	for i, test := range tests {
		result := processQuotes(test.input)
		if result != test.expected {
			t.Errorf("Test %d failed. Expected '%s', got '%s'", i+1, test.expected, result)
		}
	}
}

func TestProcessText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.",
			expected: "It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.",
		},
		{
			input:    "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			expected: "Simply add 66 and 2 and you will see the result is 68.",
		},
	}

	for i, test := range tests {
		result := processText(test.input)
		if result != test.expected {
			t.Errorf("Test %d failed.\nExpected: '%s'\nGot:      '%s'", i+1, test.expected, result)
		}
	}
}

// Вспомогательная функция для сравнения слайсов
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestExtractNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"(up, 2)", 2},
		{"(low, 5)", 5},
		{"(cap, 10)", 10},
		{"(up, 1)", 1},
	}

	for i, test := range tests {
		result := extractNumber(test.input)
		if result != test.expected {
			t.Errorf("Test %d failed. Expected %d, got %d", i+1, test.expected, result)
		}
	}
}