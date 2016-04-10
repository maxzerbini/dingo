package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	Hostname                string
	Port                    string
	DatabaseType            string
	DatabaseName            string
	Username                string
	Password                string
	BasePackage             string
	OutputPath              string
	ExcludedEntities        []string
	Entities                []string
	SkipDaoGeneration       bool
	SkipBizGeneration       bool
	SkipServiceGeneration   bool
	ForcePluralResourceName bool
	PostgresSchema          string
}

func (conf *Configuration) IsExcluded(name string) bool {
	for _, entity := range conf.ExcludedEntities {
		if entity == name {
			return true
		}
	}
	return false
}

func (conf *Configuration) IsIncluded(name string) bool {
	if len(conf.Entities) == 0 {
		return true
	} // if the list is void then all entities are included
	for _, entity := range conf.Entities {
		if entity == name {
			return true
		}
	}
	return false
}

func LoadConfiguration(path string) Configuration {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("Configuration file not found at path %s", path)
	}
	var jsontype Configuration
	if e = json.Unmarshal(file, &jsontype); e != nil {
		log.Fatalf("Invalid configuration file due to %s", e.Error())
	}
	jsontype.DatabaseType = checkDatabaseType(jsontype.DatabaseType)
	jsontype.OutputPath = correctOutputPath(jsontype.OutputPath)
	return jsontype
}

func correctOutputPath(path string) string {
	gopath := os.Getenv("GOPATH")
	path = strings.Replace(path, "$GOPATH", gopath, -1)
	path = strings.Replace(path, "%GOPATH%", gopath, -1)
	return path
}

func checkDatabaseType(databaseType string) string {
	db := strings.ToLower(databaseType)
	switch db {
	case "mysql":
		return db
	case "postgresql":
		return "postgres"
	case "postgres":
		return db

	}
	log.Fatalf("Unknow DatabaseType %s", databaseType)
	return ""
}
