package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// the configuration
var config map[string]interface{}

func set() map[string]interface{} {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return nil
	}

	var config map[string]interface{}
	err = json.Unmarshal(data, &config)
	return config
}

// get te configuration
func Get() map[string]interface{} {
	if config == nil {
		config = set()
	}
	return config
}
