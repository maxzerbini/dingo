// +build postgres
// +build !mysql

package explorer

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/maxzerbini/dingo/model"
)

func TestPostgresExploreDatabase(t *testing.T) {
	t.Log("TestExploreDatabase started")
	confPostgres := model.LoadConfiguration("../config_postgres_local.json")
	exp := NewPostgreSqlExplorer()
	schema := exp.ExploreSchema(&confPostgres)
	if schema.SchemaName == "" {
		t.Fatalf("Test faile due to %s", "schema not found")
	}
	for _, t := range schema.Tables {
		fmt.Printf("Table %+v\r\n", *t)
		for _, c := range t.Columns {
			fmt.Printf("Colum %v\r\n", *c)
		}
	}
	for _, v := range schema.Views {
		fmt.Printf("View %+v\r\n", *v)
		for _, c := range v.Columns {
			fmt.Printf("Colum %v\r\n", *c)
		}
	}
}
