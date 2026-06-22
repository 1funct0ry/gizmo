package git

import "os/exec"

type StatusTool struct{}

func (t StatusTool) Name() string { return "git_status" }

func (t StatusTool) Description() string {
	return "Return the current git repository status (staged, unstaged, and untracked files)."
}

func (t StatusTool) Execute(_ string) string {
	out, err := exec.Command("git", "status").CombinedOutput()
	if err != nil {
		return "Error: " + err.Error() + "\n" + string(out)
	}
	return string(out)
}

func (t StatusTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}
