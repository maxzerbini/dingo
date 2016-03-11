package generators

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/maxzerbini/dingo/model"
)

const (
	modelDirectory = "/model"
)

var (
	ModelTemplate = "./templates/model.tpl"
	ModelFile     = "model.go"
)

func GenerateModel(config *model.Configuration, pkg *model.ModelPackage) {
	// load template
	file, err := ioutil.ReadFile(ModelTemplate)
	if err != nil {
		log.Fatalf("Can't read file in %s", ModelTemplate)
	}
	tpl := string(file)
	// open writer
	if _, err := os.Stat(config.OutputPath); os.IsNotExist(err) {
		log.Fatalf("Output path does not exist %s", config.OutputPath)
	}
	path := config.OutputPath + modelDirectory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Fatalf("Can not create directory %s", path)
		}
	}
	path += "/" + ModelFile
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	log.Printf("Generating Model file %s\r\n", path)
	w := bufio.NewWriter(f)
	generateCode(pkg, tpl, w)
	w.Flush()
}

func generateCode(pkg interface{}, tpl string, wr io.Writer) {
	tmpl, err := template.New("modelpackage").Parse(tpl)
	if err != nil {
		log.Println("Template parse failed")
		panic(err)
	}
	err = tmpl.Execute(wr, pkg)
	if err != nil {
		panic(err)
	}
}
