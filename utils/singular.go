package utils

import (
	"strings"
)

var irregularSingularMap = map[string]string{
	"children":  "child",
	"people":    "person",
	"men":       "man",
	"women":     "woman",
	"mice":      "mouse",
	"geese":     "goose",
	"teeth":     "tooth",
	"feet":      "foot",
	"shoes":     "shoe",
	"knives":    "knife",
	"leaves":    "leaf",
	"wolves":    "wolf",
	"shelves":   "shelf",
	"loaves":    "loaf",
	"cacti":     "cactus",
	"foci":      "focus",
	"fungi":     "fungus",
	"nuclei":    "nucleus",
	"radii":     "radius",
	"bases":     "basis",
	"analyses":  "analysis",
	"diagnoses": "diagnosis",
	"theses":    "thesis",
	"crises":    "crisis",
}

func ToSingular(word string) string {
	if singular, ok := irregularSingularMap[word]; ok {
		return singular
	}

	if strings.HasSuffix(word, "es") && len(word) > 2 {
		c := word[len(word)-3]
		if c != 's' && c != 'x' && c != 'z' && c != 'h' {
			return word[:len(word)-2]
		}
	}

	if strings.HasSuffix(word, "s") {
		return word[:len(word)-1]
	}

	return word
}
