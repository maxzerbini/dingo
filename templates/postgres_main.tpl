package main

import "encoding/json"
import "io/ioutil"
import "log"
import "fmt"
import "database/sql"
import _ "github.com/lib/pq"
import "github.com/gin-gonic/gin"
import "{{.BasePackage}}/{{.PackageName}}"
import "{{.BasePackage}}/dao"

type Configuration struct {
	DatabaseHostname              string
	DatabasePort                  string
	DatabaseName          	      string
	DatabaseUsername              string
	DatabasePassword              string
	DatabaseSchema	              string
	WebHost                       string
	WebPort                       string
	WebBaseHost                   string
}

func loadConfiguration(path string) Configuration {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("Configuration file not found at %s", path)
	}
	var jsontype Configuration
	if e = json.Unmarshal(file, &jsontype); e != nil {
		log.Fatalf("Invalid configuration file due to %s", e.Error())
	}
	return jsontype
}

func main(){
	conf := loadConfiguration("config.json")
	initConnection(conf)
	// Creates a router without any middleware by default
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	gin.SetMode(gin.DebugMode)
	registerAPI(conf, router)
	registerCustomResources(conf, router)
	router.Run(conf.WebHost + ":" + conf.WebPort)
}

func initConnection(conf Configuration) {
	var err error
	connString := fmt.Sprintf("user='%s' password='%s' dbname=%s host=%s port=%s sslmode=disable",
		conf.DatabaseUsername, conf.DatabasePassword, conf.DatabaseName, conf.DatabaseHostname, conf.DatabasePort)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	} else {
		dao.Connection = conn
	}
}

func registerAPI(conf Configuration, router *gin.Engine) {
	{{range $i,$e := .ServiceTypes}}
	s{{$i}} := {{$e.PackageName}}.New{{$e.TypeName}}()
	router.GET(conf.WebBaseHost+"/{{$e.ResourceName}}", s{{$i}}.List)
	{{if .IsSimplePK}}router.POST(conf.WebBaseHost+"/{{$e.ResourceName}}", s{{$i}}.Insert)
	router.PUT(conf.WebBaseHost+"/{{$e.ResourceName}}/:id", s{{$i}}.Update)
	router.DELETE(conf.WebBaseHost+"/{{$e.ResourceName}}/:id", s{{$i}}.Delete)
	router.GET(conf.WebBaseHost+"/{{$e.ResourceName}}/:id", s{{$i}}.Find)
    {{else}}router.GET(conf.WebBaseHost+"/{{$e.ResourceName}}/count", s{{$i}}.Count)
	{{if $e.Biz.IsReadOnly}}{{else}}router.POST(conf.WebBaseHost+"/{{$e.ResourceName}}", s{{$i}}.Insert)
	router.PUT(conf.WebBaseHost+"/{{$e.ResourceName}}", s{{$i}}.Update)
	router.POST(conf.WebBaseHost+"/{{$e.ResourceName}}/delete", s{{$i}}.Delete)
	router.POST(conf.WebBaseHost+"/{{$e.ResourceName}}/find", s{{$i}}.Find)
	{{end}}{{end}}{{end}}
}
