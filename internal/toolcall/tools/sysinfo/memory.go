package sysinfo

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type SysMemoryTool struct{}

func (t SysMemoryTool) Name() string { return "sys_memory" }

func (t SysMemoryTool) Description() string {
	return "Return total, used, and free memory for the current machine."
}

func (t SysMemoryTool) Execute(_ string) string {
	if runtime.GOOS == "darwin" {
		return darwinMemory()
	}
	return linuxMemory()
}

func (t SysMemoryTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func darwinMemory() string {
	out, err := exec.Command("vm_stat").Output()
	if err != nil {
		return "Error: " + err.Error()
	}
	// Also grab total physical memory via sysctl
	total := ""
	if tot, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
		total = strings.TrimSpace(string(tot))
	}
	result := ""
	if total != "" {
		result = "Total physical memory (bytes): " + total + "\n\n"
	}
	return result + string(out)
}

func linuxMemory() string {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "Error: " + err.Error()
	}
	return string(data)
}
