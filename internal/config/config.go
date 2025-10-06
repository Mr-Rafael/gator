package config

import (
	"os"
	"fmt"
	"encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Error getting the config file path: %v", err)
	}

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

func SetUser(userName string) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Error getting the config file path: %v", err)
	}

	currentConfig, err := Read()
	if err != nil {
		return fmt.Errorf("Error reading the current configuration: %v", err)
	}
	currentConfig.CurrentUserName = userName

	jsonBytes, err := json.MarshalIndent(currentConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("Error trying to transform config to json string: %v", err)
	}

	return os.WriteFile(filePath, jsonBytes, 0644)
}

func getConfigFilePath() (string, error) {
	filePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting the home directory: %v", err)
	}
	filePath = filePath + "/" + configFileName
	return filePath, nil
}

func PrintCurrentConfig() {
	currentConfig, err := Read()
	if err != nil {
		fmt.Printf("\nError: Failed to read the current configuration. %v\n", err)
	}
	fmt.Printf("\n| Current Gator Configuration |\n")
	fmt.Printf("\nDatabase URL: %v", currentConfig.DBURL)
	fmt.Printf("\nCurrent Username: %v\n", currentConfig.CurrentUserName)
}