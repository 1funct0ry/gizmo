package sysinfo

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type SysInfoTool struct{}

func (t SysInfoTool) Name() string { return "sys_info" }

func (t SysInfoTool) Description() string {
	return "Return OS, kernel, hostname, uptime, and CPU count for the current machine."
}

func (t SysInfoTool) Execute(_ string) string {
	var sb strings.Builder

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	fmt.Fprintf(&sb, "Hostname: %s\n", hostname)
	fmt.Fprintf(&sb, "OS:       %s\n", runtime.GOOS)
	fmt.Fprintf(&sb, "Arch:     %s\n", runtime.GOARCH)
	fmt.Fprintf(&sb, "CPUs:     %d\n", runtime.NumCPU())

	if out, err := exec.Command("uname", "-srm").Output(); err == nil {
		fmt.Fprintf(&sb, "Kernel:   %s", strings.TrimSpace(string(out)))
		sb.WriteString("\n")
	}

	uptime := getUptime()
	if uptime != "" {
		fmt.Fprintf(&sb, "Uptime:   %s\n", uptime)
	}

	return sb.String()
}

func (t SysInfoTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func getUptime() string {
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("sysctl", "-n", "kern.boottime").Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(out))
	}
	out, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return ""
	}
	return fields[0] + "s"
}
