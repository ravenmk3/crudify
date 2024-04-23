package common

import (
	"crudify/utils"
)

func (s *ColumnSchema) NameCamelCase() string {
	return utils.ToCamelCase(s.Name)
}

func (s *ColumnSchema) NamePascalCase() string {
	return utils.ToPascalCase(s.Name)
}

func (s *ColumnSchema) NameSnakeCase() string {
	return utils.ToSnakeCase(s.Name)
}

func (s *ColumnSchema) NameKebabCase() string {
	return utils.ToKebabCase(s.Name)
}

func (s *TableSchema) NameCamelCase() string {
	return utils.ToCamelCase(s.Name)
}

func (s *TableSchema) NamePascalCase() string {
	return utils.ToPascalCase(s.Name)
}

func (s *TableSchema) NameSnakeCase() string {
	return utils.ToSnakeCase(s.Name)
}

func (s *TableSchema) NameKebabCase() string {
	return utils.ToKebabCase(s.Name)
}
