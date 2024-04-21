package mysql

import (
	"fmt"
	"strings"

	"crudify/db/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	DataTypeTinyInt   = "tinyint"
	DataTypeSmallInt  = "smallint"
	DataTypeMediumInt = "mediumint"
	DataTypeInt       = "int"
	DataTypeBigInt    = "bigint"
	DataTypeSingle    = "float"
	DataTypeReal      = "real"
	DataTypeDouble    = "double"
	DataTypeDecimal   = "decimal"
	DataTypeNumeric   = "numeric"
	DataTypeBit       = "bit"

	DataTypeDate      = "date"
	DataTypeTime      = "time"
	DataTypeYear      = "year"
	DataTypeDateTime  = "datetime"
	DataTypeTimeStamp = "timestamp"

	DataTypeChar       = "char"
	DataTypeVarChar    = "varchar"
	DataTypeBinary     = "binary"
	DataTypeVarBinary  = "varbinary"
	DataTypeTinyText   = "tinytext"
	DataTypeText       = "text"
	DataTypeMediumText = "mediumtext"
	DataTypeLongText   = "longtext"
	DataTypeTinyBlob   = "tinyblob"
	DataTypeBlob       = "blob"
	DataTypeMediumBlob = "mediumblob"
	DataTypeLongBlob   = "longblob"
	DataTypeEnum       = "enum"
	DataTypeSet        = "set"
	DataTypeJson       = "json"
)

var dataTypeMap = map[string]common.DataType{
	DataTypeTinyInt:    common.DataTypeByte,
	DataTypeSmallInt:   common.DataTypeInt16,
	DataTypeMediumInt:  common.DataTypeInt24,
	DataTypeInt:        common.DataTypeInt32,
	DataTypeBigInt:     common.DataTypeInt64,
	DataTypeSingle:     common.DataTypeFloat,
	DataTypeReal:       common.DataTypeFloat,
	DataTypeDouble:     common.DataTypeDouble,
	DataTypeDecimal:    common.DataTypeDecimal,
	DataTypeNumeric:    common.DataTypeDecimal,
	DataTypeBit:        common.DataTypeInt64,
	DataTypeDate:       common.DataTypeDate,
	DataTypeTime:       common.DataTypeTime,
	DataTypeYear:       common.DataTypeYear,
	DataTypeDateTime:   common.DataTypeDateTime,
	DataTypeTimeStamp:  common.DataTypeTimeStamp,
	DataTypeChar:       common.DataTypeString,
	DataTypeVarChar:    common.DataTypeString,
	DataTypeBinary:     common.DataTypeBinary,
	DataTypeVarBinary:  common.DataTypeBinary,
	DataTypeTinyText:   common.DataTypeString,
	DataTypeText:       common.DataTypeString,
	DataTypeMediumText: common.DataTypeString,
	DataTypeLongText:   common.DataTypeString,
	DataTypeTinyBlob:   common.DataTypeBinary,
	DataTypeBlob:       common.DataTypeBinary,
	DataTypeMediumBlob: common.DataTypeBinary,
	DataTypeLongBlob:   common.DataTypeBinary,
	DataTypeEnum:       common.DataTypeString,
	DataTypeSet:        common.DataTypeString,
	DataTypeJson:       common.DataTypeJson,
}

type MySqlTable struct {
	TableCatalog  string `db:"TABLE_CATALOG"`
	TableSchema   string `db:"TABLE_SCHEMA"`
	TableName     string `db:"TABLE_NAME"`
	TableType     string `db:"TABLE_TYPE"`
	Engine        string `db:"ENGINE"`
	Version       int    `db:"VERSION"`
	RowFormat     string `db:"ROW_FORMAT"`
	AutoIncrement *int64 `db:"AUTO_INCREMENT"`
	TableComment  string `db:"TABLE_COMMENT"`
}

