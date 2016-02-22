package explorer

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maxzerbini/dingo/model"
)

var conf model.Configuration

func init() {
	conf = model.Configuration{Hostname: "localhost", Port: "3306", DatabaseName: "Customers", Username: "zerbo", Password: "Mysql.2016"}
}

func TestExploreDatabase(t *testing.T) {
	t.Log("TestExploreDatabase started")
	schema := ExploreSchema(&conf)
	t.Log("%v", *schema)
}
