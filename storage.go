package main

import (
	"encoding/json"
	"os"
)

func saveConfig() {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		panic(err)
	}

	configFile, err := os.Create("data/config.json")
	if err != nil {
		panic(err)
	}

	_ = json.NewEncoder(configFile).Encode(configs)
	_ = configFile.Close()
}

func loadConfig() {
	configFile, err := os.Open("data/config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}

	err = json.NewDecoder(configFile).Decode(&configs)
	if err != nil {
		panic(err)
	}

	_ = configFile.Close()
}
