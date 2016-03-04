package generators

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/maxzerbini/dingo/model"
)

const (
	mainDirectory = ""
)

var (
	MainTemplate = "./templates/main.tpl"
	MainFile     = "main.go"
)

func GenerateMain(config *model.Configuration, pkg *model.ServicePackage) {
	// load template
	file, err := ioutil.ReadFile(MainTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", MainTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + mainDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + MainFile
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	log.Printf("Generating Main Application file %s\r\n", path)
	w := bufio.NewWriter(f)
	generateCode(pkg, tpl, w)
	w.Flush()
}
