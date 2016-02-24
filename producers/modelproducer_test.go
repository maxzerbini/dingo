package producers

import (
	"testing"

	"github.com/maxzerbini/dingo/explorer"
	"github.com/maxzerbini/dingo/model"
)

var config *model.Configuration

func init() {
	config = &model.Configuration{Hostname: "localhost", Port: "3306", DatabaseName: "Customers", Username: "zerbo", Password: "Mysql.2016"}
}

func TestProduceModelPackage(t *testing.T) {
	t.Log("TestGenerateModel started")
	schema := explorer.ExploreSchema(config)
	pkg := ProduceModelPackage(config, schema)

	t.Logf("PackageName = %s", pkg.PackageName)
	for _, mt := range pkg.ModelTypes {
		t.Logf("ModelType = %s", mt.TypeName)
	}
}
