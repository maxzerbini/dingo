package explorer

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maxzerbini/dingo/model"
)

func ExploreSchema(config *model.Configuration) (schema *model.DatabaseSchema) {
	conn, err := sql.Open("mysql", config.Username+":"+config.Password+"@tcp("+config.Hostname+":"+config.Port+")/information_schema?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	schema = &model.DatabaseSchema{}
	schema.SchemaName = config.DatabaseName
	readTables(conn, schema)
	readViews(conn, schema)
	return schema
}

func readTables(conn *sql.DB, schema *model.DatabaseSchema) {
	q := "SELECT TABLE_NAME FROM information_schema.TABLES Where TABLE_SCHEMA=? ORDER BY TABLE_NAME"
	rows, err := conn.Query(q, schema.SchemaName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		table := &model.Table{}
		err := rows.Scan(&table.TableName)
		if err != nil {
			log.Fatal(err)
		}
		schema.Tables = append(schema.Tables, table)
		log.Printf("Examining table %s\r\n", table.TableName)
		readColums(conn, schema, table.TableName, &table.Columns)
	}
}

func readViews(conn *sql.DB, schema *model.DatabaseSchema) {
	q := "SELECT TABLE_NAME FROM information_schema.VIEWS Where TABLE_SCHEMA=? ORDER BY TABLE_NAME"
	rows, err := conn.Query(q, schema.SchemaName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		view := &model.View{}
		err := rows.Scan(&view.ViewName)
		if err != nil {
			log.Fatal(err)
		}
		schema.Views = append(schema.Views, view)
		log.Printf("Examining view %s\r\n", view.ViewName)
		readColums(conn, schema, view.ViewName, &view.Columns)
	}
}

func readColums(conn *sql.DB, schema *model.DatabaseSchema, tableName string, colums *[]*model.Column) {
	q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_TYPE, COLUMN_KEY, EXTRA"
	q += " FROM information_schema.COLUMNS "
	q += " WHERE TABLE_SCHEMA=? AND TABLE_NAME=? ORDER BY ORDINAL_POSITION"
	rows, err := conn.Query(q, schema.SchemaName, tableName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		column := &model.Column{}
		nullable := "NO"
		columnKey, extra := "", ""
		err := rows.Scan(&column.ColumnName, &column.ColumnName, &nullable, &column.DataType, &column.CharacterMaximumLength, &column.NumericPrecision, &column.NumericScale, &column.ColumnType, &columnKey, &extra)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("Examining column %s\r\n", column.ColumnName)
		if "NO" == nullable {
			column.IsNullable = false
		} else {
			column.IsNullable = true
		}
		if "PRI" == columnKey {
			column.IsPrimaryKey = true
		}
		if "UNI" == columnKey {
			column.IsUnique = true
		}
		if "auto_increment" == extra {
			column.IsAutoIncrement = true
		}
		*colums = append(*colums, column)
	}
}
