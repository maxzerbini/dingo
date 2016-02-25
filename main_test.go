package main

import (
	"log"
	"os"
	"testing"

	"github.com/maxzerbini/dingo/explorer"
	"github.com/maxzerbini/dingo/generators"
	"github.com/maxzerbini/dingo/model"
	"github.com/maxzerbini/dingo/producers"
)

var config *model.Configuration

func init() {
	config = &model.Configuration{Hostname: "localhost", Port: "3306", DatabaseName: "Customers", Username: "zerbo", Password: "Mysql.2016"}
	gopath := os.Getenv("GOPATH")
	testProjectPath := gopath + "/src/github.com/maxzerbini/prjtest"
	if _, err := os.Stat(testProjectPath); os.IsNotExist(err) {
		err = os.MkdirAll(testProjectPath, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", testProjectPath)
		}
	}
	config.OutputPath = testProjectPath
	config.BasePackage = "github.com/maxzerbini/prjtest"
}

func TestGeneration(t *testing.T) {
	t.Log("TestGenerateModel started")
	schema := explorer.ExploreSchema(config)
	modelpkg := producers.ProduceModelPackage(config, schema)
	daopkg := producers.ProduceDaoPackage(config, schema, modelpkg)
	generators.GenerateModel(config, modelpkg)
	generators.GenerateDao(config, daopkg)
}
