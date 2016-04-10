package explorer

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/maxzerbini/dingo/model"
)

type PostgreSqlExplorer struct {
}

func NewPostgreSqlExplorer() *PostgreSqlExplorer {
	e := &PostgreSqlExplorer{}
	return e
}

func (e *PostgreSqlExplorer) ExploreSchema(config *model.Configuration) (schema *model.DatabaseSchema) {
	connString := fmt.Sprintf("user='%s' password='%s' dbname=%s host=%s port=%s sslmode=disable",
		config.Username, config.Password, config.DatabaseName, config.Hostname, config.Port)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	schema = &model.DatabaseSchema{}
	schema.SchemaName = config.PostgresSchema
	e.readTables(config, conn, schema)
	e.readViews(config, conn, schema)
	return schema
}

func (e *PostgreSqlExplorer) readTables(config *model.Configuration, conn *sql.DB, schema *model.DatabaseSchema) {
	q := "SELECT TABLE_NAME FROM information_schema.TABLES Where TABLE_CATALOG=$1 AND TABLE_SCHEMA=$2 AND TABLE_TYPE='BASE TABLE' ORDER BY TABLE_NAME"
	rows, err := conn.Query(q, config.DatabaseName, schema.SchemaName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		table := &model.Table{}
		err := rows.Scan(&table.TableName)
		if err != nil {
			log.Fatal(err)
		}
		if config.IsIncluded(table.TableName) && !config.IsExcluded(table.TableName) {
			schema.Tables = append(schema.Tables, table)
			log.Printf("Examining table %s\r\n", table.TableName)
			e.readColums(config, conn, schema, table.TableName, &table.Columns)
			for _, col := range table.Columns {
				if col.IsPrimaryKey {
					table.PrimaryKeys = append(table.PrimaryKeys, col)
				} else {
					table.OtherColumns = append(table.OtherColumns, col)
				}
			}
		} else {
			log.Printf("Table %s is excluded\r\n", table.TableName)
		}
	}
}

func (e *PostgreSqlExplorer) readViews(config *model.Configuration, conn *sql.DB, schema *model.DatabaseSchema) {
	q := "SELECT TABLE_NAME FROM information_schema.VIEWS Where TABLE_CATALOG=$1 AND TABLE_SCHEMA=$2 ORDER BY TABLE_NAME"
	rows, err := conn.Query(q, config.DatabaseName, schema.SchemaName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		view := &model.View{}
		err := rows.Scan(&view.ViewName)
		if err != nil {
			log.Fatal(err)
		}
		if config.IsIncluded(view.ViewName) && !config.IsExcluded(view.ViewName) {
			schema.Views = append(schema.Views, view)
			log.Printf("Examining view %s\r\n", view.ViewName)
			e.readColums(config, conn, schema, view.ViewName, &view.Columns)
		} else {
			log.Printf("View %s is excluded\r\n", view.ViewName)
		}
	}
}

func (e *PostgreSqlExplorer) readColums(config *model.Configuration, conn *sql.DB, schema *model.DatabaseSchema, tableName string, colums *[]*model.Column) {
	q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, UDT_NAME, COLUMN_DEFAULT"
	q += " FROM information_schema.COLUMNS "
	q += " WHERE TABLE_CATALOG=$1 AND TABLE_SCHEMA=$2 AND TABLE_NAME=$3 ORDER BY ORDINAL_POSITION"
	rows, err := conn.Query(q, config.DatabaseName, schema.SchemaName, tableName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		column := &model.Column{}
		nullable := "NO"
		var extra sql.NullString
		err := rows.Scan(&column.ColumnName, &column.ColumnName, &nullable, &column.DataType, &column.CharacterMaximumLength, &column.NumericPrecision, &column.NumericScale, &column.ColumnType, &extra)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("Examining column %s\r\n", column.ColumnName)
		if "NO" == nullable {
			column.IsNullable = false
		} else {
			column.IsNullable = true
		}
		if extra.Valid && strings.HasPrefix("nextval", extra.String) {
			column.IsAutoIncrement = true
		}
		*colums = append(*colums, column)
	}
}
