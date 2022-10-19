package model

import (
	"container/list"
)

// ConfigTable 表配置
type ConfigTable struct {
	Id          int
	Name        string
	SourceId    int
	TargetId    int
	SourceTable string
	TargetTable string
	Interval    int
	LastTime    string
	WhereSql    string
	Field       *list.List
}

// ConfigTableField 表字段配置
type ConfigTableField struct {
	Id          int
	ConfigId    int
	SourceField string
	TargetField string
}
