package cache

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Data struct {
	Ids      []int         `yaml:"ids"`
	Tappings []map[int]int `yaml:"tappings"`
}

func ReadYAMLFile(fileName string) *Data {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil
	}

	var config Data
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil
	}

	return &config
}

func WriteYAMLFile(filename string, ids map[int][]int, tappings []map[int]int) {
	var config Data
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return
	}

	if len(ids) > 0 {
		config.Ids = nil
		for _, values := range ids {
			config.Ids = append(config.Ids, values...)
		}
	}

	if len(tappings) > 0 {
		config.Tappings = nil
		config.Tappings = append(config.Tappings, tappings...)
	} else {
		config.Tappings = nil
	}

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return
	}

	err = os.WriteFile(filename, yamlData, 0644)
	if err != nil {
		return
	}
}

func IdExists(config *Data, id int) bool {
	for _, existingID := range config.Ids {
		if existingID == id {
			return true
		}
	}
	return false
}

func DeleteIds(filename string, yaml *Data) *Data {
	ids := make(map[int][]int)
	if yaml != nil {
		yaml.Ids = []int{}
		WriteYAMLFile(filename, ids, nil)
	}
	return yaml
}

func (y *Data) ReplaceId(oldId, newId int) error {
	for i, id := range y.Ids {
		if id == oldId {
			y.Ids[i] = newId
		}
	}
	return nil
}

func TappingIdExists(config *Data, id int) bool {
	for _, values := range config.Tappings {
		for key := range values {
			if key == id {
				return true
			}
		}
	}
	return false
}

func FindTappingIdValue(config *Data, id int) int {
	for _, values := range config.Tappings {
		for key, v := range values {
			if key == id {
				return v
			}
		}
	}
	return 0
}

func UpdateTappingValue(config *Data, id, newValue int) error {
	for _, values := range config.Tappings {
		for key := range values {
			if key == id {
				values[id] = newValue
				return nil
			}
		}
	}
	return fmt.Errorf("tapping with ID %d not found", id)
}
