package generators

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/maxzerbini/dingo/model"
)

const (
	confDirectory = ""
)

var (
	ConfTemplate = "./templates/config.tpl"
	ConfFile     = "config.json"
)

func GenerateConfig(config *model.Configuration) {
	// load template
	file, err := ioutil.ReadFile(ConfTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", ConfTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + confDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + ConfFile
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	log.Printf("Generating Configuration file %s\r\n", path)
	w := bufio.NewWriter(f)
	generateCode(config, tpl, w)
	w.Flush()
}
