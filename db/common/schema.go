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

type TableSchema struct {
	Name    string
	Columns []ColumnSchema
	Comment string
}

type SchemaProvider interface {
	io.Closer
	GetTables(database string) ([]TableSchema, error)
}
