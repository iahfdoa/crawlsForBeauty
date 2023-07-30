package util

import (
	"encoding/json"
	"os"
)

func SaveConfigToFile(filePath string, m map[string][]string) error {
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func LoadImagePathsFromConfigFile(configFilePath string) (map[string][]string, error) {
	// 从配置文件中加载 imagePaths 切片
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var imagePaths map[string][]string
	err = json.Unmarshal(data, &imagePaths)
	if err != nil {
		return nil, err
	}

	return imagePaths, nil
}
