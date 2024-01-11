package helpers

import (
	"encoding/json"
	"os"
)

func LoadJSON(filename string, v interface{}) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(file), &v)
	return err
}

func SaveJSON(filename string, data interface{}) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}