type MySqlColumn struct {
	TableCatalog           string  `db:"TABLE_CATALOG"`
	TableSchema            string  `db:"TABLE_SCHEMA"`
	TableName              string  `db:"TABLE_NAME"`
	ColumnName             string  `db:"COLUMN_NAME"`
	OrdinalPosition        int     `db:"ORDINAL_POSITION"`
	ColumnDefault          *string `db:"COLUMN_DEFAULT"`
	IsNullable             string  `db:"IS_NULLABLE"`
	DataType               string  `db:"DATA_TYPE"`
	CharacterMaximumLength *int    `db:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   *int    `db:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision       *int    `db:"NUMERIC_PRECISION"`
	NumericScale           *int    `db:"NUMERIC_SCALE"`
	DatetimePrecision      *int    `db:"DATETIME_PRECISION"`
	CharacterSetName       *string `db:"CHARACTER_SET_NAME"`
	CollationName          *string `db:"COLLATION_NAME"`
	ColumnType             string  `db:"COLUMN_TYPE"`
	ColumnKey              string  `db:"COLUMN_KEY"`
	Extra                  string  `db:"EXTRA"`
	Privileges             string  `db:"PRIVILEGES"`
	ColumnComment          string  `db:"COLUMN_COMMENT"`
}

type mySqlSchemaProvider struct {
	db *sqlx.DB
}

func NewMySqlSchemaProvider(host string, port int, username, password string) (common.SchemaProvider, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/information_schema?charset=utf8mb4&parseTime=True&loc=%s",
		username, password, host, port, "Asia%2fShanghai")

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &mySqlSchemaProvider{
		db: db.Unsafe(),
	}, nil
}

func (me *mySqlSchemaProvider) Close() error {
	return me.db.Close()
}

func (me *mySqlSchemaProvider) GetTables(database string) ([]*common.TableSchema, error) {
	tableRows := []MySqlTable{}
	tableSql := "SELECT * FROM `information_schema`.`TABLES` WHERE `TABLE_SCHEMA` = ?"
	err := me.db.Select(&tableRows, tableSql, database)
	if err != nil {
		return nil, err
	}

	tables := []*common.TableSchema{}
	for _, row := range tableRows {
		columns, err := me.GetColumns(database, row.TableName)
		if err != nil {
			return nil, err
		}

		table := &common.TableSchema{
			Name:    row.TableName,
			Columns: columns,
			Comment: row.TableComment,
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (me *mySqlSchemaProvider) GetColumns(database, table string) ([]*common.ColumnSchema, error) {
	columnRows := []MySqlColumn{}
	columnSql := "SELECT * FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"
	err := me.db.Select(&columnRows, columnSql, database, table)
	if err != nil {
		return nil, err
	}

	columns := []*common.ColumnSchema{}
	for _, row := range columnRows {
		column := toColumnSchema(&row)
		columns = append(columns, column)
	}

	return columns, nil
}

func toColumnSchema(row *MySqlColumn) *common.ColumnSchema {
	dataType := getDataType(row.DataType)
	isNullable := strings.ToUpper(row.IsNullable) == "YES"
	isAutoIncr := strings.Contains(strings.ToLower(row.Extra), "auto_increment")
	isUnsigned := strings.Contains(strings.ToLower(row.ColumnType), "unsigned")
	hasDefault := row.ColumnDefault != nil
	isPrimaryKey := strings.ToUpper(row.ColumnKey) == "PRI"

	maxLength := -1
	if row.CharacterMaximumLength != nil {
		maxLength = *row.CharacterMaximumLength
	}

	precision := -1
	if row.NumericPrecision != nil {
		precision = *row.NumericPrecision
	} else if row.DatetimePrecision != nil {
		precision = *row.DatetimePrecision
	}

	scale := -1
	if row.NumericScale != nil {
		scale = *row.NumericScale
	}

	return &common.ColumnSchema{
		Name:            row.ColumnName,
		DataType:        dataType,
		NativeType:      row.DataType,
		MaxLength:       maxLength,
		IsNullable:      isNullable,
		IsAutoIncrement: isAutoIncr,
		IsUnsigned:      isUnsigned,
		Precision:       precision,
		Scale:           scale,
		HasDefault:      hasDefault,
		IsPrimaryKey:    isPrimaryKey,
		Comment:         row.ColumnComment,
	}
}

func getDataType(mysqlDataType string) common.DataType {
	mysqlDataType = strings.ToLower(mysqlDataType)
	dataType, ok := dataTypeMap[mysqlDataType]
	if ok {
		return dataType
	}
	return common.DataTypeAny
}
