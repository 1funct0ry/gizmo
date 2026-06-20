package fs

import (
	"fmt"
	"os"
)

type CWDTool struct{}

func (C CWDTool) Name() string {
	return "current_directory"
}

func (C CWDTool) Description() string {
	return "Returns the current working directory."
}

func (C CWDTool) Execute(args string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		return "current directory could not be determined"
	}
	return fmt.Sprintf("Current working directory is %s", currentDir)
}

func (C CWDTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	}
}
