package cli

import (
	"fmt"
	"strings"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

func searchCmd(query string) error {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return fmt.Errorf("[error] Gunakan: archstack search <query>")
	}

	var found []registry.Tool
	for _, tool := range registry.AllTools() {
		hay := strings.ToLower(tool.Name + " " + tool.Pkg + " " + tool.Desc)
		if strings.Contains(hay, q) {
			found = append(found, tool)
		}
	}

	if len(found) == 0 {
		fmt.Printf("[info] Tidak ada tool yang cocok dengan '%s'\n", query)
		return nil
	}

	fmt.Printf("Hasil pencarian '%s' (%d tool):\n\n", query, len(found))
	for _, tool := range found {
		installed := system.IsInstalled(tool.Binary, tool.Pkg)
		status := "not installed"
		if installed {
			status = "installed"
		}
		fmt.Printf("  %-10s %-8s %-14s %s\n", tool.Name, tool.From, status, tool.Desc)
	}
	return nil
}
