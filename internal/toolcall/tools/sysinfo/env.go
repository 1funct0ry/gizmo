package sysinfo

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
)

type SysEnvTool struct{}

func (t SysEnvTool) Name() string { return "sys_env" }

func (t SysEnvTool) Description() string {
	return "Print one environment variable by key, or all environment variables if no key is given. Args: key (optional string)."
}

func (t SysEnvTool) Execute(args string) string {
	key := ""
	if strings.TrimSpace(args) != "" {
		var params struct {
			Key string `json:"key"`
		}
		if err := json.Unmarshal([]byte(args), &params); err == nil {
			key = params.Key
		}
	}

	if key != "" {
		val := os.Getenv(key)
		if val == "" {
			return key + " is not set"
		}
		return key + "=" + val
	}

	env := os.Environ()
	sort.Strings(env)
	return strings.Join(env, "\n") + "\n"
}

func (t SysEnvTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"key": map[string]any{
				"type":        "string",
				"description": "Environment variable name to look up. Omit to list all variables.",
			},
		},
	}
}
