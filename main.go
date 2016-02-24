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
	log.Printf("DinGo Code Generator\r\n")
	config := model.LoadConfiguration(configPath)
	schema := explorer.ExploreSchema(&config)
	pkg := producers.ProduceModelPackage(&config, schema)
	generators.GenerateModel(&config, pkg)
	log.Printf("Code generation done.\r\n")
}
