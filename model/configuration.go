package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	Hostname     string
	Port         string
	DatabaseName string
	Username     string
	Password     string
	BasePackage  string
	OutputPath   string
}

func LoadConfiguration(path string) Configuration {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("Configuration file not found at %s", path)
	}
	var jsontype Configuration
	json.Unmarshal(file, &jsontype)
	jsontype.OutputPath = correctOutputPath(jsontype.OutputPath)
	return jsontype
}

func correctOutputPath(path string) string {
	gopath := os.Getenv("GOPATH")
	path = strings.Replace(path, "$GOPATH", gopath, -1)
	path = strings.Replace(path, "%GOPATH%", gopath, -1)
	return path
}
