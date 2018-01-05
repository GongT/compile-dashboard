package lib

import (
	"io/ioutil"
	"encoding/json"
)

type configFile struct {
	Transpilers []struct {
		Command string `json:"command"`
		Title   string `json:"title"`
	} `json:"transpilers"`
	Scripts struct {
		Start   string `json:"start"`
		Restart string `json:"restart"`
		Stop    string `json:"stop"`
	} `json:"scripts"`
	Watch []string `json:"watch"`
}

var ConfigFile configFile

func init() {
	ConfigFile = decodeConfigFile("build-config.json")
}

func decodeConfigFile(file string) configFile {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		println("Error: can not open " + file + " for read.")
		panic(err)
	}
	var jsondata configFile
	json.Unmarshal(content, &jsondata)

	return jsondata
}
