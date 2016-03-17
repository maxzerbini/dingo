package generators

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/maxzerbini/dingo/model"
)

const (
	customResourcesDirectory = ""
)

var (
	CustomResourcesTemplate = "./templates/customresources.tpl"
	CustomResourcesFile     = "customresources.go"
)

func GenerateCustomResources(config *model.Configuration) {
	// load template
	file, err := ioutil.ReadFile(CustomResourcesTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", CustomResourcesTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + customResourcesDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + CustomResourcesFile
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		log.Printf("Generating Custom Resources file %s\r\n", path)
		w := bufio.NewWriter(f)
		generateCode("", tpl, w)
		w.Flush()
	} else {
		log.Printf("Custom Resources file already exists, skipping generation\r\n")
	}
}
