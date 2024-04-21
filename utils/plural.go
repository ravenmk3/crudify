package utils

import (
	"strings"
)

var irregularPluralMap = map[string]string{
	"child":     "children",
	"person":    "people",
	"man":       "men",
	"woman":     "women",
	"mouse":     "mice",
	"goose":     "geese",
	"tooth":     "teeth",
	"foot":      "feet",
	"knife":     "knives",
	"leaf":      "leaves",
	"wolf":      "wolves",
	"shelf":     "shelves",
	"loaf":      "loaves",
	"cactus":    "cacti",
	"focus":     "foci",
	"fungus":    "fungi",
	"nucleus":   "nuclei",
	"radius":    "radii",
	"basis":     "bases",
	"analysis":  "analyses",
	"diagnosis": "diagnoses",
	"thesis":    "theses",
	"crisis":    "crises",
}

var pluralSuffixForEs = []string{"s", "x", "z", "sh", "ch"}

func ToPlural(word string) string {
	word = strings.ToLower(word)
	plural, ok := irregularPluralMap[word]
	if ok {
		return plural
	}

	for _, suffix := range pluralSuffixForEs {
		if strings.HasSuffix(word, suffix) {
			return word + "es"
		}
	}

	return word + "s"
}
