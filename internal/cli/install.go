package cli

import (
	"fmt"
	"os"

	"github.com/syamsulsariphidayat7/archstack/internal/registry"
	"github.com/syamsulsariphidayat7/archstack/internal/system"
)

func listCmd() error {
	return printStatusTable(registry.AllTools())
}

func installCmd(names []string) error {
	hasError := false
	for _, name := range names {
		tool, ok := registry.GetTool(name)
		if !ok {
			fmt.Fprintf(os.Stderr, "[error] Tool tidak dikenal: %s\n", name)
			hasError = true
			continue
		}
		if system.IsInstalled(tool.Binary, tool.Pkg) {
			fmt.Printf("[ok] %s sudah terinstal\n", name)
			continue
		}
		if err := installOne(tool); err != nil {
			fmt.Fprintf(os.Stderr, "[error] %s: %v\n", name, err)
			hasError = true
		}
	}
	if hasError {
		return fmt.Errorf("beberapa tool gagal diinstal")
	}
	return nil
}

func installOne(tool *registry.Tool) error {
	fmt.Printf("[install] Menginstal %s...\n", tool.Name)
	switch tool.From {
	case registry.SourcePacman:
		return system.InstallPacman(tool.Pkg)
	case registry.SourceYay:
		return system.InstallYay(tool.Pkg)
	}
	return nil
}
