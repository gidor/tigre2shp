package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// the configuration
var config, defaults map[string]interface{}

func set() {
	if config == nil || defaults == nil {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		if config == nil {
			data, err := ioutil.ReadFile(filepath.Join(dir, "config.json"))
			if err == nil {
				err = json.Unmarshal(data, &config)
			}
		}

		if defaults == nil {
			data, err := ioutil.ReadFile(filepath.Join(dir, "defaults.json"))
			if err == nil {
				err = json.Unmarshal(data, &defaults)
			}
		}
	}
}

type confT map[string]interface{}

// Get  get a config item fo a list of keys te configuration
func Defaults(keys ...string) (interface{}, bool) {
	if defaults == nil {
		set()
	}
	return Navigate(interface{}(defaults), keys...)
}

// Get  get a config item fo a list of keys te configuration
func Get(keys ...string) (interface{}, bool) {
	if config == nil {
		set()
	}
	return Navigate(interface{}(config), keys...)
}

// Navigate an item for a list og keys
func Navigate(root interface{}, keys ...string) (interface{}, bool) {
	var (
		value interface{}
	)
	if root == nil {
		if config == nil {
			set()
		}
		value = interface{}(config)
	}
	value = root
	for _, key := range keys {
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil, false
		}
		v, ok := m[key]
		if !ok {
			return nil, false
		}
		value = interface{}(v)
	}
	return value, true
}

// GetArray
func GetArray(root interface{}, keys ...string) ([]interface{}, bool) {
	value, ok := Navigate(root, keys...)
	m, ok := value.([]interface{})
	if !ok {
		return nil, false
	}
	return m, ok
}

// GetArray
func GetMap(root interface{}, keys ...string) (map[string]interface{}, bool) {
	value, ok := Navigate(root, keys...)
	m, ok := value.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return m, ok
}

// GetInt
func GetInt(root interface{}, keys ...string) (int, bool) {
	value, ok := Navigate(root, keys...)
	if !ok {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case string:
		r, err := strconv.Atoi(v)
		if err == nil {
			return r, true
		} else {
			return 0, false
		}
	default:
		return 0, false
	}
}

// GetString
func GetString(root interface{}, keys ...string) (string, bool) {
	value, ok := Navigate(root, keys...)
	if !ok {
		return "", false
	}
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v), true
	case float64:
		return fmt.Sprintf("%f", v), true
	case string:
		return v, true
	default:
		return "", false
	}

}

// GetFloat
func GetFloat(root interface{}, keys ...string) (float64, bool) {
	value, ok := Navigate(root, keys...)
	if !ok {
		return 0, false
	}
	switch v := value.(type) {
	case int:
		return float64(v), true
	case float64:
		return v, true
	case string:
		var f float64
		_, err := fmt.Sscanf(v, "%f", f)
		return f, (err == nil)
	default:
		return 0.0, false
	}
}
