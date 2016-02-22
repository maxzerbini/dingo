package model

import (
	"database/sql"
)

type DatabaseSchema struct {
	SchemaName string
	Tables     []Table
	Views      []View
}

type Table struct {
	TableName string
	Columns   []Column
}

type View struct {
	ViewName string
	Columns  []Column
}

type Column struct {
	ColumnName             string
	IsPrimaryKey           bool
	IsNullable             bool
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnKey              string // "PRI" for primary key
	Extra                  string // "auto_increment"
}
