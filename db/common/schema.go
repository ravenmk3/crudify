package common

import (
	"io"
)

type DataType string

const (
	DataTypeBoolean   DataType = "bool"
	DataTypeByte      DataType = "byte"
	DataTypeInt16     DataType = "int16"
	DataTypeInt24     DataType = "int24"
	DataTypeInt32     DataType = "int32"
	DataTypeInt64     DataType = "int64"
	DataTypeFloat     DataType = "float"
	DataTypeDouble    DataType = "double"
	DataTypeDecimal   DataType = "decimal"
	DataTypeCurrency  DataType = "currency"
	DataTypeDate      DataType = "date"
	DataTypeTime      DataType = "time"
	DataTypeYear      DataType = "year"
	DataTypeDateTime  DataType = "datetime"
	DataTypeTimeStamp DataType = "timestamp"
	DataTypeEnum      DataType = "enum"
	DataTypeSet       DataType = "set"
	DataTypeGuid      DataType = "guid"
	DataTypeUuid      DataType = "uuid"
	DataTypeString    DataType = "string"
	DataTypeJson      DataType = "json"
	DataTypeXml       DataType = "xml"
	DataTypeBinary    DataType = "binary"
	DataTypeAny       DataType = "any"
)

type ColumnSchema struct {
	Name            string
	DataType        DataType
	NativeType      string
	MaxLength       int
	IsNullable      bool
	IsAutoIncrement bool
	IsUnsigned      bool
	Precision       int
	Scale           int
	HasDefault      bool
	IsPrimaryKey    bool
	Comment         string
}

func (s *ColumnSchema) CSharpDataType() string {
	t, ok := csharpTypeMap[s.DataType]
	if ok {
		return t
	}
	return "object"
}

func (s *ColumnSchema) JavaDataType() string {
	t, ok := javaTypeMap[s.DataType]
	if ok {
		return t
	}
	return "Object"
}

func (s *ColumnSchema) GoDataType() string {
	t, ok := goTypeMap[s.DataType]
	if ok {
		return t
	}
	return "any"
}

func (s *ColumnSchema) PythonDataType() string {
	t, ok := pythonTypeMap[s.DataType]
	if ok {
		return t
	}
	return "any"
}

type TableSchema struct {
	Name    string
	Columns []*ColumnSchema
	Comment string
}

func (s *TableSchema) PrimaryKeyColumn() *ColumnSchema {
	for _, column := range s.Columns {
		if column.IsPrimaryKey {
			return column
		}
	}
	return nil
}

type SchemaProvider interface {
	io.Closer
	GetTables(database string) ([]*TableSchema, error)
}
