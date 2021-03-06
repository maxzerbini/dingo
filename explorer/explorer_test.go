// +build mysql
// +build !postgres

package explorer

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maxzerbini/dingo/model"
)

var conf model.Configuration

func init() {
	conf = model.LoadConfiguration("../config.json")
}

func TestMySqlExploreDatabase(t *testing.T) {
	t.Log("TestExploreDatabase started")
	exp := NewMySqlExplorer()
	schema := exp.ExploreSchema(&conf)
	if schema.SchemaName == "" {
		t.Errorf("Test faile due to %s", "schema not found")
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
