package model

import (
	"database/sql"
	"strconv"
)

type DatabaseSchema struct {
	SchemaName string
	Tables     []*Table
	Views      []*View
}

type Table struct {
	TableName    string
	Columns      []*Column
	PrimaryKeys  []*Column
	OtherColumns []*Column
}

type View struct {
	ViewName string
	Columns  []*Column
}

type Column struct {
	TableName              string
	ColumnName             string
	IsPrimaryKey           bool
	IsNullable             bool
	IsAutoIncrement        bool
	IsUnique               bool
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
}

func (c Column) GetPostgresParam(i int) string {
	return "$" + strconv.Itoa(i+1)
}

func (c Column) GetPostgresParamFrom(i int, s int) string {
	return "$" + strconv.Itoa(s+i+1)
}
