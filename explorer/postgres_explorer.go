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
	q := "SELECT c.TABLE_NAME, c.COLUMN_NAME, c.IS_NULLABLE, c.DATA_TYPE, c.CHARACTER_MAXIMUM_LENGTH,c.NUMERIC_PRECISION, c.NUMERIC_SCALE, c.UDT_NAME, c.COLUMN_DEFAULT, kc.CONSTRAINT_NAME"
	q += " FROM information_schema.COLUMNS c"
	q += " LEFT JOIN information_schema.KEY_COLUMN_USAGE kc ON kc.TABLE_CATALOG = c.TABLE_CATALOG AND kc.TABLE_SCHEMA = c.TABLE_SCHEMA AND kc.TABLE_NAME = c.TABLE_NAME AND kc.COLUMN_NAME = c.COLUMN_NAME"
	q += " WHERE c.TABLE_CATALOG=$1 AND c.TABLE_SCHEMA=$2 AND c.TABLE_NAME=$3 ORDER BY c.ORDINAL_POSITION"
	rows, err := conn.Query(q, config.DatabaseName, schema.SchemaName, tableName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		column := &model.Column{}
		nullable := "NO"
		var extra, primary sql.NullString
		err := rows.Scan(&column.TableName, &column.ColumnName, &nullable, &column.DataType, &column.CharacterMaximumLength, &column.NumericPrecision, &column.NumericScale, &column.ColumnType, &extra, &primary)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("Examining column %s\r\n", column.ColumnName)
		if "NO" == nullable {
			column.IsNullable = false
		} else {
			column.IsNullable = true
		}
		if extra.Valid && strings.HasPrefix(extra.String, "nextval") {
			column.IsAutoIncrement = true
		}
		if primary.Valid {
			column.IsPrimaryKey = true
		}
		*colums = append(*colums, column)
	}
}
