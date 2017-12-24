package lib

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type ConfigFile struct {
	typescript struct {
		project string
		sources []string
	}
	scss    []string
	command string
}

func decodeConfigFile() ConfigFile {
	content, err := ioutil.ReadFile("compile.json")
	if err != nil {
		println("Error: can not open compile.json for read.")
		panic(err)
	}
	var jsondata ConfigFile
	json.Unmarshal(content, &jsondata)

	fmt.Printf("configFile: %v\n", jsondata)

	return jsondata
}
