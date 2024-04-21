package engine

import (
	"strings"

	"crudify/utils"
)

type JsFunctions struct {
}

func (f *JsFunctions) ToLower(value string) string {
	return strings.ToLower(value)
}

func (f *JsFunctions) ToUpper(value string) string {
	return strings.ToUpper(value)
}

func (f *JsFunctions) ToPlural(word string) string {
	return utils.ToPlural(word)
}

func (f *JsFunctions) ToSingular(word string) string {
	return utils.ToSingular(word)
}

func (f *JsFunctions) SplitWords(value string, lower bool) []string {
	return utils.SplitWords(value, lower)
}

func (f *JsFunctions) ToCamelCase(value string) string {
	return utils.ToCamelCase(value)
}

func (f *JsFunctions) ToPascalCase(value string) string {
	return utils.ToPascalCase(value)
}

func (f *JsFunctions) ToSnakeCase(value string) string {
	return utils.ToSnakeCase(value)
}

func (f *JsFunctions) ToKebabCase(value string) string {
	return utils.ToKebabCase(value)
}

func (f *JsFunctions) ToPluralSnakeCase(value string) string {
	return utils.ToPluralSnakeCase(value)
}

func (f *JsFunctions) ToPluralKebabCase(value string) string {
	return utils.ToPluralKebabCase(value)
}

func (f *JsFunctions) ToPluralCamelCase(value string) string {
	return utils.ToPluralCamelCase(value)
}

func (f *JsFunctions) ToPluralPascalCase(value string) string {
	return utils.ToPluralPascalCase(value)
}
