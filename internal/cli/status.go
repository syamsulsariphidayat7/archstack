package cli

import (
	"fmt"
	"strings"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

func statusCmd(names []string) error {
	if len(names) == 0 {
		return printStatusTable(registry.AllTools())
	}

	var tools []registry.Tool
	hasError := false
	for _, name := range names {
		tool, ok := registry.GetTool(name)
		if !ok {
			fmt.Printf("[error] Tool tidak dikenal: %s\n", name)
			hasError = true
			continue
		}
		tools = append(tools, *tool)
	}
	if len(tools) > 0 {
		printStatusTable(tools)
		for _, tool := range tools {
			if tool.IsServiceOnly() {
				printServiceDetail(&tool)
			}
		}
	}
	if hasError {
		return fmt.Errorf("beberapa tool tidak dikenal")
	}
	return nil
}

func printStatusTable(tools []registry.Tool) error {
	versions := system.ToolVersionStrings(tools)
	fmt.Printf("%-10s %-8s %-14s %-10s %-12s %s\n", "NAMA", "SUMBER", "STATUS", "SERVICE", "VERSI", "DESKRIPSI")
	fmt.Println(strings.Repeat("-", 84))
	for _, tool := range tools {
		installed := system.IsInstalled(tool.Binary, tool.Pkg)
		status := "not installed"
		if installed {
			status = "installed"
		}
		svc := "-"
		if tool.Service != "" {
			state, _ := system.ServiceStatusFor(&tool)
			switch state {
			case system.ServiceRunning:
				svc = "running"
			case system.ServiceStopped:
				svc = "stopped"
			default:
				svc = "-"
			}
		}
		ver := "-"
		if v, ok := versions[tool.Name]; ok && v != "" {
			ver = v
		}
		fmt.Printf("%-10s %-8s %-14s %-10s %-12s %s\n", tool.Name, tool.From, status, svc, ver, tool.Desc)
	}
	return nil
}
