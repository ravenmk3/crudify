package engine

import (
	"crudify/schema/common"
)

type GlobalTemplateData struct {
	Vars   map[string]any
	Tables []*common.TableSchema
}

type EntityTemplateData struct {
	Global *GlobalTemplateData
	Vars   map[string]any
	Table  *common.TableSchema
}
