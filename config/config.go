package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigPath struct {
	CachePath string `json:"cachePath"`
}

type Config struct {
	BFJAPI ConfigAPI     `json:"bfj_api_prod"`
	MIXAPI ConfigAPIMix  `json:"mix_api_test"`
	Path   ConfigPath    `json:"paths"`
	Auth   Authorization `json:"auth"`
}

type ConfigAPI struct {
	ApiUrl               string `json:"apiUrl"`
	ApiGetLastJournals   string `json:"apiGetLastJournals"`
	ApiGetjournal        string `json:"apiGetJournal"`
	ApiPostAuthTest      string `json:"apiPostAuthTest"`
	ApiPostAuthProd      string `json:"apiPostAuthProd"`
	ApiGetListBF         string `json:"apiGetListBF"`
	ApiGetChemCoxes      string `json:"apiGetChemCoxes"`
	ApiGetChemicalsSlags string `json:"apiGetChemicalsSlags"`
	ApiGetChemMaterials  string `json:"apiGetChemMaterials"`
	ApiGetTappings       string `json:"apiGetTappings"`
}

type ConfigAPIMix struct {
	ApiUrl                       string `json:"apiUrl"`
	ApiPostAuthTest              string `json:"apiPostAuthTest"`
	ApiGetLastJournals           string `json:"apiGetLastJournals"`
	ApiPostChemical              string `json:"apiPostChemical"`
	ApiPostLadleMovement         string `json:"apiPostLadleMovement"`
	ApiPostPouringBucketMovement string `json:"apiPostPouringBucketMovement"`
	ApiGetLadleMovement          string `json:"apiGetLadleMovement"`
}

type Authorization struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
}

var GlobalConfig *Config = getConfig()

func getConfig() *Config {
	file, err := os.Open("models/json/config.json")
	if err != nil {
		fmt.Printf("Error opening config file: %s", err)
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		fmt.Printf("error decoding config file: %s", err)
	}

	return config
}
