package sysinfo

import (
	"os/exec"
	"runtime"
)

type SysNetworkTool struct{}

func (t SysNetworkTool) Name() string { return "sys_network" }

func (t SysNetworkTool) Description() string {
	return "List active network interfaces and their IP addresses."
}

func (t SysNetworkTool) Execute(_ string) string {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("ifconfig")
	} else {
		cmd = exec.Command("ip", "addr")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "Error: " + err.Error() + "\n" + string(out)
	}
	return string(out)
}

func (t SysNetworkTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}
