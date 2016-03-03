package generators

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/maxzerbini/dingo/model"
)

const (
	viewmodelDirectory = "/viewmodel"
)

var (
	ViewModelTemplate = "./templates/viewmodel.tpl"
	ViewModelFile     = "viewmodel.go"
)

func GenerateViewModel(config *model.Configuration, pkg *model.ViewModelPackage) {
	// load template
	file, err := ioutil.ReadFile(ViewModelTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", ViewModelTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + viewmodelDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + ViewModelFile
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	log.Printf("Generating View Model file %s\r\n", path)
	w := bufio.NewWriter(f)
	generateCode(pkg, tpl, w)
	w.Flush()
}
