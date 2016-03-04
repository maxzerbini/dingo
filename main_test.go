package main

import (
	"log"
	"testing"

	"github.com/maxzerbini/dingo/explorer"
	"github.com/maxzerbini/dingo/generators"
	"github.com/maxzerbini/dingo/model"
	"github.com/maxzerbini/dingo/producers"
)

func init() {
	configPath = "config.json"
}

func TestGeneration(t *testing.T) {
	log.Printf("DinGo Code Generator\r\n")
	log.Printf("Processing configuration file %s\r\n", configPath)
	config := model.LoadConfiguration(configPath)
	schema := explorer.ExploreSchema(&config)
	modelpkg := producers.ProduceModelPackage(&config, schema)
	daopkg := producers.ProduceDaoPackage(&config, schema, modelpkg)
	viewmodelpkg := producers.ProduceViewModelPackage(&config, schema)
	bizpkg := producers.ProduceBizPackage(&config, modelpkg, daopkg, viewmodelpkg)
	srvpkg := producers.ProduceServicePackage(&config, viewmodelpkg, bizpkg)
	generators.GenerateModel(&config, modelpkg)
	generators.GenerateDao(&config, daopkg)
	generators.GenerateViewModel(&config, viewmodelpkg)
	generators.GenerateBiz(&config, bizpkg)
	generators.GenerateService(&config, srvpkg)
	log.Printf("Code generation done.\r\n")
}
