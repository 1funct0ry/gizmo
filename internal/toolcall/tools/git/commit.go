package git

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type CommitTool struct{}

func (t CommitTool) Name() string { return "git_commit" }

func (t CommitTool) Description() string {
	return "Stage all changes (git add -A) and create a commit with the given message. Destructive — requires confirmation."
}

func (t CommitTool) Execute(line string) string {
	var args struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(line), &args); err != nil {
		return "Error: invalid arguments: " + err.Error()
	}
	if strings.TrimSpace(args.Message) == "" {
		return "Error: commit message must not be empty"
	}

	if out, err := exec.Command("git", "add", "-A").CombinedOutput(); err != nil {
		return "Error staging files: " + err.Error() + "\n" + string(out)
	}

	out, err := exec.Command("git", "commit", "-m", args.Message).CombinedOutput()
	if err != nil {
		return "Error committing: " + err.Error() + "\n" + string(out)
	}
	return string(out)
}

func (t CommitTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"message": map[string]any{
				"type":        "string",
				"description": "The commit message.",
			},
		},
		"required": []string{"message"},
	}
}
