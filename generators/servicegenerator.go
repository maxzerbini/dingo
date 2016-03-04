package generators

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/maxzerbini/dingo/model"
)

const (
	serviceDirectory = "/service"
)

var (
	ServiceTemplate = "./templates/service.tpl"
	ServiceFile     = "service.go"
)

func GenerateService(config *model.Configuration, pkg *model.ServicePackage) {
	// load template
	file, err := ioutil.ReadFile(ServiceTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", ServiceTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + serviceDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + ServiceFile
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	log.Printf("Generating Service file %s\r\n", path)
	w := bufio.NewWriter(f)
	generateCode(pkg, tpl, w)
	w.Flush()
}
