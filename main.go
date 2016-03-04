package main

import (
	"flag"
	"log"
	"runtime"
	"runtime/debug"

	"github.com/maxzerbini/dingo/explorer"
	"github.com/maxzerbini/dingo/generators"
	"github.com/maxzerbini/dingo/model"
	"github.com/maxzerbini/dingo/producers"
)

var configPath string = "./config.json"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(300)
	flag.StringVar(&configPath, "conf", "./config.json", "path of the file config.json")
}

// Start the code generator
func main() {
	flag.Parse()
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
	if !config.SkipDaoGeneration {
		generators.GenerateDao(&config, daopkg)
		if !config.SkipBizGeneration {
			generators.GenerateViewModel(&config, viewmodelpkg)
			generators.GenerateBiz(&config, bizpkg)
			if !config.SkipServiceGeneration {
				generators.GenerateService(&config, srvpkg)
				generators.GenerateConfig(&config)
				generators.GenerateMain(&config, srvpkg)
			}
		}
	}
	log.Printf("Code generation done.\r\n")
}
