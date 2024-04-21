package common

var javaTypeMap = map[DataType]string{
	DataTypeBoolean:   "Boolean",
	DataTypeByte:      "Byte",
	DataTypeInt16:     "Short",
	DataTypeInt24:     "Integer",
	DataTypeInt32:     "Integer",
	DataTypeInt64:     "Long",
	DataTypeFloat:     "Float",
	DataTypeDouble:    "Double",
	DataTypeDecimal:   "BigDecimal",
	DataTypeCurrency:  "BigDecimal",
	DataTypeDate:      "LocalDate",
	DataTypeTime:      "LocalTime",
	DataTypeYear:      "Year",
	DataTypeDateTime:  "Date",
	DataTypeTimeStamp: "Date",
	DataTypeEnum:      "String",
	DataTypeSet:       "String",
	DataTypeGuid:      "String",
	DataTypeUuid:      "String",
	DataTypeString:    "String",
	DataTypeJson:      "String",
	DataTypeXml:       "String",
	DataTypeBinary:    "Byte[]",
	DataTypeAny:       "Object",
}
