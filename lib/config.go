package lib

import (
	"io/ioutil"
	"encoding/json"
	"os"
)

type ConfigFileTranspilerDefine struct {
	Command string `json:"command"`
	Title   string `json:"title"`
}

type ConfigFileType struct {
	Transpilers []ConfigFileTranspilerDefine `json:"transpilers"`
	Scripts struct {
		Start   string `json:"start"`
		Restart string `json:"restart"`
		Stop    string `json:"stop"`
	} `json:"scripts"`
	Watch []string `json:"watch"`
}

var ConfigFile ConfigFileType

func LoadConfigFile() {
	ConfigFile = decodeConfigFile("build-config.json")

	if ConfigFile.Scripts.Start == "" {
		println("Config Error: scripts.start is required.")
		os.Exit(1)
	}
}

func decodeConfigFile(file string) ConfigFileType {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		println("Error: can not open " + file + " for read.")
		panic(err)
	}
	var jsondata ConfigFileType
	json.Unmarshal(content, &jsondata)

	return jsondata
}
