package config

import (
	"os"
	"fmt"
	"encoding/json"
)

type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name`
}

func Read() (Config, error) {
	filePath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("Error getting the home directory: %v", err)
	}
	filePath = filePath + "/.gatorconfig.json"

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading the gatoroconf file: %v", err)
	}

	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error unmarshalling configuration to struct: %v", err)
	}

	return config, nil
}