package producers

import (
	"testing"

	"github.com/maxzerbini/dingo/explorer"
	"github.com/maxzerbini/dingo/model"
)

var config *model.Configuration
var exp explorer.DatabaseExplorer

func init() {
	config = &model.Configuration{Hostname: "localhost", Port: "3306", DatabaseName: "Customers", Username: "root", Password: ""}
	exp = explorer.NewMySqlExplorer()
}

func TestProduceModelPackage(t *testing.T) {
	t.Log("TestProduceModelPackage started")
	schema := exp.ExploreSchema(config)
	pkg := ProduceModelPackage(config, schema)

	t.Logf("PackageName = %s", pkg.PackageName)
	for _, mt := range pkg.ModelTypes {
		t.Logf("ModelType = %s", mt.TypeName)
	}
}

func TestProduceViewModelPackage(t *testing.T) {
	t.Log("TestProduceViewModelPackage started")
	schema := exp.ExploreSchema(config)
	pkg := ProduceViewModelPackage(config, schema)

	t.Logf("PackageName = %s", pkg.PackageName)
	for _, mt := range pkg.ViewModelTypes {
		t.Logf("ModelType = %s", mt.TypeName)
	}
}
