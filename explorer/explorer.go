package explorer

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maxzerbini/dingo/model"
)

func ExploreSchema(config *model.Configuration) (schema *model.DatabaseSchema) {
	conn, err := sql.Open("mysql", config.Username+":"+config.Password+"@tcp("+config.Hostname+":"+config.Port+")/information_schema")
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
		table := model.Table{}
		err := rows.Scan(&table.TableName)
		if err != nil {
			log.Fatal(err)
		}
		schema.Tables = append(schema.Tables, table)
	}
}

func readViews(conn *sql.DB, schema *model.DatabaseSchema) {
	q := "SELECT TABLE_NAME FROM information_schema.VIEWS Where TABLE_SCHEMA=? ORDER BY TABLE_NAME"
	rows, err := conn.Query(q, schema.SchemaName)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		view := model.View{}
		err := rows.Scan(&view.ViewName)
		if err != nil {
			log.Fatal(err)
		}
		schema.Views = append(schema.Views, view)
	}
}

func readColums(conn *sql.DB, schema *model.DatabaseSchema, tableName string, colums *[]model.Column) {

}
