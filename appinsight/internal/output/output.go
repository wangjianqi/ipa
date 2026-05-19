package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type Result struct {
	OK      bool        `json:"ok"`
	Command string      `json:"command"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewOK(command string, data interface{}) Result {
	return Result{
		OK:      true,
		Command: command,
		Data:    data,
	}
}

func NewError(command string, err string) Result {
	return Result{
		OK:      false,
		Command: command,
		Error:   err,
	}
}

func PrintJSON(r Result) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r)
}

func WriteToFile(r Result, path string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func WriteDataToFile(data interface{}, path string) error {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return os.WriteFile(path, raw, 0644)
}
