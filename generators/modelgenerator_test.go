// +build mysql
// +build !postgres

package generators

import (
	"log"
	"os"
	"testing"

	"github.com/maxzerbini/dingo/model"
)

var config *model.Configuration

func init() {
	gopath := os.Getenv("GOPATH")
	testProjectPath := gopath + "/src/github.com/maxzerbini/prjtest"
	if _, err := os.Stat(testProjectPath); os.IsNotExist(err) {
		err = os.MkdirAll(testProjectPath, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", testProjectPath)
		}
	}
	config = &model.Configuration{OutputPath: testProjectPath}
	ModelTemplate = "../templates/model.tpl"
}

func TestGenerateModel(t *testing.T) {
	t.Log("TestGenerateModel started")
	pkg := &model.ModelPackage{PackageName: "model"}
	mt := &model.ModelType{TypeName: "Customer"}
	mt.Fields = append(mt.Fields, &model.ModelField{FieldName: "Id", FieldType: "int"})
	mt.Fields = append(mt.Fields, &model.ModelField{FieldName: "Name", FieldType: "string"})
	mt.Fields = append(mt.Fields, &model.ModelField{FieldName: "CreationDate", FieldType: "time.Time"})
	pkg.ModelTypes = append(pkg.ModelTypes, mt)
	mt2 := &model.ModelType{TypeName: "Product"}
	mt2.Fields = append(mt2.Fields, &model.ModelField{FieldName: "Id", FieldType: "int"})
	mt2.Fields = append(mt2.Fields, &model.ModelField{FieldName: "ProductName", FieldType: "string"})
	mt2.Fields = append(mt2.Fields, &model.ModelField{FieldName: "CreationDate", FieldType: "time.Time"})
	pkg.ModelTypes = append(pkg.ModelTypes, mt2)
	mt3 := &model.ModelType{TypeName: "Order"}
	mt3.Fields = append(mt3.Fields, &model.ModelField{FieldName: "Id", FieldType: "int"})
	mt3.Fields = append(mt3.Fields, &model.ModelField{FieldName: "OrderDate", FieldType: "time.Time"})
	mt3.Fields = append(mt3.Fields, &model.ModelField{FieldName: "TotalCost", FieldType: "float64"})
	mt3.Fields = append(mt3.Fields, &model.ModelField{FieldName: "UpdateDate", FieldType: "time.Time"})
	pkg.ModelTypes = append(pkg.ModelTypes, mt3)
	pkg.ImportPackages = append(pkg.ImportPackages, "time")
	// pkg.ImportPackages = append(pkg.ImportPackages, "log")
	GenerateModel(config, pkg)
}
