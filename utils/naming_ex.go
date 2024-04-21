package utils

import (
	"strings"
)

func ToPluralSnakeCase(value string) string {
	words := SplitWords(value, true)
	i := len(words) - 1
	words[i] = ToPlural(words[i])
	return strings.Join(words, "_")
}

func ToPluralKebabCase(value string) string {
	words := SplitWords(value, true)
	i := len(words) - 1
	words[i] = ToPlural(words[i])
	return strings.Join(words, "-")
}

func ToPluralCamelCase(value string) string {
	words := SplitWords(value, true)
	i := len(words) - 1
	words[i] = ToPlural(words[i])
	for i, word := range words {
		if i == 0 {
			continue
		}
		words[i] = firstRuneToUpper(word)
	}
	return strings.Join(words, "")
}

func ToPluralPascalCase(value string) string {
	words := SplitWords(value, true)
	i := len(words) - 1
	words[i] = ToPlural(words[i])
	for i, word := range words {
		words[i] = firstRuneToUpper(word)
	}
	return strings.Join(words, "")
}
