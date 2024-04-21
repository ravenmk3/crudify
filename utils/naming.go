package utils

import (
	"regexp"
	"strings"

	"github.com/fatih/camelcase"
)

var reNotLetterOrDigit = regexp.MustCompile("[^A-Za-z0-9]+")

func SplitWords(value string, lower bool) []string {
	result := []string{}

	parts := reNotLetterOrDigit.Split(value, 100)
	for _, part := range parts {
		entries := camelcase.Split(part)
		for _, entry := range entries {
			if lower {
				entry = strings.ToLower(entry)
			}
			result = append(result, entry)
		}
	}

	return result
}

func ToSnakeCase(value string) string {
	words := SplitWords(value, true)
	return strings.Join(words, "_")
}

func ToKebabCase(value string) string {
	words := SplitWords(value, true)
	return strings.Join(words, "-")
}

func ToCamelCase(value string) string {
	words := SplitWords(value, true)
	for i, word := range words {
		if i == 0 {
			continue
		}
		words[i] = firstRuneToUpper(word)
	}
	return strings.Join(words, "")
}

func ToPascalCase(value string) string {
	words := SplitWords(value, true)
	for i, word := range words {
		words[i] = firstRuneToUpper(word)
	}
	return strings.Join(words, "")
}

func firstRuneToUpper(s string) string {
	runes := []rune(s)
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}
