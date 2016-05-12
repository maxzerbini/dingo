// +build postgres
// +build !mysql

package producers

import (
	"testing"

	"github.com/maxzerbini/dingo/explorer"
	"github.com/maxzerbini/dingo/model"
)

var configPostgres model.Configuration
var expPostgres explorer.DatabaseExplorer

func init() {
	configPostgres = model.LoadConfiguration("../config_postgres_local.json")
	expPostgres = explorer.NewPostgreSqlExplorer()
}

func TestPostgresProduceModelPackage(t *testing.T) {
	t.Log("TestPostgresProduceModelPackage started")
	schema := expPostgres.ExploreSchema(&configPostgres)
	pkg := ProduceModelPackage(&configPostgres, schema)

	t.Logf("PackageName = %s", pkg.PackageName)
	for _, mt := range pkg.ModelTypes {
		t.Logf("ModelType = %s", mt.TypeName)
	}
}

func TestPostgresProduceViewModelPackage(t *testing.T) {
	t.Log("TestPostgresProduceViewModelPackage started")
	schema := expPostgres.ExploreSchema(&configPostgres)
	pkg := ProduceViewModelPackage(&configPostgres, schema)

	t.Logf("PackageName = %s", pkg.PackageName)
	for _, mt := range pkg.ViewModelTypes {
		t.Logf("ModelType = %s", mt.TypeName)
	}
}
