package explorer

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maxzerbini/dingo/model"
)

var conf model.Configuration

func init() {
	conf = model.Configuration{Hostname: "localhost", Port: "3306", DatabaseName: "Customers", Username: "root", Password: ""}
}

func TestExploreDatabase(t *testing.T) {
	t.Log("TestExploreDatabase started")
	schema := ExploreSchema(&conf)
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
